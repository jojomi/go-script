package main

import (
	"fmt"
	"strings"

	script "github.com/jojomi/go-script"
)

func main() {
	// make sure panics are printed in a human friendly way
	defer script.RecoverFunc()

	sc := script.NewContext()
	err := sc.SetWorkingDirTemp()
	if err != nil {
		panic(err)
	}
	name, args := script.SplitCommand("ls -lahr /")
	pr, err := sc.ExecuteFullySilent(name, args...)
	if err != nil {
		panic(err)
	}
	if strings.Contains(pr.Output(), "etc") {
		fmt.Println("/etc is in output!")
	}
}
