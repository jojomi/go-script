package script

import (
	"os"
	"path"
)

func (x File) Dir() Dir {
	return x.context.DirAt(path.Dir(x.path))
}

func (x File) EnsureDir(perm os.FileMode) error {
	dir := x.Dir()
	return dir.Ensure(perm)
}

func (x File) MustEnsureDir(perm os.FileMode) File {
	err := x.EnsureDir(perm)
	if err != nil {
		panic(err)
	}
	return x
}
