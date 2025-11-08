package types

import "time"

type SecretData interface {
	Validate() error
}

type Secreter interface {
	GetUUID() string
	GetLastModified() time.Time
	GetHash() string
	GetData() []byte
}
