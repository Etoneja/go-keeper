package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/etoneja/go-keeper/internal/crypto"
)

type BaseSecret struct {
	Type     string
	Name     string
	Metadata string
}

type LocalSecret struct {
	UUID         string    `json:"uuid"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
	Hash         string    `json:"hash"`
	Data         []byte    `json:"data"`
	Metadata     string    `json:"metadata,omitempty"`
}

func (s *LocalSecret) ParseData() (SecretData, error) {
	return parseSecretData(s.Type, s.Data)
}

func (s *LocalSecret) SetData(cryptor crypto.Cryptor, data SecretData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// TODO: can be more efficient
	hashData := fmt.Sprintf("%s%s", string(jsonData), s.Metadata)
	hash := cryptor.CalculateDataHash([]byte(hashData))

	s.Data = jsonData
	s.Hash = hash
	return nil
}
