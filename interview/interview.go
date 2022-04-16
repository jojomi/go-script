package interview

import (
	"errors"
	"fmt"
	"github.com/jojomi/go-script/v2/terminal"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// Confirm allows querying for confirmation with a default value that is used when no answer is typed.
func Confirm(question string, defaultValue bool) (result bool, err error) {
	prompt := &survey.Confirm{
		Message: question,
		Default: defaultValue,
	}
	err = survey.AskOne(prompt, &result, nil)
	return
}

// ConfirmNoDefault allows querying for confirmation without a default value, so the user needs to answer the question explicitly.
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

// ChooseOneString queries the user to choose one string from a list of strings
func ChooseOneString(question string, options []string) (result string, err error) {
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &result, nil)
	return
}

// ChooseMultiStrings queries the user to choose from a list of strings allowing multiple selection
func ChooseMultiStrings(question string, options []string) (results []string, err error) {
	prompt := &survey.MultiSelect{
		Message: question,
		Options: options,
	}
	err = survey.AskOne(prompt, &results, nil)
	return
}

func WithAutoShortcuts(actions []*Action) []*Action {
	result := make([]*Action, len(actions))
	usedShortcuts := make(map[string]struct{}, 0)

	for i, action := range actions {
		action.Shortcut = getFirstFreeShortcut(action.Label, usedShortcuts)
		result[i] = action
	}

	return result
}

func getFirstFreeShortcut(label string, shortcuts map[string]struct{}) string {
	runes := []rune(label)
	var s string
	for _, r := range runes {
		s = strings.ToLower(string(r))
		if _, ok := shortcuts[s]; ok {
			continue
		}
		shortcuts[s] = struct{}{}
		return s
	}
	return ""
}

func SelectAction(actions []*Action) (*Action, error) {
	return selectAction(actions, nil)
}

func SelectActionWithDefault(actions []*Action, defaultAction *Action) (*Action, error) {
	return selectAction(actions, defaultAction)
}

func selectAction(actions []*Action, defaultAction *Action) (*Action, error) {
	if !terminal.IsInteractive() {
		if defaultAction != nil {
			return defaultAction, nil
		} else {
			return nil, errors.New("could not ask for an action because the terminal is not interactive")
		}
	}

	// print actions
	var s strings.Builder
	for i, action := range actions {
		if action == defaultAction {
			s.WriteString(action.ShortcutStringAsDefault())
		} else {
			s.WriteString(action.ShortcutString())
		}
		if i < len(actions)-1 {
			s.WriteString(" | ")
		}
	}
	s.WriteString(":")

	var answer string
	q := &survey.Input{Message: s.String()}
	err := survey.AskOne(q, &answer, survey.WithValidator(actionValidator(actions, defaultAction)))
	if err != nil {
		return nil, err
	}

	// match and return
	if defaultAction != nil && answer == "" {
		return defaultAction, nil
	}
	for _, action := range actions {
		if strings.EqualFold(action.Label, answer) {
			return action, nil
		}
		if strings.EqualFold(action.Shortcut, answer) {
			return action, nil
		}
	}

	// should never be reached due to survey filter and matching above
	panic("unreachable code reached")
}

func actionValidator(actions []*Action, defaultAction *Action) survey.Validator {
	return func(ans interface{}) error {
		answer := ans.(string)

		if defaultAction != nil && answer == "" {
			return nil
		}

		for _, action := range actions {
			if strings.EqualFold(action.Label, answer) {
				return nil
			}
			if strings.EqualFold(action.Shortcut, answer) {
				return nil
			}
		}
		return fmt.Errorf("action not found: %v", ans)
	}
}
