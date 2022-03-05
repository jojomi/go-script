package terminal

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsInteractive returns if this program is run inside an interactive terminal
func IsInteractive() bool {
	return !(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())))
}

var NonInteractiveTerminalErrorValue NonInteractiveTerminalError

type NonInteractiveTerminalError struct{}

func (x *NonInteractiveTerminalError) Error() string {
	return "non-interactive terminal"
}

func IsNonInteractiveTerminalError(err error) bool {
	return err == &NonInteractiveTerminalErrorValue
}

func SetNonInteractiveDefault[T any](err error, value, defaultValue T) T {
	if !IsNonInteractiveTerminalError(err) {
		return value
	}
	return defaultValue
}
