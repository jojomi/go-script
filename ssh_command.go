package script

import (
	"fmt"
	"strings"
)

type SSHCommand struct {
	target     string
	sshOptions map[string]string
	elements   []string
}

func NewSSHCommand(target string) *SSHCommand {
	l := &SSHCommand{
		target: target,
	}
	return l
}

func SSHCommandFrom(target, command string) *SSHCommand {
	l := NewSSHCommand(target)
	c, args := SplitCommand(command)
	l.Add(c)
	l.AddAll(args...)
	return l
}

func (l *SSHCommand) AddAll(input ...string) {
	for _, i := range input {
		l.Add(i)
	}
}

func (l *SSHCommand) Add(input string) {
	if l.elements == nil {
		l.elements = make([]string, 0)
	}
	l.elements = append(l.elements, input)
}

func (l *SSHCommand) AddOpt(opt, value string) {
	if l.sshOptions == nil {
		l.sshOptions = make(map[string]string, 0)
	}
	l.sshOptions[opt] = value
}

func (l *SSHCommand) Binary() string {
	return l.allElements()[0]
}

func (l *SSHCommand) Args() []string {
	return l.allElements()[1:]
}

func (l *SSHCommand) allElements() []string {
	result := []string{"ssh"}
	for opt, val := range l.sshOptions {
		result = append(result, "-o", fmt.Sprintf("%s=%s", opt, val))
	}
	result = append(result, l.target)
	result = append(result, l.elements...)
	return result
}

func (l *SSHCommand) String() string {
	var b strings.Builder
	for i, e := range l.allElements() {
		if i > 0 {
			b.WriteString(" ")
		}
		// contains double quotes? escape them!
		if strings.Contains(e, `"`) {
			e = strings.ReplaceAll(e, `"`, `\"`)
		}
		// contains Whitespace? wrap with double quotes
		if strings.Contains(e, ` `) {
			if !isWrapped(e, `"`) && !isWrapped(e, `'`) {
				e = `"` + e + `"`
			}
		}
		b.WriteString(e)
	}
	return b.String()
}
