package types

import (
	"time"
)

type RemoteSecret struct {
	UUID         string
	LastModified time.Time
	Hash         string
	Data         []byte
}
