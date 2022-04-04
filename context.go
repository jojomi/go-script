package script

import (
	"fmt"
	"io"
	"os"
	"time"

	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/afero"
)

// Context for script operations. A Context includes the working directory and provides
// access the buffers and results of commands run in the Context.
// Using different Contexts it is possible to handle multiple separate environments.
type Context struct {
	workingDir string
	env        map[string]string
	fs         afero.Fs
	stdout     io.Writer
	stderr     io.Writer
	stdin      io.Reader
	isTTY      bool
	start      time.Time

	commandLogger *CommandLogger
	logKey        string
}

// NewContext returns a pointer to a new Context.
func NewContext() (context *Context) {
	// initialize Context
	context = &Context{
		env:    make(map[string]string, 0),
		fs:     afero.NewOsFs(),
		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  os.Stdin,

		start: time.Now(),
	}

	cwd, err := os.Getwd()
	if err == nil {
		context.SetWorkingDir(cwd)
	}
	return
}

// SetWorkingDir changes the current working dir
func (c *Context) SetWorkingDir(workingDir string) {
	c.workingDir = workingDir
}

// WorkingDir retrieves the current working dir
func (c *Context) WorkingDir() string {
	return c.workingDir
}

// SetWorkingDirTemp changes the current working dir to a temporary directory, returning an error in case something went wrong
func (c *Context) SetWorkingDirTemp() error {
	dir, err := c.TempDir()
	if err != nil {
		return err
	}
	c.SetWorkingDir(dir)
	return nil
}

// IsUserRoot checks if a user is root priviledged (Linux and Mac only? Windows?)
func (c *Context) IsUserRoot() bool {
	return os.Geteuid() == 0
}

// IsTerminal returns if this program is run inside an interactive terminal
func (c Context) IsTerminal() bool {
	return !(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())))
}

// SetCommandLogger sets a custom logger for commands being executed
func (c *Context) SetCommandLogger(commandLogger *CommandLogger) *Context {
	c.commandLogger = commandLogger
	return c
}

// SetLogKey defines the current key for logging.
func (c *Context) SetLogKey(logKey string) *Context {
	c.logKey = logKey
	return c
}

// LogKey returns the current key for logging.
func (c *Context) LogKey() string {
	return c.logKey
}

// LogCommand logs a command being executed. Automatically called by the Exec* function family.
func (c *Context) LogCommand(command Command) {
	if c.commandLogger == nil {
		return
	}

	l := *c.commandLogger
	err := l(*c, command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not log command execution: %s", err.Error())
	}
}

func (c *Context) GetStart() time.Time {
	return c.start
}
