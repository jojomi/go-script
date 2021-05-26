package main

import (
	"fmt"
	"os"

	"github.com/jojomi/go-script/v2/interview"
	"github.com/jojomi/go-script/v2/print"
)

func main() {
	shouldContinue, err := interview.Confirm("Should we continue?", true)
	if !shouldContinue {
		os.Exit(1)
	}

	level, err := interview.ChooseOneString("What is your expertise level?", []string{"Novice", "Learner", "Professional"})
	if err != nil {
		return
	}
	switch level {
	case "Novice":
		print.Errorln(level)
	case "Learner":
		fmt.Println(level)
	case "Professional":
		print.Successln(level)
	}

	sports, err := interview.ChooseMultiStrings("What sport do you watch on tv?", []string{"Football", "Basketball", "Soccer"})
	if err != nil {
		return
	}
	fmt.Println("Your tv sports:", sports)
}
