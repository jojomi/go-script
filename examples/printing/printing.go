package main

import (
	"fmt"

	"github.com/jojomi/go-script/v2/print"
)

func main() {
	print.Boldln("trying", "A")

	fmt.Print("Yes,", "indeed ")
	print.SuccessCheck("\n")
	print.SuccessCheck(" oh yes!\n")
	print.Successln("It worked")
	print.Successf("very %s\n", "well")

	fmt.Println()
	print.Boldln("B too?")

	fmt.Print("No,", "no ")
	print.ErrorCross("\n")
	print.ErrorCross(" oh no!\n")
	print.Error("It did")
	print.Errorf(" not %s", "work")
	print.Errorln()

	fmt.Println()
	fmt.Println("I'm done.", "Really.")
}
