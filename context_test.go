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
	assert.Equal(t, workingDir, sc.WorkingDirPath(), fmt.Sprintf("Expected working directory not set (should be %s)", workingDir))
}

func TestIsUserRoot(t *testing.T) {
	sc := NewContext()
	assert.False(t, sc.IsUserRoot())
}

func TestSetWorkingDirTemp(t *testing.T) {
	assert := assert.New(t)

	sc := NewContext()
	err := sc.SetWorkingDirTemp()
	assert.Nil(err)
	wd1 := sc.WorkingDirPath()
	err = sc.SetWorkingDirTemp()
	assert.Nil(err)
	wd2 := sc.WorkingDirPath()

	assert.NotEqual(wd1, wd2)
}
