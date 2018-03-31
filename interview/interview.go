package interview

import (
	survey "gopkg.in/AlecAivazis/survey.v1"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
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

// ChooseOneString func
func ChooseOneString(question string, options []string) (result string, err error) {
	surveyCore.QuestionIcon = "?"

	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &result, nil)
	return
}

// ChooseMultiStrings func
func ChooseMultiStrings(question string, options []string) (results []string, err error) {
	surveyCore.QuestionIcon = "?"

	prompt := &survey.MultiSelect{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &results, nil)
	return
}
