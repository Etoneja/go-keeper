package ctl

import (
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/types"
	"github.com/manifoldco/promptui"
)

func PromptForLocalOnlyAction(localSecret *types.LocalSecret) (ActionType, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Secret '%s' exists locally but not on server", localSecret.UUID),
		Items: LocalOnlyActions,
	}
	return runActionPrompt(prompt)
}

func PromptForRemoteOnlyAction(remoteSecret *types.RemoteSecret) (ActionType, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Secret '%s' exists on server but not locally", remoteSecret.UUID),
		Items: RemoteOnlyActions,
	}
	return runActionPrompt(prompt)
}

func PromptForConflictCheckPairAction(secretCheckPair *types.SecretCheckPair) (ActionType, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Different versions of secret '%s' exists locally and on server", secretCheckPair.Local.UUID),
		Items: ConflictCheckPairActions,
	}
	return runActionPrompt(prompt)
}

func runActionPrompt(prompt promptui.Select) (ActionType, error) {
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	actionMap := map[string]ActionType{
		ActionDeleteLocal.String():   ActionDeleteLocal,
		ActionCreateRemote.String():  ActionCreateRemote,
		ActionCreateLocal.String():   ActionCreateLocal,
		ActionDeleteRemote.String():  ActionDeleteRemote,
		ActionReplaceLocal.String():  ActionReplaceLocal,
		ActionReplaceRemote.String(): ActionReplaceRemote,
		ActionSkip.String():          ActionSkip,
	}

	if action, exists := actionMap[result]; exists {
		return action, nil
	}

	return "", fmt.Errorf("unknown action: %s", result)
}
