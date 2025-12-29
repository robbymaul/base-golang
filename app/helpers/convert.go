package helpers

import (
	"encoding/json"
	"gorm.io/datatypes"
)

func ConvertAnyToDatatypeJson(v any) (datatypes.JSON, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(marshal), nil
}
