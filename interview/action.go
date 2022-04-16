package interview

import (
	"fmt"
	"regexp"
)

type Action struct {
	Label    string
	Shortcut string
}

func (x Action) WithShortcut(shortcut string) Action {
	return Action{
		Label:    x.Label,
		Shortcut: shortcut,
	}
}

func (x Action) ShortcutString() string {
	r := regexp.MustCompile(fmt.Sprintf(`(?i)^(.*?)(%s)`, x.Shortcut))
	return r.ReplaceAllString(x.Label, `$1[$2]`)
}

func (x Action) ShortcutStringAsDefault() string {
	r := regexp.MustCompile(fmt.Sprintf(`(?i)^(.*?)(%s)`, x.Shortcut))
	return r.ReplaceAllString(x.Label, `$1{$2}`)
}

func (x Action) String() string {
	return fmt.Sprintf("%s [shortcut %s]", x.Label, x.Shortcut)
}
