package script

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	commandInPath    = "ls"
	commandNotInPath = "binary-not-in-path"

	nonExistingBinary = "./not-existing"
)

/* OUTPUT HANDLING */

func TestProcessOutputErrorCatching(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)

	pr, err := sc.ExecuteFullySilent("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", pr.Output())
	assert.Equal(t, "\"abc\"\n", pr.Error())

	pr, err = sc.ExecuteDebug("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", pr.Output())
	assert.Equal(t, "\"abc\"\n", pr.Error())

	pr, err = sc.ExecuteSilent("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", pr.Output())
	assert.Equal(t, "\"abc\"\n", pr.Error())
}

func TestProcessStdoutStderr(t *testing.T) {
	sc := processContext()

	outBuffer, errBuffer := setOutputBuffers(sc)
	_, err := sc.ExecuteFullySilent("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "", outBuffer.String())
	assert.Equal(t, "", errBuffer.String())

	outBuffer, errBuffer = setOutputBuffers(sc)
	_, err = sc.ExecuteSilent("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "", outBuffer.String())
	assert.Equal(t, "\"abc\"\n", errBuffer.String())

	outBuffer, errBuffer = setOutputBuffers(sc)
	_, err = sc.ExecuteDebug("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Equal(t, "hello this is me\nwhatever\n", outBuffer.String())
	assert.Equal(t, "\"abc\"\n", errBuffer.String())
}

/* COMMAND EXECUTION */

func TestProcessRunFailure(t *testing.T) {
	sc := processContext()
	_, err := sc.ExecuteFullySilent(nonExistingBinary)
	assert.NotNil(t, err)
}

func TestProcessStateString(t *testing.T) {
	sc := processContext()
	pr, err := sc.ExecuteFullySilent("./bin", "basic-output")
	assert.Nil(t, err)
	assert.Regexp(t, `^PID: \d+, Exited: true, Exit Code: 0, Success: true, User Time: \d+(\.\d+)?[mµ]?s$`, pr.StateString())
}

func TestProcessMustExecuteDebug(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)
	pr := sc.MustExecuteDebug("./bin", "basic-output")
	assert.NotNil(t, pr)

	assert.Panics(t, func() {
		sc.MustExecuteDebug(nonExistingBinary)
	})
}

func TestProcessMustExecuteSilent(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)
	pr := sc.MustExecuteSilent("./bin", "basic-output")
	assert.NotNil(t, pr)

	assert.Panics(t, func() {
		sc.MustExecuteSilent(nonExistingBinary)
	})
}

func TestProcessMustExecuteFullySilent(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)
	pr := sc.MustExecuteFullySilent("./bin", "basic-output")
	assert.NotNil(t, pr)

	assert.Panics(t, func() {
		sc.MustExecuteFullySilent(nonExistingBinary)
	})
}

/* COMMAND HANDLING */

func TestProcessCommandExists(t *testing.T) {
	sc := NewContext()
	sc.CommandExists(commandInPath)
}

func TestProcessCommandPath(t *testing.T) {
	sc := NewContext()
	path := sc.CommandPath(commandInPath)
	assert.NotEqual(t, "", path)
}

func TestProcessCommandPathFailure(t *testing.T) {
	sc := NewContext()
	path := sc.CommandPath(commandNotInPath)
	assert.Equal(t, "", path)
}

func TestProcessCommandPathFailurePanic(t *testing.T) {
	sc := NewContext()
	assert.Panics(t, func() {
		sc.MustCommandExist(commandNotInPath)
	})
}

/* EXIT CODE HANDLING */

func TestProcessSuccessful(t *testing.T) {
	sc := processContext()
	pr, err := sc.ExecuteFullySilent(commandInPath)
	assert.Nil(t, err)
	assert.Equal(t, true, pr.Successful())
}

func TestProcessExitCode(t *testing.T) {
	sc := processContext()
	pr, err := sc.ExecuteFullySilent("./bin", "exit-code-error")
	assert.NotNil(t, err)
	assert.Equal(t, 28, pr.ExitCode())
}

/* HELPER FUNCTIONS */

func processContext() *Context {
	sc := NewContext()
	sc.SetWorkingDir("./test/bin")
	return sc
}

func setOutputBuffers(sc *Context) (out, err *bytes.Buffer) {
	stdoutBuffer := bytes.NewBuffer(make([]byte, 0, 100))
	sc.stdout = stdoutBuffer
	stderrBuffer := bytes.NewBuffer(make([]byte, 0, 100))
	sc.stderr = stderrBuffer
	return stdoutBuffer, stderrBuffer
}
