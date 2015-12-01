// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/termie/go-shutil"
)

// FileExists checks if a given filename exists (being a file).
func (c *Context) FileExists(filename string) bool {
	filename = c.AbsPath(filename)
	fi, err := os.Stat(filename)
	return !os.IsNotExist(err) && !fi.IsDir()
}

// MustFileExist ensures if a given filename exists (being a file), panics otherwise.
func (c *Context) MustFileExist(filename string) {
	if !c.FileExists(filename) {
		panic(fmt.Errorf("File %s does not exist.", filename))
	}
}

// DirExists checks if a given filename exists (being a directory).
func (c *Context) DirExists(path string) bool {
	path = c.AbsPath(path)
	fi, err := os.Stat(path)
	return !os.IsNotExist(err) && fi.IsDir()
}

// MustDirExist checks if a given filename exists (being a directory).
func (c *Context) MustDirExist(path string) {
	if !c.DirExists(path) {
		panic(fmt.Errorf("Directory %s does not exist.", path))
	}
}

// MustGetTempFile guarantees to return a temporary file, panics otherwise
func (c *Context) MustGetTempFile() (tempFile *os.File) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	return
}

// MustGetTempDir guarantees to return a temporary directory, panics otherwise
func (c *Context) MustGetTempDir() (tempDir string) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return
}

// AbsPath returns the absolute path of the path given. If the input path
// is absolute, it is returned untouched. Otherwise the absolute path is
// built relative to the current working directory of the Context.
func (c *Context) AbsPath(filename string) string {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return filename
	}
	isAbsolute := absPath == filename
	if !isAbsolute {
		absPath, err := filepath.Abs(path.Join(c.workingDir, filename))
		if err != nil {
			return filename
		}
		return absPath
	}
	return filename
}

// ResolveSymlinks resolve symlinks in a directory. All symlinked files are
// replaced with copies of the files they point to. Only one level symlinks
// are currently supported.
func (c *Context) ResolveSymlinks(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// symlink?
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			// resolve
			linkTargetPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				panic(err)
			}
			targetInfo, err := os.Stat(linkTargetPath)
			if err != nil {
				panic(err)
			}
			os.Remove(path)
			// directory?
			if targetInfo.IsDir() {
				c.CopyDir(linkTargetPath, path)
			} else {
				c.CopyFile(linkTargetPath, path)
			}
		}
		return err
	})
	return err
}

/* Move/Copy Files and Directories */

// MoveFile moves a file. Cross-device moving is supported, so files
// can be moved from and to tmpfs mounts.
func (c *Context) MoveFile(from, to string) error {
	from = c.AbsPath(from)
	to = c.AbsPath(to)

	// work around "invalid cross-device link" for os.Rename
	err := shutil.CopyFile(from, to, true)
	if err != nil {
		return err
	}
	err = os.Remove(from)
	if err != nil {
		return err
	}
	return nil
}

// MoveDir moves a directory. Cross-device moving is supported, so directories
// can be moved from and to tmpfs mounts.
func (c *Context) MoveDir(from, to string) error {
	from = c.AbsPath(from)
	to = c.AbsPath(to)

	// work around "invalid cross-device link" for os.Rename
	options := &shutil.CopyTreeOptions{
		Symlinks:               true,
		Ignore:                 nil,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: false,
	}
	err := shutil.CopyTree(from, to, options)
	if err != nil {
		return err
	}
	err = os.RemoveAll(from)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile copies a file. Cross-device copying is supported, so files
// can be copied from and to tmpfs mounts.
func (c *Context) CopyFile(from, to string) error {
	return shutil.CopyFile(from, to, true) // don't follow symlinks
}

// CopyDir copies a directory. Cross-device copying is supported, so directories
// can be copied from and to tmpfs mounts.
func (c *Context) CopyDir(src, dst string) error {
	options := &shutil.CopyTreeOptions{
		Symlinks:               true,
		Ignore:                 nil,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: false,
	}
	err := shutil.CopyTree(src, dst, options)
	return err
}
