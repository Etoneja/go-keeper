package types

type SecretCheckPair struct {
	Local  *LocalSecret
	Remote *RemoteSecret
}

type SecretsDiff struct {
	LocalOnly  []*LocalSecret
	RemoteOnly []*RemoteSecret
	Both       []*SecretCheckPair
}

func (p *SecretCheckPair) IsIdentical() bool {
	return p.Local.Hash == p.Remote.Hash &&
		p.Local.LastModified.Equal(p.Remote.LastModified)
}
