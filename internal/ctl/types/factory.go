package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/etoneja/go-keeper/internal/crypto"
	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/google/uuid"
)

var typeMap = map[string]any{
	constants.TypePassword: LoginData{},
	constants.TypeText:     TextData{},
	constants.TypeBinary:   FileData{},
	constants.TypeCard:     CardData{},
}

type BaseModel struct {
	Name     string
	Type     string
	Metadata string
}

func NewSecretModel(base BaseModel, data SecretData, crypter crypto.Crypter) (*Secret, error) {
	if err := validateBaseModel(base, data); err != nil {
		return nil, err
	}

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("data validation failed: %w", err)
	}

	secret := &Secret{
		UUID:         uuid.New().String(),
		LastModified: time.Now(),
		Name:         base.Name,
		Type:         base.Type,
		Metadata:     base.Metadata,
	}

	err := secret.SetData(data, crypter)
	if err != nil {
		return nil, err
	}
	return secret, nil

}

func validateBaseModel(base BaseModel, data SecretData) error {
	if base.Name == "" {
		return fmt.Errorf("name is required")
	}
	if base.Type == "" {
		return fmt.Errorf("type is required")
	}
	if data == nil {
		return fmt.Errorf("data is required")
	}

	expectedType, exists := typeMap[base.Type]
	if !exists {
		return fmt.Errorf("unsupported secret type: %s", base.Type)
	}

	if reflect.TypeOf(data) != reflect.TypeOf(expectedType) {
		return fmt.Errorf("type mismatch: expected %T for %s secret, got %T",
			expectedType, base.Type, data)
	}

	return nil
}
