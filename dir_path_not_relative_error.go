package script

import "fmt"

type DirPathNotRelativeError struct {
	path string
}

func NewDirPathNotRelativeError(path string) *DirPathNotRelativeError {
	return &DirPathNotRelativeError{
		path: path,
	}
}

func (x DirPathNotRelativeError) Error() string {
	return fmt.Sprintf("dir path was expected to be relative, but it was not: %s", x.path)
}
