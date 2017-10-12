package script

import (
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
