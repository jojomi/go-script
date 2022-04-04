// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// CommandConfig defines details of command execution.
type CommandConfig struct {
	RawStdout    bool
	RawStderr    bool
	OutputStdout bool
	OutputStderr bool
	ConnectStdin bool
	Detach       bool
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
func (c *Context) ExecuteRaw(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		RawStdout:    true,
		RawStderr:    true,
		ConnectStdin: true,
	}, command)
	return
}

// ExecuteDebug executes a system command, stdout and stderr are piped
func (c *Context) ExecuteDebug(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: true,
		OutputStderr: true,
		ConnectStdin: true,
	}, command)
	return
}

// ExecuteSilent executes a  system command without outputting stdout (it is
// still captured and can be retrieved using the returned ProcessResult)
func (c *Context) ExecuteSilent(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: true,
		ConnectStdin: true,
	}, command)
	return
}

// ExecuteFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using the returned ProcessResult)
func (c *Context) ExecuteFullySilent(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: false,
		ConnectStdin: true,
	}, command)
	return
}

// MustExecuteDebug ensures a system command to be executed, otherwise panics
func (c *Context) MustExecuteDebug(command Command) (pr *ProcessResult) {
	pr, err := c.Execute(CommandConfig{
		OutputStdout: true,
		OutputStderr: true,
		ConnectStdin: true,
	}, command)
	if err != nil {
		panic(err)
	}
	return
}

// MustExecuteSilent ensures a system command to be executed without outputting
// stdout, otherwise panics
func (c *Context) MustExecuteSilent(command Command) (pr *ProcessResult) {
	pr, err := c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: true,
		ConnectStdin: true,
	}, command)
	if err != nil {
		panic(err)
	}
	return
}

// MustExecuteFullySilent ensures a system command to be executed without
// outputting stdout and stderr, otherwise panics
func (c *Context) MustExecuteFullySilent(command Command) (pr *ProcessResult) {
	pr, err := c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: false,
		ConnectStdin: true,
	}, command)
	if err != nil {
		panic(err)
	}
	return
}

// Execute executes a system command according to given CommandConfig.
func (c *Context) Execute(cc CommandConfig, command Command) (pr *ProcessResult, err error) {
	cmd, pr := c.prepareCommand(cc, command)

	if cc.Detach {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
	}

	// logging
	c.LogCommand(command)

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
func (c *Context) ExecuteDetachedDebug(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: true,
		OutputStderr: true,
		Detach:       true,
	}, command)
	return
}

// ExecuteDetachedSilent executes a system command without outputting stdout (it is
// still captured and can be retrieved using the returned ProcessResult).
// The command is executed in the background (detached).
func (c *Context) ExecuteDetachedSilent(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: true,
		Detach:       true,
	}, command)
	return
}

// ExecuteDetachedFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using the returned ProcessResult).
// The command is executed in the background (detached).
func (c *Context) ExecuteDetachedFullySilent(command Command) (pr *ProcessResult, err error) {
	pr, err = c.Execute(CommandConfig{
		OutputStdout: false,
		OutputStderr: false,
		Detach:       true,
	}, command)
	return
}

func (c Context) prepareCommand(cc CommandConfig, command Command) (*exec.Cmd, *ProcessResult) {
	pr := NewProcessResult()

	cmd := exec.Command(command.Binary(), command.Args()...)
	pr.Cmd = cmd

	cmd.Dir = c.workingDir
	cmd.Env = c.GetFullEnv()

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

func Errify(pr *ProcessResult, err error) error {
	if err != nil {
		return err
	}

	if pr.Successful() {
		return nil
	}
	return fmt.Errorf("command execution failed: %s", pr.Error())
}
