package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/etoneja/go-keeper/internal/crypto"
)

// Domain model
type Secret struct {
	UUID         string    `json:"uuid"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	LastModified time.Time `json:"last_modified"`
	Hash         string    `json:"hash"`
	Data         []byte    `json:"data"`
	Metadata     string    `json:"metadata,omitempty"`
}

func (s *Secret) GetUUID() string {
	return s.UUID
}

func (s *Secret) GetLastModified() time.Time {
	return s.LastModified
}

func (s *Secret) GetHash() string {
	return s.Hash
}

func (s *Secret) GetData() []byte {
	return s.Data
}

// ParseData parses secret data based on type
func (s *Secret) ParseData() (SecretData, error) {
	return ParseSecretData(s.Type, s.Data)
}

// SetData serializes data to JSON based on type
func (s *Secret) SetData(data SecretData, crypter crypto.Crypter) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// TODO: can be more efficient?
	hashData := fmt.Sprintf("%s%s", string(jsonData), s.Metadata)
	hash := crypter.GenerateHash([]byte(hashData))

	s.Data = jsonData
	s.Hash = hash
	return nil
}
