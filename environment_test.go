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

	assert.False(t, inStringArray(sc.getFullEnv(), envKeyValue))
	sc.SetEnv(envKey, envValue)
	assert.True(t, inStringArray(sc.getFullEnv(), envKeyValue))
}

func inStringArray(haystack []string, needle string) bool {
	for _, f := range haystack {
		if f == needle {
			return true
		}
	}
	return false
}
