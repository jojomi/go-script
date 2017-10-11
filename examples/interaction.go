package main

import (
	"fmt"
	"os"

	script "github.com/jojomi/go-script"
)

func main() {
	// make sure panics are printed in a human friendly way
	defer script.RecoverFunc()

	sc := script.NewContext()

	shouldContinue, err := sc.Confirm("Should we continue?", true)
	if !shouldContinue {
		os.Exit(1)
	}

	level, err := sc.ChooseOneString("What is your expertise level?", []string{"Novice", "Learner", "Professional"})
	if err != nil {
		return
	}
	switch level {
	case "Novice":
		sc.PrintlnError(level)
	case "Learner":
		fmt.Println(level)
	case "Professional":
		sc.PrintlnSuccess(level)
	}

	sports, err := sc.ChooseMultiStrings("What sport do you watch on tv?", []string{"Football", "Basketball", "Soccer"})
	if err != nil {
		return
	}
	fmt.Println("Your tv sports:", sports)
}
