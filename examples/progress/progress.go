package main

import (
	"fmt"
	"time"

	script "github.com/jojomi/go-script"
)

func main() {
	sc := script.NewContext()
	ai := sc.ActivityIndicator("Loading")
	ai.Start()
	time.Sleep(3 * time.Second)
	ai.Text(" Finished.")
	ai.Persist()
	fmt.Println("All done, goodbye!")
}
