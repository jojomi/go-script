package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		return
	}

	var err error
	file, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var (
		input    string
		value    string
		intValue int
	)
	for scanner.Scan() {
		input = scanner.Text()

		if strings.HasPrefix(input, "out: ") {
			value = input[len("out: "):]
			os.Stdout.Write([]byte(value + "\n"))
			continue
		}
		if strings.HasPrefix(input, "err: ") {
			value = input[len("err: "):]
			os.Stderr.Write([]byte(value + "\n"))
			continue
		}
		if strings.HasPrefix(input, "echo: ") {
			stdinScanner := bufio.NewScanner(os.Stdin)
			stdinScanner.Scan()
			value = stdinScanner.Text() + "\n"
			os.Stdout.Write([]byte(value))
			os.Stderr.Write([]byte(value))
			continue
		}
		if strings.HasPrefix(input, "exit: ") {
			value = input[len("exit: "):]
			intValue, err = strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			os.Exit(intValue)
			continue
		}
		if strings.HasPrefix(input, "sleep: ") {
			value = input[len("sleep: "):]
			intValue, err = strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Millisecond * time.Duration(intValue))
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
