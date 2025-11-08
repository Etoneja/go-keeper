package types

type CheckPair struct {
	Local  *LocalSecret
	Remote *RemoteSecret
}

type SecretsDiff struct {
	LocalOnly  []*LocalSecret
	RemoteOnly []*RemoteSecret
	Both       []*CheckPair
}
