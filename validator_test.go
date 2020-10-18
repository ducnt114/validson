package validson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestValidate_Required_Success(t *testing.T) {
	schema, err := LoadSchemaFromFile("./test_data/schemas/schema_1.json")
	if err != nil {
		t.Fatal(err)
	}

	type str1 struct {
		ID    int64   `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	jsonInput := &str1{}

	inputData := `
{
    "id": 1,
    "name": "Lampshade",
    "price": 9.4
}
`
	err = json.Unmarshal([]byte(inputData), &jsonInput)
	if err != nil {
		t.Fatal(err)
	}

	res1, _ := schema.Validate(jsonInput)
	if res1 == false {
		t.FailNow()
	}
}

func TestValidate_Required_Fail(t *testing.T) {
	schema, err := LoadSchemaFromFile("./test_data/schemas/schema_1.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.SetToErrorMessageFunc(func(rule string, fName string, ruleVal interface{}) string {
		switch rule {
		case RuleRequire:
			return fmt.Sprintf("Field \"%v\" is required", fName)
		default:
			return "Validate fail"
		}
	})

	type str1 struct {
		ID    int64   `json:"id"`
		Price float64 `json:"price"`
	}

	jsonInput := &str1{}

	inputData := `
{
    "id": 1,
    "price": 9.6
}
`
	_ = json.Unmarshal([]byte(inputData), &jsonInput)

	res1, _ := schema.Validate(jsonInput)
	if res1 == true {
		t.FailNow()
	}
}

func TestValidate_Type_Fail(t *testing.T) {
	schema, err := LoadSchemaFromFile("./test_data/schemas/schema_1.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.SetToErrorMessageFunc(func(rule string, fName string, ruleVal interface{}) string {
		switch rule {
		case RuleType:
			return fmt.Sprintf("Type of \"%v\" must be %v", fName, ruleVal)
		default:
			return "Validate fail"
		}
	})

	type str1 struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Price string `json:"price"`
	}

	jsonInput := &str1{}

	inputData := `
{
    "id": 1,
	"name": "Product Name"
    "price": "123"
}
`
	_ = json.Unmarshal([]byte(inputData), &jsonInput)

	res1, errMsgList := schema.Validate(jsonInput)
	if res1 == true {
		t.FailNow()
	}
	if errMsgList[0] != "Type of \"price\" must be number" {
		t.FailNow()
	}
}
