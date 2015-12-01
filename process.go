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
		panic(fmt.Errorf("Command %s is not available. Please make sure it is installed and accessible. Ouptut of which: %s", name, c.LastOutput()))
	}
}

// ExecuteNormal execute a system command, stdout and stderr are output just
// normally
func (c *Context) ExecuteDebug(name string, args ...string) (err error) {
	err = c.Execute(false, false, name, args...)
	return
}

// ExecuteSilent executes a  system command without outputting stdout (it is
// still captured and can be retrieved using LastOutput())
func (c *Context) ExecuteSilent(name string, args ...string) (err error) {
	err = c.Execute(true, false, name, args...)
	return
}

// ExecuteFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using LastOutput() and
// LastError())
func (c *Context) ExecuteFullySilent(name string, args ...string) (err error) {
	err = c.Execute(true, true, name, args...)
	return
}

// MustExecute ensures a system command to be executed, otherwise panics
func (c *Context) MustExecuteDebug(name string, args ...string) {
	err := c.Execute(false, false, name, args...)
	if err != nil {
		panic(err)
	}
}

// MustExecuteSilent ensures a system command to be executed without outputting
// stdout, otherwise panics
func (c *Context) MustExecuteSilent(name string, args ...string) {
	err := c.ExecuteSilent(name, args...)
	if err != nil {
		panic(err)
	}
}

// MustExecuteFullySilent ensures a system command to be executed without
// outputting stdout and stderr, otherwise panics
func (c *Context) MustExecuteFullySilent(name string, args ...string) {
	err := c.ExecuteFullySilent(name, args...)
	if err != nil {
		panic(err)
	}
}

// Execute exceutes a system command with configurable stdout and stderr output
// https://github.com/golang/go/issues/9307
func (c *Context) Execute(stdoutSilent bool, stderrSilent bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Dir = c.workingDir

	// handling Stdout and Stderr
	// idea from http://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	var wg sync.WaitGroup
	wg.Add(2)

	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	outScanner := bufio.NewScanner(cmdOutReader)
	c.stdoutBuffer = new(bytes.Buffer)
	go func() {
		defer wg.Done()
		outputHandler(outScanner, !stdoutSilent, c.stdoutBuffer)
	}()

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	errScanner := bufio.NewScanner(cmdErrReader)
	c.stderrBuffer = new(bytes.Buffer)
	go func() {
		defer wg.Done()
		outputHandler(errScanner, !stderrSilent, c.stderrBuffer)
	}()
	cmd.Start()

	err = cmd.Wait()
	// make sure all output is captured and processed before continuing
	wg.Wait()

	c.lastProcessState = cmd.ProcessState
	return err
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

// LastOutput returns the output buffer of the last command executed using one
// of the Execute* functions. stdout is captured for any command run.
func (c *Context) LastOutput() string {
	return c.stdoutBuffer.String()
}

// LastError returns the error buffer of the last command executed using one
// of the Execute* functions. stderr is captured for any command run.
func (c *Context) LastError() string {
	return c.stderrBuffer.String()
}

// LastSuccessful determins if the last command run was successful. Success is
// defined as the process' return code being zero.
func (c *Context) LastSuccessful() bool {
	return c.LastExitCode() == 0
}

// LastExitCode returns the exit code of the last command run.
func (c *Context) LastExitCode() int {
	var waitStatus syscall.WaitStatus
	if exitError, ok := c.lastProcessError.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
	} else {
		waitStatus = c.lastProcessState.Sys().(syscall.WaitStatus)
	}
	fmt.Println(waitStatus)
	return waitStatus.ExitStatus()
}

// LastProcessState returns the *os.ProcessState of the last command run.
func (c *Context) LastProcessState() *os.ProcessState {
	return c.lastProcessState
}

// PrintLastState conveniently prints statistics on the last command run. Useful
// for debugging purposes.
func (c *Context) PrintLastState() {
	state := c.lastProcessState
	fmt.Println(
		"PID:", state.Pid(),
		"Exited:", state.Exited(),
		"Exit Code:", c.LastExitCode(),
		"Success:", state.Success(),
		"User Time:", state.UserTime())
}
