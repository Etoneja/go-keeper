package types

import "time"

type Secret struct {
	ID           string
	UserID       string
	LastModified time.Time
	Hash         string
	Data         []byte
}
