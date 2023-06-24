package script

import (
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// Dir is a filesystem directory.
type Dir struct {
	context *Context
	path    string
	fs      afero.Fs
}

func DirAt(path string) Dir {
	return NewContext().DirAt(path)
}

func (c *Context) DirAt(path string) Dir {
	// replace home dir path
	path, _ = homedir.Expand(path)

	// remove trailing path seperator if it exists
	path = strings.TrimRight(path, string(os.PathSeparator))

	return Dir{
		path:    path,
		context: c,
		fs:      afero.OsFs{},
	}
}

func (x Dir) Create(perm os.FileMode) error {
	// TODO implement
	return nil
}

func (x Dir) Ensure(perm os.FileMode) error {
	return x.Create(perm)
}

func (x Dir) Exists() bool {
	fi, err := x.context.fs.Stat(x.AbsPath())
	return !os.IsNotExist(err) && fi.IsDir()
}

func (x Dir) NotExists() bool {
	return !x.Exists()
}

func (x Dir) IsAbs() bool {
	return x.context.AbsPath(x.path) == x.path
}

func (x Dir) IsRel() bool {
	return !x.IsAbs()
}

func (x Dir) AbsPath() string {
	return x.context.AbsPath(x.path)
}

func (x Dir) RelPath() string {
	if x.IsRel() {
		return x.path
	}

	result, err := filepath.Rel(x.context.WorkingDir(), x.AbsPath())
	if err != nil {
		return ""
	}
	return result
}

func (x Dir) MustFileAt(relativePath string) File {
	file, err := x.FileAt(relativePath)
	if err != nil {
		panic(err)
	}
	return file
}

func (x Dir) FileAt(relativePath string) (File, error) {
	if x.context.FileAt(relativePath).IsAbs() {
		return File{}, NewFilePathNotRelativeError(relativePath)
	}

	return x.context.FileAt(filepath.Join(x.path, relativePath)), nil
}

func (x Dir) MustDirAt(relativePath string) Dir {
	dir, err := x.DirAt(relativePath)
	if err != nil {
		panic(err)
	}
	return dir
}

func (x Dir) DirAt(relativePath string) (Dir, error) {
	if x.context.DirAt(relativePath).IsAbs() {
		return Dir{}, NewDirPathNotRelativeError(relativePath)
	}

	return x.context.DirAt(filepath.Join(x.path, relativePath)), nil
}

func (x Dir) WithTrailingPathSeparator() string {
	return x.path + string(os.PathSeparator)
}

func (x Dir) WithoutTrailingPathSeparator() string {
	return x.path
}

func (x Dir) String() string {
	return x.path
}
