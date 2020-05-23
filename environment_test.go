package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	envKey := "MY_ENV_KEY"
	envValue := "environment value"
	envKeyValue := envKey + "=" + envValue
	sc := NewContext()

	assert.False(t, inStringArray(sc.GetFullEnv(), envKeyValue))
	assert.Empty(t, sc.GetCustomEnvValue(envKey))
	assert.Empty(t, sc.GetCustomEnv())
	sc.SetEnv(envKey, envValue)
	assert.True(t, inStringArray(sc.GetCustomEnv(), envKeyValue))
	assert.True(t, inStringArray(sc.GetFullEnv(), envKeyValue))
}

func inStringArray(haystack []string, needle string) bool {
	for _, f := range haystack {
		if f == needle {
			return true
		}
	}
	return false
}
