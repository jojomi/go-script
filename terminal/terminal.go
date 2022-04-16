package terminal

import (
	"github.com/mattn/go-isatty"
	"os"
)

// IsInteractive returns if this program is run inside an interactive terminal
func IsInteractive() bool {
	return !(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())))
}
