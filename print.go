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
var printlnBold = colorBold.FprintlnFunc()

var colorSuccess = color.New(color.Bold, color.FgGreen)
var printSuccess = colorSuccess.FprintFunc()
var printlnSuccess = colorSuccess.FprintlnFunc()

var colorError = color.New(color.Bold, color.FgRed)
var printError = colorError.FprintFunc()
var printlnError = colorError.FprintlnFunc()

// PrintBold func
func (c Context) PrintBold(input ...interface{}) {
	c.terminalize(c.stdout, printBold, fmt.Fprint, input...)
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
