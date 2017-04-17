package script

import (
	"fmt"
	"os"
)

// SetEnv sets a certain environment variable for this context
func (c *Context) SetEnv(key, value string) {
	c.env[key] = value
}

func (c *Context) getFullEnv() []string {
	env := os.Environ()
	for key, value := range c.env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
