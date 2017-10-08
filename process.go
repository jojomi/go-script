// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// ProcessResult contains the results of a process execution be it successful or not.
type ProcessResult struct {
	Cmd          *exec.Cmd
	Process      *os.Process
	ProcessState *os.ProcessState
	ProcessError error
	stdoutBuffer *bytes.Buffer
	stderrBuffer *bytes.Buffer
}

// NewProcessResult creates a new empty ProcessResult
func NewProcessResult() *ProcessResult {
	p := &ProcessResult{}
	p.stdoutBuffer = bytes.NewBuffer(make([]byte, 0, 100))
	p.stderrBuffer = bytes.NewBuffer(make([]byte, 0, 100))
	return p
}

// Output returns a string representation of the output of the process denoted
// by this struct.
func (pr *ProcessResult) Output() string {
	return pr.stdoutBuffer.String()
}

// Error returns a string representation of the stderr output of the process denoted
// by this struct.
func (pr *ProcessResult) Error() string {
	return pr.stderrBuffer.String()
}

// Successful returns true iff the process denoted by this struct was run
// successfully. Success is defined as the exit code being set to 0.
func (pr *ProcessResult) Successful() bool {
	return pr.ExitCode() == 0
}

// StateString returns a string representation of the process denoted by
// this struct
func (pr *ProcessResult) StateString() string {
	state := pr.ProcessState
	return fmt.Sprintf("PID: %d, Exited: %t, Exit Code: %d, Success: %t, User Time: %s", state.Pid(), state.Exited(), pr.ExitCode(), state.Success(), state.UserTime())
}

// ExitCode returns the exit code of the command denoted by this struct
func (pr *ProcessResult) ExitCode() int {
	var waitStatus syscall.WaitStatus
	if exitError, ok := pr.ProcessError.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
	} else {
		waitStatus = pr.ProcessState.Sys().(syscall.WaitStatus)
	}
	return waitStatus.ExitStatus()
}

// CommandPath finds the full path of a binary given its name. This requires the wich command to be present in the system.
func (c *Context) CommandPath(name string) string {
	cmd := exec.Command("which", name)
	cmdOutput, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.Trim(string(cmdOutput), "\n")
}

// CommandExists checks if a given binary exists in PATH.
func (c *Context) CommandExists(name string) bool {
	return c.CommandPath(name) != ""
}

// MustCommandExist ensures a given binary exists in PATH, otherwise panics.
func (c *Context) MustCommandExist(name string) {
	if !c.CommandExists(name) {
		panic(fmt.Errorf("Command %s is not available. Please make sure it is installed and accessible.", name))
	}
}

// ExecuteDebug executes a system command, stdout and stderr are piped
func (c *Context) ExecuteDebug(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(false, false, name, args...)
	return
}

// ExecuteSilent executes a  system command without outputting stdout (it is
// still captured and can be retrieved using the returned ProcessResult)
func (c *Context) ExecuteSilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(true, false, name, args...)
	return
}

// ExecuteFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using the returned ProcessResult)
func (c *Context) ExecuteFullySilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(true, true, name, args...)
	return
}

// MustExecuteDebug ensures a system command to be executed, otherwise panics
func (c *Context) MustExecuteDebug(name string, args ...string) (pr *ProcessResult) {
	pr, err := c.Execute(false, false, name, args...)
	if err != nil {
		panic(err)
	}
	return
}

// MustExecuteSilent ensures a system command to be executed without outputting
// stdout, otherwise panics
func (c *Context) MustExecuteSilent(name string, args ...string) (pr *ProcessResult) {
	pr, err := c.ExecuteSilent(name, args...)
	if err != nil {
		panic(err)
	}
	return
}

// MustExecuteFullySilent ensures a system command to be executed without
// outputting stdout and stderr, otherwise panics
func (c *Context) MustExecuteFullySilent(name string, args ...string) (pr *ProcessResult) {
	pr, err := c.ExecuteFullySilent(name, args...)
	if err != nil {
		panic(err)
	}
	return
}

// Execute executes a system command with configurable stdout and stderr output
func (c *Context) Execute(stdoutSilent bool, stderrSilent bool, name string, args ...string) (pr *ProcessResult, err error) {
	cmd, pr := c.prepareCommand(stdoutSilent, stderrSilent, name, args...)

	err = cmd.Start()
	if err != nil {
		return
	}
	pr.Process = cmd.Process
	err = cmd.Wait()
	pr.ProcessState = cmd.ProcessState
	pr.ProcessError = err

	return
}

// ExecuteDetached executes the given command in this context in the background (detached). This means the script execution instantly continues.
func (c *Context) ExecuteDetached(name string, args ...string) (cmd *exec.Cmd, pr *ProcessResult, err error) {
	cmd, pr = c.prepareCommand(true, true, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	err = cmd.Start()
	pr.Process = cmd.Process
	return
}

func (c Context) prepareCommand(stdoutSilent bool, stderrSilent bool, name string, args ...string) (*exec.Cmd, *ProcessResult) {
	pr := NewProcessResult()

	cmd := exec.Command(name, args...)
	pr.Cmd = cmd

	cmd.Dir = c.workingDir
	cmd.Env = c.getFullEnv()

	if stdoutSilent {
		cmd.Stdout = pr.stdoutBuffer
	} else {
		cmd.Stdout = io.MultiWriter(c.stdout, pr.stdoutBuffer)
	}
	if stderrSilent {
		cmd.Stderr = pr.stderrBuffer
	} else {
		cmd.Stderr = io.MultiWriter(c.stderr, pr.stderrBuffer)
	}
	return cmd, pr
}
