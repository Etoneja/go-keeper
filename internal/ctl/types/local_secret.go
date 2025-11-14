package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/etoneja/go-keeper/internal/ctl/crypto"
)

type BaseSecret struct {
	Type     string
	Name     string
	Metadata string
}

type LocalSecret struct {
	UUID         string
	Type         string
	Name         string
	LastModified time.Time
	Hash         string
	Data         []byte
	Metadata     string
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
