package script

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// File is a file in the filesystem.
type File struct {
	context           *Context
	path              string
	createPermissions os.FileMode
	fs                afero.Fs
}

func FileAt(path string) File {
	return NewContext().FileAt(path)
}

func (c *Context) FileAt(path string) File {
	// replace home dir path
	path, _ = homedir.Expand(path)

	return File{
		path:              path,
		context:           c,
		createPermissions: 0640,
		fs:                afero.OsFs{},
	}
}

// CreatePermissions allows you to define the FileMode used when creating this file (if it did not exist).
func (x File) CreatePermissions(perm os.FileMode) File {
	x.createPermissions = perm
	return x
}

func (x File) Exists() bool {
	fi, err := x.context.fs.Stat(x.AbsPath())
	return !os.IsNotExist(err) && !fi.IsDir()
}

func (x File) AssertExists() File {
	if !x.Exists() {
		panic(fmt.Errorf("file %s should have existed", x))
	}
	return x
}

func (x File) NotExists() bool {
	return !x.Exists()
}

func (x File) AssertNotExists() File {
	if !x.NotExists() {
		panic(fmt.Errorf("file %s should not have existed", x))
	}
	return x
}

func (x File) IsAbs() bool {
	return x.context.AbsPath(x.path) == x.path
}

func (x File) IsRel() bool {
	return !x.IsAbs()
}

func (x File) AbsPath() string {
	return x.context.AbsPath(x.path)
}

func (x File) RelPath() string {
	result, err := filepath.Rel(x.context.WorkingDir(), x.AbsPath())
	if err != nil {
		return x.AbsPath()
	}
	return result
}

func (x File) IsHidden() bool {
	return strings.HasPrefix(x.Filename(), ".")
}

func (x File) Equals(otherFile File) bool {
	return x.AbsPath() == otherFile.AbsPath()
}

func (x File) String() string {
	return x.path
}

func (x File) Filename() string {
	return path.Base(x.path)
}
