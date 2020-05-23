package script

import (
	"fmt"
	"os"
)

// SetEnv sets a certain environment variable for this context
func (c *Context) SetEnv(key, value string) {
	c.env[key] = value
}

func (c *Context) GetCustomEnvValue(key string) string {
	return c.env[key]
}

func (c *Context) GetCustomEnv() []string {
	env := make([]string, 0)
	for key, value := range c.env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

func (c *Context) GetFullEnv() []string {
	env := os.Environ()
	env = append(env, c.GetCustomEnv()...)
	return env
}
