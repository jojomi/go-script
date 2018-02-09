# go-script

Go library facilitating the creation of programs that resemble bash scripts.


## Rationale

[Go](https://golang.org)'s advantages like static binding and a huge modern standard library do suggest its usage for little tools that used to be implemented as shell scripts.

This library is intended as a wrapper for typical tasks shell scripts include and aimed at bringing the LOC size closer to unparalleled `bash` shortness.

`go-script` uses several other libraries that enable you to create scripts with a good user feedback and user interface on the command line.

This library strives for a good test coverage even though it is not always easy for user facing code like this.


## Methods

[![GoDoc](https://godoc.org/github.com/jojomi/go-script?status.svg)](https://godoc.org/github.com/jojomi/go-script)
![![CircleCI](https://circleci.com/gh/jojomi/go-script.svg?style=svg)](https://circleci.com/gh/jojomi/go-script)
[![Coverage Status](https://coveralls.io/repos/github/jojomi/go-script/badge.svg?branch=master)](https://coveralls.io/github/jojomi/go-script?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jojomi/go-script)](https://goreportcard.com/report/github.com/jojomi/go-script)

The methods include helpers for [executing external commands](process.go) (including [environment variables](environment.go)), maintaining a [working directory](context.go), handling [files and directories](filesystem.go) (cp/mv), and evaluating [command output](process.go) (exit code, stdout/stderr). You can use methods for [requesting input](interaction.go) from users, print [progress bars and activity indicators](progress.go), and use helpers for [printing colorful or bold text](print.go).


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

More example can be found in the `examples` directory, execute them like this:

`go run examples/command-checking/command-checking.go`


## Warning

This library's API is not yet stable. Use at your own discretion.

You should be prepared for future API changes of any kind.

In doubt, fork
away to keep a certain API status or use vendoring ([dep](https://github.com/golang/dep)) to keep your desired state.


## On The Shoulders or Giants

### Libraries Used in `go-script`

* [go-isatty](github.com/mattn/go-isatty) to detect terminal capabilities
* [survey](gopkg.in/AlecAivazis/survey.v1) for user interactions
* [wow](github.com/gernest/wow) for activity indicators
* [pb](gopkg.in/cheggaaa/pb.v1) for progress bars
* [color](https://github.com/fatih/color) for printing colorful and bold output
* [go-shutil](https://github.com/termie/go-shutil) (forked) for copying data

* [afero](github.com/spf13/afero) for abstracting filesystem for easier testing

### Other Libraries

Some libraries have proven highly useful in conjunction with `go-script`:

* [termtables](https://github.com/apcera/termtables)

More inspiration can be found at [awesome-go](https://github.com/avelino/awesome-go#command-line).


## Development

Comments, issues, and of course pull requests are highly welcome.

If you create a Merge Request, be sure to execute `./precommit.sh` beforehand.


## License

see [LICENSE](LICENSE)
