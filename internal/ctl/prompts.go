package ctl

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func PromptForLocalOnlyAction(secretID string) (ActionType, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Secret '%s' exists locally but not on server", secretID),
		Items: []string{
			ActionDeleteLocal.String(),
			ActionCreateRemote.String(),
			ActionSkip.String(),
		},
	}
	return runPrompt(prompt)
}

func PromptForRemoteOnlyAction(secretID string) (ActionType, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Secret '%s' exists on server but not locally", secretID),
		Items: []string{
			ActionCreateLocal.String(),
			ActionDeleteRemote.String(),
			ActionSkip.String(),
		},
	}
	return runPrompt(prompt)
}

func runPrompt(prompt promptui.Select) (ActionType, error) {
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	switch result {
	case ActionDeleteLocal.String():
		return ActionDeleteLocal, nil
	case ActionCreateRemote.String():
		return ActionCreateRemote, nil
	case ActionCreateLocal.String():
		return ActionCreateLocal, nil
	case ActionDeleteRemote.String():
		return ActionDeleteRemote, nil
	case ActionSkip.String():
		return ActionSkip, nil
	default:
		return "", fmt.Errorf("unknown action: %s", result)
	}
}
