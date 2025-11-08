package types

import (
	"encoding/json"
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
)

func parseSecretData(secretType string, inData []byte) (SecretData, error) {
	switch secretType {
	case constants.SecretTypePassword:
		var outData LoginData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case constants.SecretTypeText:
		var outData TextData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case constants.SecretTypeBinary:
		var outData FileData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case constants.SecretTypeCard:
		var outData CardData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	default:
		return nil, fmt.Errorf("unknown secret type: %s", secretType)
	}
}
