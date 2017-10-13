package script

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// FileExists checks if a given filename exists (being a file).
func (c *Context) FileExists(filename string) bool {
	filename = c.AbsPath(filename)
	fi, err := c.fs.Stat(filename)
	return !os.IsNotExist(err) && !fi.IsDir()
}

// EnsureDirExists ensures a directory with the given name exists.
// This function panics if it is unable to find or create a directory as requested.
// TODO also check if permissions are less than requested and update if possible
func (c *Context) EnsureDirExists(dirname string, perm os.FileMode) error {
	fullPath := c.AbsPath(dirname)
	if !c.DirExists(fullPath) {
		err := c.fs.MkdirAll(fullPath, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsurePathForFile guarantees the path for a given filename to exist.
// If the directory is not yet existing, it will be created using the permission
// mask given.
// TODO also check if permissions are less than requested and update if possible
func (c *Context) EnsurePathForFile(filename string, perm os.FileMode) error {
	return c.EnsureDirExists(filepath.Dir(filename), perm)
}

// DirExists checks if a given filename exists (being a directory).
func (c *Context) DirExists(path string) bool {
	path = c.AbsPath(path)
	fi, err := c.fs.Stat(path)
	return !os.IsNotExist(err) && fi.IsDir()
}

// TempFile returns a temporary file and an error if one occurred
func (c *Context) TempFile() (*os.File, error) {
	file, err := c.tempFileInternal()
	if err != nil {
		return nil, err
	}
	return file.(*os.File), nil
}

func (c *Context) tempFileInternal() (afero.File, error) {
	return afero.TempFile(c.fs, "", "")
}

// TempDir returns a temporary directory and an error if one occurred
func (c *Context) TempDir() (string, error) {
	return afero.TempDir(c.fs, "", "")
}

// AbsPath returns the absolute path of the path given. If the input path
// is absolute, it is returned untouched except for removing trailing path separators.
// Otherwise the absolute path is built relative to the current working directory of the Context.
// This function always returns a path *without* path separator at the end. See AbsPathSep for one that adds it.
func (c *Context) AbsPath(filename string) string {
	filename = c.WithoutTrailingPathSep(filename)
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

// AbsPathSep is a variant of AbsPath that always adds a trailing path separator
func (c *Context) AbsPathSep(filename string) string {
	return c.WithTrailingPathSep(c.AbsPath(filename))
}

// WithoutTrailingPathSep trims trailing os.PathSeparator from a string
func (c *Context) WithoutTrailingPathSep(input string) string {
	return strings.TrimRight(input, string(os.PathSeparator))
}

// WithTrailingPathSep adds a trailing os.PathSeparator to a string if it is missing
func (c *Context) WithTrailingPathSep(input string) string {
	if strings.HasSuffix(input, string(os.PathSeparator)) {
		return input
	}
	return input + string(os.PathSeparator)
}

// ResolveSymlinks resolve symlinks in a directory. All symlinked files are
// replaced with copies of the files they point to. Only one level symlinks
// are currently supported.
func (c *Context) ResolveSymlinks(dir string) error {
	var (
		err            error
		linkTargetPath string
		targetInfo     os.FileInfo
	)
	// directory does not exist -> nothing to do
	dir = c.AbsPath(dir)
	if !c.DirExists(dir) {
		return nil
	}
	err = afero.Walk(c.fs, dir, func(path string, info os.FileInfo, err error) error {
		// symlink?
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			// resolve
			linkTargetPath, err = filepath.EvalSymlinks(path)
			if err != nil {
				panic(err)
			}
			targetInfo, err = c.fs.Stat(linkTargetPath)
			if err != nil {
				panic(err)
			}
			c.fs.Remove(path)
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
	err := CopyFile(c.fs, from, to, true)
	if err != nil {
		return err
	}
	err = c.fs.Remove(from)
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
	options := &CopyTreeOptions{
		Ignore:       nil,
		CopyFunction: Copy,
	}
	err := CopyTree(c.fs, from, to, options)
	if err != nil {
		return err
	}
	err = c.fs.RemoveAll(from)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile copies a file. Cross-device copying is supported, so files
// can be copied from and to tmpfs mounts.
func (c *Context) CopyFile(from, to string) error {
	return CopyFile(c.fs, from, to, true) // don't follow symlinks
}

// CopyDir copies a directory. Cross-device copying is supported, so directories
// can be copied from and to tmpfs mounts.
func (c *Context) CopyDir(src, dst string) error {
	options := &CopyTreeOptions{
		Ignore:       nil,
		CopyFunction: Copy,
	}
	err := CopyTree(c.fs, src, dst, options)
	return err
}
