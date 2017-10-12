package script

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestPrintBold(t *testing.T) {
	sc := NewContext()

	color.NoColor = true
	stdout, stderr := setOutputBuffers(sc)
	sc.PrintlnBold("in", "put")
	assert.Equal(t, "in put\n", stdout.String())
	assert.Equal(t, "", stderr.String())

	color.NoColor = false
	stdout, stderr = setOutputBuffers(sc)
	sc.PrintlnBold("in", "put")
	sc.PrintBold("more")
	sc.PrintfBold("%dst", 1)
	assert.Equal(t, "\x1b[1min put\n\x1b[0m\x1b[1mmore\x1b[0m\x1b[1m1st\x1b[0m", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestPrintSuccess(t *testing.T) {
	sc := NewContext()

	color.NoColor = true
	stdout, stderr := setOutputBuffers(sc)
	sc.PrintlnSuccess("in", "put")
	assert.Equal(t, "in put\n", stdout.String())
	assert.Equal(t, "", stderr.String())

	color.NoColor = false
	stdout, stderr = setOutputBuffers(sc)
	sc.PrintlnSuccess("in", "put")
	sc.PrintSuccess("more")
	sc.PrintfSuccess("%dst", 1)
	assert.Equal(t, "\x1b[1;32min put\n\x1b[0m\x1b[1;32mmore\x1b[0m\x1b[1;32m1st\x1b[0m", stdout.String())
	assert.Equal(t, "", stderr.String())

	stdout, stderr = setOutputBuffers(sc)
	sc.PrintSuccessCheck(" mor", "e\n")
	sc.PrintSuccessCheck()
	assert.Equal(t, "\x1b[1;32m✓ more\n\x1b[0m\x1b[1;32m✓\x1b[0m", stdout.String())
	assert.Equal(t, "", stderr.String())
}

func TestPrintError(t *testing.T) {
	sc := NewContext()

	color.NoColor = true
	stdout, stderr := setOutputBuffers(sc)
	sc.PrintlnError("in", "put")
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "in put\n", stderr.String())

	color.NoColor = false
	stdout, stderr = setOutputBuffers(sc)
	sc.PrintlnError("in", "put")
	sc.PrintError("more")
	sc.PrintfError("%dst", 1)
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "\x1b[1;31min put\n\x1b[0m\x1b[1;31mmore\x1b[0m\x1b[1;31m1st\x1b[0m", stderr.String())

	stdout, stderr = setOutputBuffers(sc)
	sc.PrintErrorCross(" mo", "re\n")
	sc.PrintErrorCross()
	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "\x1b[1;31m✗ more\n\x1b[0m\x1b[1;31m✗\x1b[0m", stderr.String())
}
