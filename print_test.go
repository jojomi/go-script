package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintBold(t *testing.T) {
	sc := NewContext()
	stdout, stderr := setOutputBuffers(sc)

	sc.PrintlnBold("in", "put")
	assert.Equal(t, "in put\n", stdout.String())
	assert.Equal(t, "", stderr.String())
}
