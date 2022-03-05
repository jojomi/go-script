package interview

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// Confirm func
func Confirm(question string, defaultValue bool) (result bool, err error) {
	prompt := &survey.Confirm{
		Message: question,
		Default: defaultValue,
	}
	err = survey.AskOne(prompt, &result, nil)
	return
}

func ConfirmNoDefault(question string) (result bool, err error) {
	q := &survey.Input{
		Message: question + " (y/n)",
	}
	v := func(val interface{}) error {
		str := strings.ToLower(val.(string))
		if str != "y" && str != "yes" && str != "n" && str != "no" {
			return fmt.Errorf("Invalid input. Please type \"y\" for yes or \"n\" for no.")
		}
		return nil
	}
	var res string
	err = survey.AskOne(q, &res, survey.WithValidator(v))

	if strings.ToLower(res) == "y" || strings.ToLower(res) == "yes" {
		result = true
		return
	}
	if strings.ToLower(res) == "n" || strings.ToLower(res) == "no" {
		result = false
		return
	}

	return
}

// ChooseOneString func
func ChooseOneString(question string, options []string) (result string, err error) {
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &result, nil)
	return
}

// ChooseMultiStrings func
func ChooseMultiStrings(question string, options []string) (results []string, err error) {
	prompt := &survey.MultiSelect{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &results, nil)
	return
}
