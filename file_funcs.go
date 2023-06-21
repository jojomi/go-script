package script

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func (c Context) ReplaceInFile(filename, searchRegexp, replacement string) error {
	absoluteFilename := c.AbsPath(filename)

	// read file to string
	b, err := ioutil.ReadFile(absoluteFilename)
	if err != nil {
		return err
	}

	// replace
	re, err := regexp.Compile(searchRegexp)
	if err != nil {
		return err
	}
	b = re.ReplaceAll(b, []byte(replacement))

	// write back
	file, err := os.Open(absoluteFilename)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(absoluteFilename, b, fileInfo.Mode())
	if err != nil {
		return err
	}
	return nil
}

// FileHasContent func
func (c Context) FileHasContent(filename, search string) (bool, error) {
	fileContents, err := ioutil.ReadFile(c.AbsPath(filename))
	if err != nil {
		return false, err
	}
	return strings.Contains(string(fileContents), search), nil
}

// FileHasContentRegexp func
func (c Context) FileHasContentRegexp(filename, searchRegexp string) (bool, error) {
	fileContents, err := ioutil.ReadFile(c.AbsPath(filename))
	if err != nil {
		return false, err
	}
	r, err := regexp.Compile(searchRegexp)
	if err != nil {
		return false, err
	}
	results := r.FindStringIndex(string(fileContents))
	return len(results) > 0, nil
}

/*func (c Context) FileComment(filename, searchRegexp string) {

}

func (c Context) FileUncomment(filename, searchRegexp string) {

}*/
