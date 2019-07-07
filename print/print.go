package print

import (
	"os"

	"github.com/fatih/color"
)

var (
	// SucessChar is a single character signalling a successful operation
	SuccessChar = "✓"
	// ErrorChar is a single character signalling an error
	ErrorChar = "✗"
)

var colorTitle = color.New(color.Bold, color.FgGreen)
var printTitle = colorTitle.FprintFunc()
var colorSubtitle = color.New(color.Bold, color.FgBlue)
var printSubtitle = colorSubtitle.FprintFunc()

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

// Title func
func Title(input ...interface{}) {
	printTitle(os.Stdout, "» ")
	printTitle(os.Stdout, input...)
	printTitle(os.Stdout, "\n")
}

// Subtitle func
func Subtitle(input ...interface{}) {
	printSubtitle(os.Stdout, "› ")
	printSubtitle(os.Stdout, input...)
	printSubtitle(os.Stdout, "\n")
}

// Bold func
func Bold(input ...interface{}) {
	printBold(os.Stdout, input...)
}

// Boldf func
func Boldf(format string, input ...interface{}) {
	printfBold(os.Stdout, format, input...)
}

// Boldln func
func Boldln(input ...interface{}) {
	printlnBold(os.Stdout, input...)
}

// Success func
func Success(input ...interface{}) {
	printSuccess(os.Stdout, input...)
}

// Successf func
func Successf(format string, input ...interface{}) {
	printfSuccess(os.Stdout, format, input...)
}

// Successln func
func Successln(input ...interface{}) {
	printlnSuccess(os.Stdout, input...)
}

// SuccessCheck func
func SuccessCheck(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = SuccessChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	printSuccess(os.Stdout, input...)
}

// Error func
func Error(input ...interface{}) {
	printError(os.Stderr, input...)
}

// Errorf func
func Errorf(format string, input ...interface{}) {
	printfError(os.Stderr, format, input...)
}

// Errorln func
func Errorln(input ...interface{}) {
	printlnError(os.Stderr, input...)
}

// ErrorCross func
func ErrorCross(inputSuffix ...interface{}) {
	input := make([]interface{}, len(inputSuffix)+1)
	input[0] = ErrorChar
	for index, i := range inputSuffix {
		input[index+1] = i
	}
	printError(os.Stderr, input...)
}
