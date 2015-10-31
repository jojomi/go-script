# go-script
Go library facilitating the creation of programs that resemble bash scripts.


## Rationale

[Go](https://golang.org)'s advantages like static binding and a huge modern standard library do suggest its usage for little tools that used to be implemented as shell scripts.

This library is intended as a wrapper for typical tasks shell scripts include and aimed at bringing the LOC size closer to unparalleled bash shortness.


## Methods

[![GoDoc](https://godoc.org/github.com/jojomi/go-script?status.svg)](https://godoc.org/github.com/jojomi/go-script)

Up-to-date list of methods (in source): [script.go](script.go)

The methods include helpers for executing external commands, maintaining a working directory, handling files and directories (cp/mv), and evaluating command output (exit code, stdout/stderr).


## Usage

```go
package main

import (
	"fmt"
	"github.com/jojomi/go-script"
)

func main() {
	sc := script.NewContext()
	sc.MustCommandExist("date")
	sc.SetWorkingDir("/tmp")
	sc.MustExecuteSilent("date", "-R")
	fmt.Print("The current date: ", sc.LastOutput().String())
	sc.PrintLastState()
}
```


## Warning

This library's API is not yet stable. Consider it ALPHA software at this time.
Thus you should be prepared for future API changes of any kind. In doubt, fork
away to keep a certain API status.


## Development

Comments, issues, and of course pull requests are highly welcome.


## License

see [LICENSE](LICENSE)
