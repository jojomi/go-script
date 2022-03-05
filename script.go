package script

import (
	"fmt"
	"os"
)

// RecoverFunc prints any panic message to Stderr
var RecoverFunc = func() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "%v\n", r)
		os.Exit(1)
	}
}
