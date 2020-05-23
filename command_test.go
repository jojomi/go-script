package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalCommandAdd(t *testing.T) {
	c := NewLocalCommand()
	assert.Equal(t, "", c.Binary())
	c.Add("ssh")
	assert.Equal(t, "ssh", c.Binary())
	c.Add("myhost")
	assert.Equal(t, "ssh", c.Binary())
	assert.Equal(t, 1, len(c.Args()))
	assert.Equal(t, "myhost", c.Args()[0])
}

func TestLocalCommandAddAll(t *testing.T) {
	c := NewLocalCommand()
	assert.Equal(t, "", c.Binary())
	c.AddAll("ssh", "myhost", "remotecommand")
	assert.Equal(t, "ssh", c.Binary())
	assert.Equal(t, 2, len(c.Args()))
	assert.Equal(t, []string{"myhost", "remotecommand"}, c.Args())
}

func TestLocalCommandString(t *testing.T) {
	tests := []struct {
		Elements    []string
		ValidOutput string
	}{
		{
			Elements:    []string{"ls", "-al", "file"},
			ValidOutput: `ls -al file`,
		},
		{
			Elements:    []string{"ls", `my file.txt`},
			ValidOutput: `ls "my file.txt"`,
		},
		{
			Elements:    []string{"ls", `*.test`},
			ValidOutput: `ls *.test`,
		},
		{
			Elements:    []string{"ls", `weird".file`},
			ValidOutput: `ls weird\".file`,
		},
		{
			Elements:    []string{"ls", `'my custom file'`},
			ValidOutput: `ls 'my custom file'`,
		},
	}

	for _, test := range tests {
		c := NewLocalCommand()
		c.AddAll(test.Elements...)
		assert.Equal(t, test.ValidOutput, c.String())
	}
}

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		input   string
		command string
		args    []string
	}{
		// simple cases
		{"ls -la", "ls", []string{"-la"}},
		{"./bin exit-code-error second_ARG", "./bin", []string{"exit-code-error", "second_ARG"}},
		// special cases
		{"", "", []string{}},
		// quoting
		{`"quoted bin" "fir st" 'sec ond'`, "quoted bin", []string{"fir st", "sec ond"}},
		{`bin -p  "fir st"   "sec ond"`, "bin", []string{"-p", "fir st", "sec ond"}},
		{`"\"bin" 'par am"'`, "\"bin", []string{"par am\""}},
	}

	for _, test := range tests {
		command, args := SplitCommand(test.input)
		assert.Equal(t, test.command, command)
		assert.Equal(t, test.args, args)
	}
}
