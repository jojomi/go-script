package interview

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jojomi/go-script/v2/terminal"
)

func ChooseOneStringableWithDefault[T fmt.Stringer](question string, options []T, defaultValue T) (T, error) {
	return ChooseOneWithMapperAndDefault(question, options, defaultValue, func(t T) string {
		return t.String()
	})
}

func ChooseOneStringable[T fmt.Stringer](question string, options []T) (T, error) {
	return ChooseOneWithMapper(question, options, func(t T) string {
		return t.String()
	})
}

func ChooseOneWithMapperAndDefault[T any](question string, options []T, defaultValue T, mapper func(t T) string) (T, error) {
	t, err := ChooseOneWithMapper(question, options, mapper)
	if err != nil && !terminal.IsNonInteractiveTerminalError(err) {
		return t, err
	}
	return terminal.SetNonInteractiveDefault(err, t, defaultValue), nil
}

func ChooseOneWithMapper[T any](question string, options []T, mapper func(t T) string) (T, error) {
	var tResult T

	if !terminal.IsInteractive() {
		return tResult, &terminal.NonInteractiveTerminalErrorValue
	}

	// map options
	stringOpts := make([]string, len(options))
	for i, option := range options {
		stringOpts[i] = mapper(option)
	}

	prompt := &survey.Select{
		Message: question,
		Options: stringOpts,
	}

	var resultIndex int
	err := survey.AskOne(prompt, &resultIndex, nil)
	if err != nil {
		return tResult, err
	}

	// map back
	tResult = options[resultIndex]

	return tResult, nil
}
