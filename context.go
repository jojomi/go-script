// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"os"
)

// Context for script operations. A Context includes the working directory and provides
// access the buffers and results of commands run in the Context.
// Using different Contexts it is possible to handle multiple separate environments.
type Context struct {
	workingDir string
	env        []string
}

// NewContext returns a pointer to a new Context.
func NewContext() (context *Context) {
	// initialize Context
	context = &Context{}

	cwd, err := os.Getwd()
	if err == nil {
		context.SetWorkingDir(cwd)
	}
	return
}

// SetWorkingDir changes the current working dir.
func (c *Context) SetWorkingDir(workingDir string) {
	c.workingDir = workingDir
}

// WorkingDir retrieves the current working dir.
func (c *Context) WorkingDir() string {
	return c.workingDir
}
