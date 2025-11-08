package types

import "time"

type RemoteSecret struct {
	UUID         string
	LastModified time.Time
	Hash         string
	Data         []byte
}

func (r *RemoteSecret) GetUUID() string {
	return r.UUID
}

func (r *RemoteSecret) GetLastModified() time.Time {
	return r.LastModified
}

func (r *RemoteSecret) GetHash() string {
	return r.Hash
}

func (r *RemoteSecret) GetData() []byte {
	return r.Data
}
