package types

type CheckPair struct {
	Local  Secreter
	Remote Secreter
}

type SecretDiff struct {
	LocalOnly  []Secreter
	RemoteOnly []Secreter
	Both       []*CheckPair
}
