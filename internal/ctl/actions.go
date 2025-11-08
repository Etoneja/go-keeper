package ctl

type ActionType string

const (
	ActionDeleteLocal  ActionType = "delete_local"
	ActionCreateRemote ActionType = "create_remote"

	ActionCreateLocal  ActionType = "create_local"
	ActionDeleteRemote ActionType = "delete_remote"

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

func (a ActionType) String() string {
	return string(a)
}
