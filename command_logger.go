package script

import (
	"fmt"
	"github.com/jojomi/strtpl"
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

func TemplateFilename(sc Context, filename string) string {
	logKey := sc.LogKey()
	if logKey == "" {
		return filename
	}
	return strtpl.MustEval(filename, map[string]any{
		"logKey": logKey,
		"start":  sc.GetStart(),
	})
}

func fileCommandLoggerWithFlags(filename string, flags int) CommandLogger {
	return func(sc Context, c Command) error {
		f, err := os.OpenFile(TemplateFilename(sc, filename), flags, 0600)
		if err != nil {
			return err
		}
		defer f.Close()
		message := getCommandLogMessage(sc, c)
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
