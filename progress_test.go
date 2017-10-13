package script

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivityIndicator(t *testing.T) {
	sc := NewContext()
	stdout, stderr := setOutputBuffers(sc)

	ai := sc.ActivityIndicator("indicator")
	ai.Start()
	time.Sleep(1 * time.Second)
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())
}
