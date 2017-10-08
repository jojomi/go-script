package main

import (
	"fmt"

	script "github.com/jojomi/go-script"
)

func main() {
	// make sure panics are printed in a human friendly way
	defer script.RecoverFunc()

	sc := script.NewContext()
	if sc.CommandExists("ls") {
		fmt.Println("ls found, listing files and directories should be possible.")
	}
	if !sc.CommandExists("customjava") {
		fmt.Println("No customjava found, continuing still.")
	}
	sc.MustCommandExist("custompython")
	fmt.Println("custompython found, ready to go!")
}
