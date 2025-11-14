package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/etoneja/go-keeper/internal/ctl/crypto"
	"github.com/google/uuid"
)

var typeMap = map[string]any{
	constants.SecretTypePassword: LoginData{},
	constants.SecretTypeText:     TextData{},
	constants.SecretTypeBinary:   FileData{},
	constants.SecretTypeCard:     CardData{},
}

func NewSecretModel(base BaseSecret, data SecretData, cryptor crypto.Cryptor) (*LocalSecret, error) {
	if err := validateBaseSecret(base, data); err != nil {
		return nil, err
	}

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("data validation failed: %w", err)
	}

	secret := &LocalSecret{
		UUID:         uuid.New().String(),
		Type:         base.Type,
		Name:         base.Name,
		LastModified: time.Now().UTC().Truncate(time.Microsecond),
		Metadata:     base.Metadata,
	}

	if err := secret.SetData(cryptor, data); err != nil {
		return nil, err
	}
	return secret, nil

}

func validateBaseSecret(base BaseSecret, data SecretData) error {
	if base.Type == "" {
		return fmt.Errorf("type is required")
	}
	if base.Name == "" {
		return fmt.Errorf("name is required")
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
