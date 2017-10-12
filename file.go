package script

import (
	"io/ioutil"
	"regexp"
	"strings"
)

/*func (c Context) ReplaceInFile(filename, searchRegexp, replacement string) error {
	return nil
}*/

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
