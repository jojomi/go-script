package script

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var colorBold = color.New(color.Bold)
var printBold = colorBold.FprintFunc()
var printfBold = colorBold.FprintfFunc()
var printlnBold = colorBold.FprintlnFunc()

var colorSuccess = color.New(color.Bold, color.FgGreen)
var printSuccess = colorSuccess.FprintFunc()
var printfSuccess = colorSuccess.FprintfFunc()
var printlnSuccess = colorSuccess.FprintlnFunc()

var colorError = color.New(color.Bold, color.FgRed)
var printError = colorError.FprintFunc()
var printfError = colorError.FprintfFunc()
var printlnError = colorError.FprintlnFunc()

// PrintBold func
func (c Context) PrintBold(input ...interface{}) {
	c.terminalize(c.stdout, printBold, fmt.Fprint, input...)
}

// PrintfBold func
func (c Context) PrintfBold(format string, input ...interface{}) {
	c.terminalizef(c.stdout, printfBold, fmt.Fprintf, format, input...)
}

// PrintlnBold func
func (c Context) PrintlnBold(input ...interface{}) {
	c.terminalize(c.stdout, printlnBold, fmt.Fprintln, input...)
}

// PrintBoldCheck func
func (c Context) PrintBoldCheck(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = c.successChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	c.terminalize(c.stdout, printBold, fmt.Fprint, input...)
}

// PrintSuccess func
func (c Context) PrintSuccess(input ...interface{}) {
	c.terminalize(c.stdout, printSuccess, fmt.Fprint, input...)
}

// PrintfSuccess func
func (c Context) PrintfSuccess(format string, input ...interface{}) {
	c.terminalizef(c.stdout, printfSuccess, fmt.Fprintf, format, input...)
}

// PrintlnSuccess func
func (c Context) PrintlnSuccess(input ...interface{}) {
	c.terminalize(c.stdout, printlnSuccess, fmt.Fprintln, input...)
}

// PrintSuccessCheck func
func (c Context) PrintSuccessCheck(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = c.successChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	c.terminalize(c.stdout, printSuccess, fmt.Fprint, input...)
}

// PrintError func
func (c Context) PrintError(input ...interface{}) {
	c.terminalize(c.stderr, printError, fmt.Fprint, input...)
}

// PrintfError func
func (c Context) PrintfError(format string, input ...interface{}) {
	c.terminalizef(c.stdout, printfError, fmt.Fprintf, format, input...)
}

// PrintlnError func
func (c Context) PrintlnError(input ...interface{}) {
	c.terminalize(c.stderr, printlnError, fmt.Fprintln, input...)
}

// PrintErrorCross func
func (c Context) PrintErrorCross(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = c.errorChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	c.terminalize(c.stderr, printError, fmt.Fprint, input...)
}

// IsTerminal returns if this program is run inside an interactive terminal
func (c Context) IsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func (c Context) terminalize(w io.Writer, candy func(w io.Writer, input ...interface{}), basic func(w io.Writer, input ...interface{}) (int, error), input ...interface{}) {
	c.output(w, c.IsTerminal(), candy, basic, input...)
}

func (c Context) output(w io.Writer, isTerminal bool, candy func(w io.Writer, input ...interface{}), basic func(w io.Writer, input ...interface{}) (int, error), input ...interface{}) {
	if !isTerminal {
		basic(w, input...)
		return
	}
	candy(w, input...)
}

func (c Context) terminalizef(w io.Writer, candy func(w io.Writer, format string, input ...interface{}), basic func(w io.Writer, format string, input ...interface{}) (int, error), format string, input ...interface{}) {
	c.outputf(w, c.IsTerminal(), candy, basic, format, input...)
}

func (c Context) outputf(w io.Writer, isTerminal bool, candy func(w io.Writer, format string, input ...interface{}), basic func(w io.Writer, format string, input ...interface{}) (int, error), format string, input ...interface{}) {
	if c.PrintDetectTTY && !isTerminal {
		basic(w, format, input...)
		return
	}
	candy(w, format, input...)
}
