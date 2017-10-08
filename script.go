package script

import (
	"fmt"
	"os"
)

// RecoverFunc prints any panic message to Stderr
var RecoverFunc = func() {
	if r := recover(); r != nil {
		os.Stderr.WriteString(fmt.Sprintf("%v\n", r))
	}
}
