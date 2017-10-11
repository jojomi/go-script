package script

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkingDir(t *testing.T) {
	workingDir := "/tmp/working-dir"
	sc := NewContext()
	sc.SetWorkingDir(workingDir)
	assert.Equal(t, workingDir, sc.WorkingDir(), fmt.Sprintf("Expected working directory not set (should be %s)", workingDir))
}

func TestIsUserRoot(t *testing.T) {
	sc := NewContext()
	assert.False(t, sc.IsUserRoot())
}

func TestSetWorkingDirTemp(t *testing.T) {
	sc := NewContext()
	err := sc.SetWorkingDirTemp()
	assert.Nil(t, err)
	assert.Regexp(t, `/tmp/\d+`, sc.WorkingDir())
}
