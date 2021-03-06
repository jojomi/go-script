package script

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	commandInPath    = "ls"
	commandNotInPath = "binary-not-in-path"

	nonExistingBinary = "./not-existing"

	basicOutputStdout = "hello this is me\nwhatever\n"
	basicOutputStderr = "\"abc\"\n"
)

/* OUTPUT HANDLING */

func TestProcessOutputErrorCatching(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)

	pr, err := sc.ExecuteFullySilent(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Equal(t, basicOutputStdout, pr.Output())
	assert.Equal(t, basicOutputStderr, pr.Error())

	pr, err = sc.ExecuteDebug(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Equal(t, basicOutputStdout, pr.Output())
	assert.Equal(t, basicOutputStderr, pr.Error())

	pr, err = sc.ExecuteSilent(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Equal(t, basicOutputStdout, pr.Output())
	assert.Equal(t, basicOutputStderr, pr.Error())
}

func TestProcessStdoutStderr(t *testing.T) {
	sc := processContext()

	outBuffer, errBuffer := setOutputBuffers(sc)
	_, err := sc.ExecuteFullySilent(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Equal(t, "", outBuffer.String())
	assert.Equal(t, "", errBuffer.String())

	outBuffer, errBuffer = setOutputBuffers(sc)
	_, err = sc.ExecuteSilent(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Equal(t, "", outBuffer.String())
	assert.Equal(t, basicOutputStderr, errBuffer.String())

	outBuffer, errBuffer = setOutputBuffers(sc)
	_, err = sc.ExecuteDebug(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Equal(t, basicOutputStdout, outBuffer.String())
	assert.Equal(t, basicOutputStderr, errBuffer.String())
}

func TestProcessStdin(t *testing.T) {
	input := "my input"
	sc := processContext()

	sc.stdin = strings.NewReader(input)
	pr, err := sc.ExecuteFullySilent(LocalCommandFrom("./bin echo"))
	assert.Nil(t, err)
	assert.Equal(t, input+"\n", pr.Output())
	assert.Equal(t, input+"\n", pr.Error())
}

/* COMMAND EXECUTION */

func TestProcessRunFailure(t *testing.T) {
	sc := processContext()
	_, err := sc.ExecuteFullySilent(LocalCommandFrom(nonExistingBinary))
	assert.NotNil(t, err)
}

func TestProcessStateString(t *testing.T) {
	sc := processContext()
	pr, err := sc.ExecuteFullySilent(LocalCommandFrom("./bin basic-output"))
	assert.Nil(t, err)
	assert.Regexp(t, `^PID: \d+, Exited: true, Exit Code: 0, Success: true, User Time: \d+(\.\d+)?[mµ]?s$`, pr.StateString())
}

func TestProcessMustExecuteDebug(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)
	pr := sc.MustExecuteDebug(LocalCommandFrom("./bin basic-output"))
	assert.NotNil(t, pr)

	assert.Panics(t, func() {
		sc.MustExecuteDebug(LocalCommandFrom(nonExistingBinary))
	})
}

func TestProcessMustExecuteSilent(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)
	pr := sc.MustExecuteSilent(LocalCommandFrom("./bin basic-output"))
	assert.NotNil(t, pr)

	assert.Panics(t, func() {
		sc.MustExecuteSilent(LocalCommandFrom(nonExistingBinary))
	})
}

func TestProcessMustExecuteFullySilent(t *testing.T) {
	sc := processContext()
	setOutputBuffers(sc)
	pr := sc.MustExecuteFullySilent(LocalCommandFrom("./bin basic-output"))
	assert.NotNil(t, pr)

	assert.Panics(t, func() {
		sc.MustExecuteFullySilent(LocalCommandFrom(nonExistingBinary))
	})
}

/* DETACHED COMMANDS */

func TestProcessExecuteDetachedDebug(t *testing.T) {
	sc := processContext()
	stdout, stderr := setOutputBuffers(sc)
	pr, err := sc.ExecuteDetachedDebug(LocalCommandFrom("./bin sleep"))
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.IsType(t, int(0), pr.Process.Pid, "Not seen a PID on a detached process. Did it even start?") // int

	_, exitErr := pr.ExitCode()
	assert.NotNil(t, exitErr)
	assert.False(t, pr.Successful())

	sc.ExecuteDebug(LocalCommandFrom("./bin basic-output"))

	sc.WaitCmd(pr)
	assert.True(t, pr.Successful())
	assert.Equal(t, "before\n"+basicOutputStdout+"after\n", stdout.String())
	assert.Equal(t, "error-before\n"+basicOutputStderr+"error-after\n", stderr.String())
}

func TestProcessExecuteDetachedSilent(t *testing.T) {
	sc := processContext()
	stdout, stderr := setOutputBuffers(sc)
	pr, err := sc.ExecuteDetachedSilent(LocalCommandFrom("./bin sleep"))
	assert.Nil(t, err)
	assert.NotNil(t, pr)

	sc.ExecuteDebug(LocalCommandFrom("./bin basic-output"))

	sc.WaitCmd(pr)
	assert.Equal(t, basicOutputStdout, stdout.String())
	assert.Equal(t, "error-before\n"+basicOutputStderr+"error-after\n", stderr.String())
}

func TestProcessExecuteDetachedFullySilent(t *testing.T) {
	sc := processContext()
	stdout, stderr := setOutputBuffers(sc)
	pr, err := sc.ExecuteDetachedFullySilent(LocalCommandFrom("./bin sleep"))
	assert.Nil(t, err)
	assert.NotNil(t, pr)

	sc.ExecuteDebug(LocalCommandFrom("./bin basic-output"))

	sc.WaitCmd(pr)
	assert.Equal(t, basicOutputStdout, stdout.String())
	assert.Equal(t, basicOutputStderr, stderr.String())
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
	pr, err := sc.ExecuteFullySilent((LocalCommandFrom(commandInPath)))
	assert.Nil(t, err)
	assert.Equal(t, true, pr.Successful())
}

func TestProcessExitCode(t *testing.T) {
	sc := processContext()
	pr, err := sc.ExecuteFullySilent(LocalCommandFrom("./bin exit-code-error"))
	assert.Nil(t, err)
	exitCode, exitErr := pr.ExitCode()
	assert.Nil(t, exitErr)
	assert.Equal(t, 28, exitCode)
}

/* HELPER FUNCTIONS */

func processContext() *Context {
	sc := NewContext()
	sc.SetWorkingDir("./test/bin")
	return sc
}
