package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintBold(t *testing.T) {
	sc := NewContext()
	sc.PrintDetectTTY = true

	sc.isTTY = false
	stdout, stderr := setOutputBuffers(sc)
	sc.PrintlnBold("in", "put")
	assert.Equal(t, "in put\n", stdout.String())
	assert.Equal(t, "", stderr.String())

	sc.isTTY = true
	stdout, stderr = setOutputBuffers(sc)
	sc.PrintlnBold("in", "put")
	assert.Equal(t, "\x1b[1min put\n\x1b[0m", stdout.String())
	assert.Equal(t, "", stderr.String())
}
