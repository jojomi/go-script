// forked from https://github.com/termie/go-shutil
package script

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type SameFileError struct {
	Src string
	Dst string
}

func (e SameFileError) Error() string {
	return fmt.Sprintf("%s and %s are the same file", e.Src, e.Dst)
}

type SpecialFileError struct {
	File     string
	FileInfo os.FileInfo
}

func (e SpecialFileError) Error() string {
	return fmt.Sprintf("`%s` is a named pipe", e.File)
}

type NotADirectoryError struct {
	Src string
}

func (e NotADirectoryError) Error() string {
	return fmt.Sprintf("`%s` is not a directory", e.Src)
}

type AlreadyExistsError struct {
	Dst string
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("`%s` already exists", e.Dst)
}

func samefile(fs afero.Fs, src string, dst string) bool {
	srcInfo, _ := fs.Stat(src)
	dstInfo, _ := fs.Stat(dst)
	return os.SameFile(srcInfo, dstInfo)
}

func specialfile(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CopyFile copies data from src to dst
func CopyFile(fs afero.Fs, src, dst string, followSymlinks bool) error {
	if samefile(fs, src, dst) {
		return &SameFileError{src, dst}
	}

	// Make sure src exists and neither are special files
	srcStat, err := fs.Stat(src)
	if err != nil {
		return err
	}
	if specialfile(srcStat) {
		return &SpecialFileError{src, srcStat}
	}

	dstStat, err := fs.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		if specialfile(dstStat) {
			return &SpecialFileError{dst, dstStat}
		}
	}

	// do the actual copy
	fsrc, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	size, err := io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}

	if size != srcStat.Size() {
		return fmt.Errorf("%s: %d/%d copied", src, size, srcStat.Size())
	}

	return nil
}

// CopyMode copies mode bits from src to dst.
func CopyMode(fs afero.Fs, src, dst string, followSymlinks bool) error {
	srcStat, err := fs.Stat(src)
	if err != nil {
		return err
	}

	// get the actual file stats
	srcStat, _ = fs.Stat(src)
	err = fs.Chmod(dst, srcStat.Mode())
	return err
}

// Copy data and mode bits ("cp src dst"). Return the file's destination.
//
// The destination may be a directory.
//
// If followSymlinks is false, symlinks won't be followed. This
// resembles GNU's "cp -P src dst".
//
// If source and destination are the same file, a SameFileError will be
// rased.
func Copy(fs afero.Fs, src, dst string, followSymlinks bool) (string, error) {
	dstInfo, err := fs.Stat(dst)

	if err == nil && dstInfo.Mode().IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err != nil && !os.IsNotExist(err) {
		return dst, err
	}

	err = CopyFile(fs, src, dst, followSymlinks)
	if err != nil {
		return dst, err
	}

	err = CopyMode(fs, src, dst, followSymlinks)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

type CopyTreeOptions struct {
	CopyFunction func(afero.Fs, string, string, bool) (string, error)
	Ignore       func(string, []os.FileInfo) []string
}

// Recursively copy a directory tree.
//
// The destination directory must not already exist.
//
// If the optional Symlinks flag is true, symbolic links in the
// source tree result in symbolic links in the destination tree; if
// it is false, the contents of the files pointed to by symbolic
// links are copied. If the file pointed by the symlink doesn't
// exist, an error will be returned.
//
// You can set the optional IgnoreDanglingSymlinks flag to true if you
// want to silence this error. Notice that this has no effect on
// platforms that don't support os.Symlink.
//
// The optional ignore argument is a callable. If given, it
// is called with the `src` parameter, which is the directory
// being visited by CopyTree(), and `names` which is the list of
// `src` contents, as returned by ioutil.ReadDir():
//
//   callable(src, entries) -> ignoredNames
//
// Since CopyTree() is called recursively, the callable will be
// called once for each directory that is copied. It returns a
// list of names relative to the `src` directory that should
// not be copied.
//
// The optional copyFunction argument is a callable that will be used
// to copy each file. It will be called with the source path and the
// destination path as arguments. By default, Copy() is used, but any
// function that supports the same signature (like Copy2() when it
// exists) can be used.
func CopyTree(fs afero.Fs, src, dst string, options *CopyTreeOptions) error {
	if options == nil {
		options = &CopyTreeOptions{
			Ignore:       nil,
			CopyFunction: Copy,
		}
	}

	srcFileInfo, err := fs.Stat(src)
	if err != nil {
		return err
	}

	if !srcFileInfo.IsDir() {
		return &NotADirectoryError{src}
	}

	_, err = fs.Open(dst)
	if !os.IsNotExist(err) {
		return &AlreadyExistsError{dst}
	}

	entries, err := afero.ReadDir(fs, src)
	if err != nil {
		return err
	}

	err = fs.MkdirAll(dst, srcFileInfo.Mode())
	if err != nil {
		return err
	}

	ignoredNames := []string{}
	if options.Ignore != nil {
		ignoredNames = options.Ignore(src, entries)
	}

	for _, entry := range entries {
		if stringInSlice(entry.Name(), ignoredNames) {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryFileInfo, err := fs.Stat(srcPath)
		if err != nil {
			return err
		}

		if entryFileInfo.IsDir() {
			err = CopyTree(fs, srcPath, dstPath, options)
			if err != nil {
				return err
			}
		} else {
			_, err = options.CopyFunction(fs, srcPath, dstPath, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
