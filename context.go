package script

import (
	"io"
	"os"

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
}

// NewContext returns a pointer to a new Context.
func NewContext() (context *Context) {
	// initialize Context
	context = &Context{}
	context.env = make(map[string]string, 0)
	context.fs = afero.NewOsFs()
	context.stdout = os.Stdout
	context.stderr = os.Stderr

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

// IsUserRoot checks if a user is root priviledged (Linux and Mac only? Windows?)
func (c *Context) IsUserRoot() bool {
	return os.Geteuid() == 0
}
