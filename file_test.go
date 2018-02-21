package script

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileHasContent(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("./test")

	_, err := sc.FileHasContent("non-existing-file.cfg", "github.com")
	assert.NotNil(t, err)

	has, err := sc.FileHasContent("file.cfg", "github.com")
	assert.Nil(t, err)
	assert.True(t, has)

	has, err = sc.FileHasContent("file.cfg", "gitlab.com")
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestFileHasContentRegexp(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("./test")

	_, err := sc.FileHasContentRegexp("non-existing-file.cfg", `git[a-z]*\.com`)
	assert.NotNil(t, err)

	_, err = sc.FileHasContentRegexp("file.cfg", `git\p[a-z]*\.com`) // invalid regexp
	assert.NotNil(t, err)

	has, err := sc.FileHasContentRegexp("file.cfg", `git[a-z]*\.com`)
	assert.Nil(t, err)
	assert.True(t, has)

	has, err = sc.FileHasContentRegexp("file.cfg", `gitlab\.com`)
	assert.Nil(t, err)
	assert.False(t, has)
}

func TestReplaceInFile(t *testing.T) {
	tests := []struct {
		content  string
		replace  string
		with     string
		expected string
	}{
		// string based (attention, always using RegExp syntax internally!)
		{
			content:  `useTLS: yes`,
			replace:  `useTLS: yes`,
			with:     `useTLS: no`,
			expected: `useTLS: no`,
		},
		// using backreferences
		{
			content:  `useTLS: yes`,
			replace:  `(useTLS): yes`,
			with:     `$1: no`,
			expected: `useTLS: no`,
		},
		// multiple replacements
		{
			content:  `I have many many more ideas.`,
			replace:  `many\s*`,
			with:     ``,
			expected: `I have more ideas.`,
		},
	}

	sc := NewContext()
	sc.SetWorkingDir(".")

	filename := "test/service.cfg"

	for _, test := range tests {
		makeFile(sc, filename, test.content)
		err := sc.ReplaceInFile(filename, test.replace, test.with)
		assert.Nil(t, err)
		output, _ := ioutil.ReadFile(sc.AbsPath(filename))
		assert.Equal(t, test.expected, string(output))
		os.Remove(filename)
	}
}
