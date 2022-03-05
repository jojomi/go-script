package script

import (
	"strings"
)

type LocalCommand struct {
	elements []string
}

func NewLocalCommand() *LocalCommand {
	l := LocalCommand{}
	return &l
}

func LocalCommandFrom(command string) *LocalCommand {
	c, args := SplitCommand(command)
	l := NewLocalCommand()
	l.Add(c)
	l.AddAll(args...)
	return l
}

func (l *LocalCommand) AddAll(input ...string) {
	for _, i := range input {
		l.Add(i)
	}
}

func (l *LocalCommand) Add(input string) {
	if l.elements == nil {
		l.elements = make([]string, 0)
	}
	l.elements = append(l.elements, input)
}

func (l *LocalCommand) Binary() string {
	if l.elements == nil || len(l.elements) == 0 {
		return ""
	}
	return l.elements[0]
}

func (l *LocalCommand) Args() []string {
	if l.elements == nil || len(l.elements) < 2 {
		return []string{}
	}
	return l.elements[1:]
}

func (l *LocalCommand) String() string {
	var b strings.Builder
	for i, e := range l.elements {
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
