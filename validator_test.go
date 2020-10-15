package validson

import (
	"encoding/json"
	"testing"
)

func TestValidate_Required_Success(t *testing.T) {
	validator, err := NewValidatorWithSchema("./test_data/schemas/schema_1.json")
	if err != nil {
		t.Fatal(err)
	}

	jsonInput := make(map[string]interface{})

	inputData := `
{
    "id": 1,
    "name": "Lampshade",
    "price": 9
}
`
	_ = json.Unmarshal([]byte(inputData), &jsonInput)

	res1 := validator.Validate(jsonInput)
	if res1 == false {
		t.FailNow()
	}
}

func TestValidate_Required_Fail(t *testing.T) {
	validator, err := NewValidatorWithSchema("./test_data/schemas/schema_1.json")
	if err != nil {
		t.Fatal(err)
	}

	jsonInput := make(map[string]interface{})

	inputData := `
{
    "id": 1,
    "price": 9
}
`
	_ = json.Unmarshal([]byte(inputData), &jsonInput)

	res1 := validator.Validate(jsonInput)
	if res1 == true {
		t.FailNow()
	}
}
