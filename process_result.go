package script

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

// MustCombinedOutput returns a string representation of all the output of the process denoted
// by this struct.
func (pr ProcessResult) MustCombinedOutput() string {
	out, err := pr.Cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return string(out)
}

// Output returns a string representation of the output of the process denoted
// by this struct.
func (pr ProcessResult) Output() string {
	return pr.stdoutBuffer.String()
}

// TrimmedOutput returns a string representation of the output of the process denoted
// by this struct with surrounding whitespace removed.
func (pr ProcessResult) TrimmedOutput() string {
	return strings.TrimSpace(pr.Output())
}

// TrimmedError returns a string representation of the error output of the process denoted
// by this struct with surrounding whitespace removed.
func (pr ProcessResult) TrimmedError() string {
	return strings.TrimSpace(pr.Error())
}

// Error returns a string representation of the stderr output of the process denoted
// by this struct.
func (pr ProcessResult) Error() string {
	return pr.stderrBuffer.String()
}

// Successful returns true iff the process denoted by this struct was run
// successfully. Success is defined as the exit code being set to 0.
func (pr *ProcessResult) Successful() bool {
	code, err := pr.ExitCode()
	if err != nil {
		return false
	}
	return code == 0
}

// StateString returns a string representation of the process denoted by
// this struct
func (pr ProcessResult) StateString() string {
	state := pr.ProcessState
	exitCode, err := pr.ExitCode()
	exitCodeString := "?"
	if err == nil {
		exitCodeString = strconv.Itoa(exitCode)
	}
	return fmt.Sprintf("PID: %d, Exited: %t, Exit Code: %s, Success: %t, User Time: %s", state.Pid(), state.Exited(), exitCodeString, state.Success(), state.UserTime())
}

// ExitCode returns the exit code of the command denoted by this struct
func (pr ProcessResult) ExitCode() (int, error) {
	var (
		waitStatus syscall.WaitStatus
		exitError  *exec.ExitError
	)
	ok := false
	if pr.ProcessError != nil {
		exitError, ok = pr.ProcessError.(*exec.ExitError)
	}
	if ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
	} else {
		if pr.ProcessState == nil {
			return -1, errors.New("no exit code available")
		}
		waitStatus = pr.ProcessState.Sys().(syscall.WaitStatus)
	}
	return waitStatus.ExitStatus(), nil
}
