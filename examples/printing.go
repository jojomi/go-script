package main

import (
	"fmt"

	script "github.com/jojomi/go-script"
)

func main() {
	// make sure panics are printed in a human friendly way
	/// defer script.RecoverFunc()

	sc := script.NewContext()

	fmt.Println()
	sc.PrintlnBold("trying", "A")

	fmt.Print("Yes,", "indeed ")
	sc.PrintSuccessCheck("\n")
	sc.PrintSuccessCheck(" oh yes!\n")
	sc.PrintlnSuccess("It worked")

	fmt.Println()
	sc.PrintlnBold("B too?")

	fmt.Print("No,", "no ")
	sc.PrintErrorCross("\n")
	sc.PrintErrorCross(" oh no!\n")
	sc.PrintlnError("It did not work")

	fmt.Println()
	fmt.Println("I'm done.", "Really.")
}
