package types

import (
	"encoding/json"
	"fmt"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
)

func ParseSecretData(dataType string, inData []byte) (SecretData, error) {
	switch dataType {
	case constants.TypePassword:
		var outData LoginData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case constants.TypeText:
		var outData TextData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case constants.TypeBinary:
		var outData FileData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case constants.TypeCard:
		var outData CardData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	default:
		return nil, fmt.Errorf("unknown secret type: %s", dataType)
	}
}
