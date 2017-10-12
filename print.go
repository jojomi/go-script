package script

import (
	"github.com/fatih/color"
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
	printBold(c.stdout, input...)
}

// PrintfBold func
func (c Context) PrintfBold(format string, input ...interface{}) {
	printfBold(c.stdout, format, input...)
}

// PrintlnBold func
func (c Context) PrintlnBold(input ...interface{}) {
	printlnBold(c.stdout, input...)
}

// PrintSuccess func
func (c Context) PrintSuccess(input ...interface{}) {
	printSuccess(c.stdout, input...)
}

// PrintfSuccess func
func (c Context) PrintfSuccess(format string, input ...interface{}) {
	printfSuccess(c.stdout, format, input...)
}

// PrintlnSuccess func
func (c Context) PrintlnSuccess(input ...interface{}) {
	printlnSuccess(c.stdout, input...)
}

// PrintSuccessCheck func
func (c Context) PrintSuccessCheck(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = c.successChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	printSuccess(c.stdout, input...)
}

// PrintError func
func (c Context) PrintError(input ...interface{}) {
	printError(c.stderr, input...)
}

// PrintfError func
func (c Context) PrintfError(format string, input ...interface{}) {
	printfError(c.stderr, format, input...)
}

// PrintlnError func
func (c Context) PrintlnError(input ...interface{}) {
	printlnError(c.stderr, input...)
}

// PrintErrorCross func
func (c Context) PrintErrorCross(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = c.errorChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	printError(c.stderr, input...)
}
