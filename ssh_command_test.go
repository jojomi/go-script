package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHCommandAdd(t *testing.T) {
	c := NewSSHCommand("host")
	assert.Equal(t, "ssh", c.Binary())
	c.Add("ls")
	assert.Equal(t, "ssh", c.Binary())
	c.Add("-l")
	assert.Equal(t, "ssh", c.Binary())
	assert.Equal(t, 3, len(c.Args()))
	assert.Equal(t, "-l", c.Args()[2])
}

func TestSSHCommandAddAll(t *testing.T) {
	c := NewSSHCommand("myhost")
	assert.Equal(t, "ssh", c.Binary())
	c.AddAll("remotecommand", "param1", "param2")
	assert.Equal(t, "ssh", c.Binary())
	assert.Equal(t, 4, len(c.Args()))
	assert.Equal(t, []string{"myhost", "remotecommand", "param1", "param2"}, c.Args())
}

func TestSSHCommandString(t *testing.T) {
	tests := []struct {
		Elements    []string
		ValidOutput string
	}{
		{
			Elements:    []string{"ls", "-al", "file"},
			ValidOutput: `ssh root@golang.org ls -al file`,
		},
		{
			Elements:    []string{"ls", `my file.txt`},
			ValidOutput: `ssh root@golang.org ls "my file.txt"`,
		},
		{
			Elements:    []string{"ls", `*.test`},
			ValidOutput: `ssh root@golang.org ls *.test`,
		},
		{
			Elements:    []string{"ls", `weird".file`},
			ValidOutput: `ssh root@golang.org ls weird\".file`,
		},
		{
			Elements:    []string{"ls", `'my custom file'`},
			ValidOutput: `ssh root@golang.org ls 'my custom file'`,
		},
	}

	for _, test := range tests {
		c := NewSSHCommand("root@golang.org")
		c.AddAll(test.Elements...)
		assert.Equal(t, test.ValidOutput, c.String())
	}
}

func TestSSHCommand_AddOpt(t *testing.T) {
	c := NewSSHCommand("host")
	c.AddOpt("ConnectTimeout", "1")
	c.Add("date")
	assert.Equal(t, "ssh -o ConnectTimeout=1 host date", c.String())
}
