// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// ProcessResult contains the results of a process execution be it successful or not.
type ProcessResult struct {
	Cmd          *exec.Cmd
	ProcessState *os.ProcessState
	ProcessError error
	stdoutBuffer *bytes.Buffer
	stderrBuffer *bytes.Buffer
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
	fmt.Println(pr.ExitCode())
	return pr.ExitCode() == 0
}

// StateString returns a string representation of the process denoted by
// this struct
func (pr *ProcessResult) StateString() string {
	state := pr.ProcessState
	return fmt.Sprintf("PID: %q, Exited: %t, Exit Code: %q, Success: %t, User Time: %q", state.Pid(), state.Exited(), pr.ExitCode(), state.Success(), state.UserTime())
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

// CommandPath finds the full path of a binary given its name.
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
// still captured and can be retrieved using LastOutput())
func (c *Context) ExecuteSilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(true, false, name, args...)
	return
}

// ExecuteFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using LastOutput() and
// LastError())
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

// Execute exceutes a system command with configurable stdout and stderr output
// https://github.com/golang/go/issues/9307
func (c *Context) Execute(stdoutSilent bool, stderrSilent bool, name string, args ...string) (pr *ProcessResult, err error) {
	pr = &ProcessResult{}

	cmd := exec.Command(name, args...)

	cmd.Dir = c.workingDir

	// handling Stdout and Stderr
	// idea from http://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	var wg sync.WaitGroup

	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	outScanner := bufio.NewScanner(cmdOutReader)
	pr.stdoutBuffer = new(bytes.Buffer)
	wg.Add(1)
	go func() {
		defer wg.Done()
		outputHandler(outScanner, !stdoutSilent, pr.stdoutBuffer)
	}()

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	errScanner := bufio.NewScanner(cmdErrReader)
	pr.stderrBuffer = new(bytes.Buffer)
	wg.Add(1)
	go func() {
		defer wg.Done()
		outputHandler(errScanner, !stderrSilent, pr.stderrBuffer)
	}()
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	// make sure all output is captured and processed before continuing
	wg.Wait()

	pr.Cmd = cmd
	pr.ProcessState = cmd.ProcessState
	pr.ProcessError = err

	return pr, err
}

// internal
func outputHandler(scanner *bufio.Scanner, output bool, buffer *bytes.Buffer) {
	for scanner.Scan() {
		text := scanner.Text()
		if buffer.Len() > 0 {
			buffer.WriteString("\n")
		}
		_, err := buffer.WriteString(text)
		if err != nil {
			panic(err)
		}
		if output {
			fmt.Println(text, buffer.String())
		}
	}
}
