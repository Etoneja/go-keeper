package stypes

import "time"

type Secret struct {
	ID           string
	UserID       string
	LastModified time.Time
	Hash         string
	Data         []byte
}

type User struct {
	ID           string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}
