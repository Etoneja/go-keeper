package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/etoneja/go-keeper/internal/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/constants"
)

// Domain model
type Secret struct {
	UUID         string    `json:"uuid"`
	LastModified time.Time `json:"last_modified"`
	Hash         string    `json:"hash"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Data         []byte    `json:"data"`
	Metadata     string    `json:"metadata,omitempty"`
}

// ParseData parses secret data based on type
func (s *Secret) ParseData() (SecretData, error) {
	switch s.Type {
	case constants.TypePassword:
		var data LoginData
		err := json.Unmarshal(s.Data, &data)
		return data, err
	case constants.TypeText:
		var data TextData
		err := json.Unmarshal(s.Data, &data)
		return data, err
	case constants.TypeBinary:
		var data FileData
		err := json.Unmarshal(s.Data, &data)
		return data, err
	case constants.TypeCard:
		var data CardData
		err := json.Unmarshal(s.Data, &data)
		return data, err
	default:
		return nil, fmt.Errorf("unknown secret type: %s", s.Type)
	}
}

// SetData serializes data to JSON based on type
func (s *Secret) SetData(data SecretData, crypter crypto.Crypter) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	hashData := fmt.Sprintf("%s%s%s%s", s.Name, s.Type, string(jsonData), s.Metadata)
	hash := crypter.GenerateHash([]byte(hashData))

	s.Data = jsonData
	s.Hash = hash
	return nil
}
