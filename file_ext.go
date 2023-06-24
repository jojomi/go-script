package script

import (
	"fmt"
	"regexp"
	"strings"
)

func (x File) HasAnyExtension() bool {
	filename := x.Filename()
	return strings.Contains(filename[1:], ".")
}

func (x File) HasExtension(fileExtension FileExtension) bool {
	return strings.HasSuffix(x.path, fileExtension.WithDot())
}

func (x File) AssertExtension(fileExtension FileExtension) File {
	if !x.HasExtension(fileExtension) {
		panic(fmt.Errorf("file %s should have had file extension %s", x, fileExtension))
	}
	return x
}

func (x File) WithExtension(fileExtension FileExtension) File {
	if !x.HasAnyExtension() {
		return x.context.FileAt(x.path + fileExtension.WithDot())
	}

	r := regexp.MustCompile(`(\.tar)?\.[^.]+$`)
	newFilename := r.ReplaceAllString(x.Filename(), fileExtension.WithDot())

	return x.context.FileAt(newFilename)
}
