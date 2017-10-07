package script

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessOutputErrorCatching(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("./test/bin")
	setOutputBuffers(sc)

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

func TestProcessStdoutStderr(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("./test/bin")

	outBuffer, errBuffer := setOutputBuffers(sc)
	_, err := sc.ExecuteFullySilent("bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "", outBuffer.String())
	assert.Equal(t, "", errBuffer.String())

	outBuffer, errBuffer = setOutputBuffers(sc)
	_, err = sc.ExecuteSilent("bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "", outBuffer.String())
	assert.Equal(t, "\"abc\"\n", errBuffer.String())

	outBuffer, errBuffer = setOutputBuffers(sc)
	_, err = sc.ExecuteDebug("bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", outBuffer.String())
	assert.Equal(t, "\"abc\"\n", errBuffer.String())
}

func TestProcessCommandExists(t *testing.T) {
	sc := NewContext()
	sc.CommandExists("ls")
}

func TestProcessSuccessful(t *testing.T) {
	sc := NewContext()
	pr, err := sc.ExecuteFullySilent("ls")
	assert.Nil(t, err)
	assert.Equal(t, true, pr.Successful())
}

func setOutputBuffers(sc *Context) (out, err *bytes.Buffer) {
	stdoutBuffer := bytes.NewBuffer(make([]byte, 0, 100))
	sc.stdout = stdoutBuffer
	stderrBuffer := bytes.NewBuffer(make([]byte, 0, 100))
	sc.stderr = stderrBuffer
	return stdoutBuffer, stderrBuffer
}
