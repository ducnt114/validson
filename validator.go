package validson

import (
	"encoding/json"
	"github.com/ducnt114/validson/models"
	"io/ioutil"
	"os"
	"reflect"
)

const (
	RuleConvert = "convert"
	RuleRequire = "require"
	RuleType    = "type"
)

type ToErrorMessageFunc func(rule string, fName string, ruleVal interface{}) string

type Schema interface {
	Validate(input interface{}) (bool, []string)
	SetToErrorMessageFunc(f ToErrorMessageFunc)
}

type schemaImpl struct {
	jsonSchema         *models.JsonSchema
	toErrorMessageFunc ToErrorMessageFunc
}

func LoadSchemaFromFile(schemaFilePath string) (Schema, error) {
	jsonFile, err := os.Open(schemaFilePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytesVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	jsonSchema := &models.JsonSchema{}
	err = json.Unmarshal(bytesVal, &jsonSchema)
	if err != nil {
		return nil, err
	}

	return &schemaImpl{
		jsonSchema: jsonSchema,
	}, nil
}

func (s *schemaImpl) SetToErrorMessageFunc(f ToErrorMessageFunc) {
	s.toErrorMessageFunc = f
}

func (s *schemaImpl) Validate(input interface{}) (bool, []string) {
	errList := make([]string, 0)
	jsonInput, err := s.convertToJsonMap(input)
	if err != nil {
		errList = append(errList, s.toErrorMessageFunc(RuleConvert, "", ""))
		return false, errList
	}

	// Check required fields
	missingRequiredField := false
	for _, requireField := range s.jsonSchema.Required {
		if _, existed := jsonInput[requireField]; !existed {
			missingRequiredField = true
			errList = append(errList, s.toErrorMessageFunc(requireField, RuleRequire, ""))
		}
	}
	if missingRequiredField {
		return false, errList
	}

	hasFieldInvalidate := false
	for fName, fVal := range jsonInput {
		if prop, existed := s.jsonSchema.Properties[fName]; existed {
			fieldValid, errMsg := s.validateProperties(prop, fName, fVal)
			if !fieldValid {
				hasFieldInvalidate = true
				errList = append(errList, errMsg)
			}
		}
	}
	if hasFieldInvalidate {
		return false, errList
	}

	return true, []string{}
}

func (s *schemaImpl) convertToJsonMap(input interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	val := reflect.ValueOf(input).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := val.Field(i)
		jsonTag := typeField.Tag.Get("json")
		res[jsonTag] = valueField.Interface()
	}

	return res, nil
}

func (s *schemaImpl) validateProperties(prop *models.JsonProperty, fName string, fVal interface{}) (bool, string) {
	// check type
	switch prop.Type {
	case "string":
		if _, ok := fVal.(string); !ok {
			return false, s.toErrorMessageFunc(RuleType, fName, "string")
		}
	case "number":
		if _, ok := fVal.(int64); !ok {
			if _, ok = fVal.(float64); !ok {
				return false, s.toErrorMessageFunc(RuleType, fName, "number")
			}
		}
	case "boolean":
		if _, ok := fVal.(bool); !ok {
			return false, s.toErrorMessageFunc(RuleType, fName, "boolean")
		}
	default:
		return false, s.toErrorMessageFunc(RuleType, fName, "")
	}

	return true, ""
}
