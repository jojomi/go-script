# go-script
Go library facilitating the creation of programs that resemble bash scripts.


## Rationale

[Go](https://golang.org)'s advantages like static binding and a huge modern standard library do suggest its usage for little tools that used to be implemented as shell scripts.

This library is intended as a wrapper for typical tasks shell scripts include and aimed at bringing the LOC size closer to unparalleled bash shortness.


## Methods

[![GoDoc](https://godoc.org/github.com/jojomi/go-script?status.svg)](https://godoc.org/github.com/jojomi/go-script) ![Coverage](http://gocover.io/_badge/github.com/jojomi/go-script)

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
	pr := sc.MustExecuteSilent("date", "-R")
	fmt.Print("The current date: ", pr.Output())
	fmt.Println(pr.StateString())
}
```


## Warning

This library's API is not yet stable. Consider it ALPHA software at this time.
Thus you should be prepared for future API changes of any kind. In doubt, fork
away to keep a certain API status.


## Useful Libraries

Some libraries have proven highly useful in conjunction with go-script:

* [color](https://github.com/fatih/color)
* [termtables](https://github.com/apcera/termtables)
* [uiprogress](https://github.com/gosuri/uiprogress)

More inspiration can be found at [awesome-go](https://github.com/avelino/awesome-go#command-line).


## Development

Comments, issues, and of course pull requests are highly welcome.

If you create a Merge Request, be sure to execute `./precommit.sh` beforehand.


## License

see [LICENSE](LICENSE)
