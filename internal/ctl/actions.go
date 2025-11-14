package ctl

type ActionType string

func (a ActionType) String() string {
	return string(a)
}

const (
	ActionDeleteLocal  ActionType = "delete_local"
	ActionCreateRemote ActionType = "create_remote"

	ActionCreateLocal  ActionType = "create_local"
	ActionDeleteRemote ActionType = "delete_remote"

	ActionReplaceLocal  ActionType = "replace_local"
	ActionReplaceRemote ActionType = "replace_remote"

	ActionSkip ActionType = "ignore"
)

var LocalOnlyActions = []ActionType{
	ActionDeleteLocal,
	ActionCreateRemote,
	ActionSkip,
}

var RemoteOnlyActions = []ActionType{
	ActionCreateLocal,
	ActionDeleteRemote,
	ActionSkip,
}

var ConflictCheckPairActions = []ActionType{
	ActionReplaceLocal,
	ActionReplaceRemote,
	ActionSkip,
}
