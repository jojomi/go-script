package script

import (
	"fmt"
	"os"
)

type CommandLogger func(Context, Command) error

func FileCommandLogger(filename string) *CommandLogger {
	f := fileCommandLoggerWithFlags(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE)
	return &f
}

func EnvFileCommandLogger(envKey string) *CommandLogger {
	filename := os.Getenv(envKey)
	if filename == "" {
		return nil
	}
	f := fileCommandLoggerWithFlags(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE)
	return &f
}

func fileCommandLoggerWithFlags(filename string, flags int) CommandLogger {
	return func(sc Context, c Command) error {
		f, err := os.OpenFile(filename, flags, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		message := getCommandLogMessage(sc, c)
		fmt.Println("writing", message)
		_, err = f.WriteString(message)
		if err != nil {
			return err
		}
		return nil
	}
}

func getCommandLogMessage(sc Context, c Command) string {
	return fmt.Sprintf("%s: %s\n", sc.WorkingDir(), c.String())
}
