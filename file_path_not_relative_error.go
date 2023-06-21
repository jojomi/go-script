package script

import "fmt"

type FilePathNotRelativeError struct {
	path string
}

func NewFilePathNotRelativeError(path string) *FilePathNotRelativeError {
	return &FilePathNotRelativeError{
		path: path,
	}
}

func (x FilePathNotRelativeError) Error() string {
	return fmt.Sprintf("file path was expected to be relative, but it was not: %s", x.path)
}
