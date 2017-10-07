package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessOutputErrorCatching(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("./test/bin")

	pr, err := sc.ExecuteFullySilent("bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", pr.Output())
	assert.Equal(t, "\"abc\"\n", pr.Error())

	pr, err = sc.ExecuteDebug("bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", pr.Output())
	assert.Equal(t, "\"abc\"\n", pr.Error())

	pr, err = sc.ExecuteSilent("bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", pr.Output())
	assert.Equal(t, "\"abc\"\n", pr.Error())
}
