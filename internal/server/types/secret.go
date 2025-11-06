package types

import "time"

type Secret struct {
	ID           string
	UserID       string
	Data         []byte
	Hash         string
	LastModified time.Time
}
