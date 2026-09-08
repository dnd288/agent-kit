package prompt

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// AskString prompts for a string value.
func AskString(label, defaultVal string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
	}
	return prompt.Run()
}

// AskStringRequired prompts for a required string value.
func AskStringRequired(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("value is required")
			}
			return nil
		},
	}
	return prompt.Run()
}

// AskSelect prompts the user to select one item from a list.
func AskSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	return result, err
}

// AskConfirm asks a yes/no question.
func AskConfirm(label string, defaultYes bool) (bool, error) {
	defVal := "n"
	if defaultYes {
		defVal = "y"
	}

	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   defVal,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// AskMultiSelect prompts the user to select multiple items.
// Returns the selected items.
func AskMultiSelect(label string, items []string) ([]string, error) {
	var selected []string

	fmt.Printf("\n%s (space to toggle, enter to confirm):\n", label)
	for i, item := range items {
		fmt.Printf("  %d. %s\n", i+1, item)
	}

	// promptui doesn't have native multi-select, so we use a loop
	for _, item := range items {
		yes, err := AskConfirm(fmt.Sprintf("  Include %s?", item), true)
		if err != nil {
			return nil, err
		}
		if yes {
			selected = append(selected, item)
		}
	}

	return selected, nil
}
