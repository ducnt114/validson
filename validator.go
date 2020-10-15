package validson

import (
	"encoding/json"
	"github.com/ducnt114/validson/models"
	"io/ioutil"
	"os"
)

type Validator interface {
	Validate(input map[string]interface{}) bool
}

type validatorImpl struct {
	schema *models.Schema
}

func NewValidatorWithSchema(schemaFilePath string) (Validator, error) {
	jsonFile, err := os.Open(schemaFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytesVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	schema := &models.Schema{}
	err = json.Unmarshal(bytesVal, &schema)
	if err != nil {
		return nil, err
	}

	return &validatorImpl{
		schema: schema,
	}, nil
}

func (v *validatorImpl) Validate(input map[string]interface{}) bool {
	// Check required fields
	for _, requireField := range v.schema.Required {
		if _, existed := input[requireField]; !existed {
			return false
		}
	}

	// TODO: check properties

	return true
}
