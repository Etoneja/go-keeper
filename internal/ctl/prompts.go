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

	actionMap := map[string]ActionType{
		ActionDeleteLocal.String():  ActionDeleteLocal,
		ActionCreateRemote.String(): ActionCreateRemote,
		ActionCreateLocal.String():  ActionCreateLocal,
		ActionDeleteRemote.String(): ActionDeleteRemote,
		ActionSkip.String():         ActionSkip,
	}

	if action, exists := actionMap[result]; exists {
		return action, nil
	}

	return "", fmt.Errorf("unknown action: %s", result)
}
