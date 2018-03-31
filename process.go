// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

// CommandConfig defines details of command execution.
type CommandConfig struct {
	RawStdout    bool
	RawStderr    bool
	OutputStdout bool
	OutputStderr bool
	ConnectStdin bool
	Detach       bool
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
	code, err := pr.ExitCode()
	if err != nil {
		return false
	}
	return code == 0
}

// StateString returns a string representation of the process denoted by
// this struct
func (pr *ProcessResult) StateString() string {
	state := pr.ProcessState
	exitCode, err := pr.ExitCode()
	exitCodeString := "?"
	if err == nil {
		exitCodeString = strconv.Itoa(exitCode)
	}
	return fmt.Sprintf("PID: %d, Exited: %t, Exit Code: %s, Success: %t, User Time: %s", state.Pid(), state.Exited(), exitCodeString, state.Success(), state.UserTime())
}

// ExitCode returns the exit code of the command denoted by this struct
func (pr *ProcessResult) ExitCode() (int, error) {
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

// CommandPath finds the full path of a binary given its name.
// also see https://golang.org/pkg/os/exec/#LookPath
func (c *Context) CommandPath(name string) (path string) {
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return
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

// ExecuteRaw executes a system command without touching stdout and stderr.
func (c *Context) ExecuteRaw(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		RawStdout:    true,
		RawStderr:    true,
		ConnectStdin: true,
	}, name, args...)
	return
}

// ExecuteDebug executes a system command, stdout and stderr are piped
func (c *Context) ExecuteDebug(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: true,
		OutputStderr: true,
		ConnectStdin: true,
	}, name, args...)
	return
}

// ExecuteSilent executes a  system command without outputting stdout (it is
// still captured and can be retrieved using the returned ProcessResult)
func (c *Context) ExecuteSilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: true,
		ConnectStdin: true,
	}, name, args...)
	return
}

// ExecuteFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using the returned ProcessResult)
func (c *Context) ExecuteFullySilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: false,
		ConnectStdin: true,
	}, name, args...)
	return
}

// MustExecuteDebug ensures a system command to be executed, otherwise panics
func (c *Context) MustExecuteDebug(name string, args ...string) (pr *ProcessResult) {
	pr, err := c.Execute(CommandConfig{
		OutputStdout: true,
		OutputStderr: true,
		ConnectStdin: true,
	}, name, args...)
	if err != nil {
		panic(err)
	}
	return
}

// MustExecuteSilent ensures a system command to be executed without outputting
// stdout, otherwise panics
func (c *Context) MustExecuteSilent(name string, args ...string) (pr *ProcessResult) {
	pr, err := c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: true,
		ConnectStdin: true,
	}, name, args...)
	if err != nil {
		panic(err)
	}
	return
}

// MustExecuteFullySilent ensures a system command to be executed without
// outputting stdout and stderr, otherwise panics
func (c *Context) MustExecuteFullySilent(name string, args ...string) (pr *ProcessResult) {
	pr, err := c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: false,
		ConnectStdin: true,
	}, name, args...)
	if err != nil {
		panic(err)
	}
	return
}

// Execute executes a system command according to given CommandConfig.
func (c *Context) Execute(cc CommandConfig, name string, args ...string) (pr *ProcessResult, err error) {
	cmd, pr := c.prepareCommand(cc, name, args...)

	if cc.Detach {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
	}

	err = cmd.Start()
	if err != nil {
		return
	}
	pr.Process = cmd.Process

	if !cc.Detach {
		c.WaitCmd(pr)
	}

	return
}

// ExecuteDetachedDebug executes a system command, stdout and stderr are piped.
// The command is executed in the background (detached).
func (c *Context) ExecuteDetachedDebug(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: true,
		OutputStderr: true,
		Detach:       true,
	}, name, args...)
	return
}

// ExecuteDetachedSilent executes a  system command without outputting stdout (it is
// still captured and can be retrieved using the returned ProcessResult).
// The command is executed in the background (detached).
func (c *Context) ExecuteDetachedSilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: true,
		Detach:       true,
	}, name, args...)
	return
}

// ExecuteDetachedFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using the returned ProcessResult).
// The command is executed in the background (detached).
func (c *Context) ExecuteDetachedFullySilent(name string, args ...string) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: false,
		Detach:       true,
	}, name, args...)
	return
}

func (c Context) prepareCommand(cc CommandConfig, name string, args ...string) (*exec.Cmd, *ProcessResult) {
	pr := NewProcessResult()

	cmd := exec.Command(name, args...)
	pr.Cmd = cmd

	cmd.Dir = c.workingDir
	cmd.Env = c.getFullEnv()

	if cc.RawStdout {
		cmd.Stdout = os.Stdout
	} else {
		if !cc.OutputStdout {
			cmd.Stdout = pr.stdoutBuffer
		} else {
			cmd.Stdout = io.MultiWriter(c.stdout, pr.stdoutBuffer)
		}
	}
	if cc.RawStderr {
		cmd.Stderr = os.Stderr
	} else {
		if !cc.OutputStderr {
			cmd.Stderr = pr.stderrBuffer
		} else {
			cmd.Stderr = io.MultiWriter(c.stderr, pr.stderrBuffer)
		}
	}

	if cc.ConnectStdin {
		cmd.Stdin = c.stdin
	}
	return cmd, pr
}

// WaitCmd waits for a command to be finished (useful on detached processes).
func (c Context) WaitCmd(pr *ProcessResult) {
	err := pr.Cmd.Wait()
	pr.ProcessState = pr.Cmd.ProcessState
	pr.ProcessError = err
}

// SplitCommand helper splits a string to command and arbitrarily many args.
// Does handle bash-like escaping (\) and string delimiters " and '.
func SplitCommand(input string) (command string, args []string) {
	quotes := []string{`"`, `'`}

	var (
		ok     bool
		length int
		value  string
		index  = 0
	)
	args = make([]string, 0)

outerloop:
	for {
		if index >= len(input) {
			break
		}

		ok, length, _ = parseWhitespace(input[index:])
		if ok {
			index += length
			continue
		}

		for _, quote := range quotes {
			ok, length, value = parseQuoted(input[index:], quote, `\`+quote)
			if ok {
				if command == "" {
					command = value
				} else {
					args = append(args, value)
				}
				index += length
				continue outerloop
			}
		}

		ok, length, value = parseUnquoted(input[index:])
		if ok {
			if command == "" {
				command = value
			} else {
				args = append(args, value)
			}
			index += length
			continue
		}
	}
	return
}

func parseQuoted(input, quoteString, escapeString string) (ok bool, length int, value string) {
	if !strings.HasPrefix(input, quoteString) {
		return
	}

	length = len(quoteString)
	for {
		if length >= len(input) {
			break
		}
		// escaped quoteString? (continue!)
		if strings.HasPrefix(input[length:], escapeString) {
			length += len(escapeString)
			value += quoteString
		}
		// quoteString (end!)
		if strings.HasPrefix(input[length:], quoteString) {
			length += len(quoteString)
			ok = true
			return
		}

		// otherwise inner content
		value += input[length : length+1]
		length++
	}

	return ok, length, value
}

func parseUnquoted(input string) (ok bool, length int, value string) {
	length = 0
	for {
		if length >= len(input) {
			ok = true
			return
		}
		// whitespace (end!) // TODO all whitespace!
		if strings.HasPrefix(input[length:], " ") {
			length++
			ok = true
			return
		}

		// otherwise inner content
		value += input[length : length+1]
		length++
	}
}

func parseWhitespace(input string) (ok bool, length int, value string) {
	length = 0
	for {
		if length >= len(input) {
			break
		}
		// no whitespace (end!) // TODO all whitespace!
		if !strings.HasPrefix(input[length:], " ") {
			ok = length > 0
			return
		}

		// otherwise inner content (whitespace)
		value += input[length : length+1]
		length++
	}

	return ok, length, value
}
