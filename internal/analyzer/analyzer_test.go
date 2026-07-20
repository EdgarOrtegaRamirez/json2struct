package analyzer

import (
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	data := map[string]interface{}{
		"name":    "John",
		"age":     float64(30),
		"active":  true,
		"address": map[string]interface{}{"city": "NYC"},
		"tags":    []interface{}{"admin", "user"},
		"score":   nil,
	}

	schema := Analyze(data)

	if schema.Type != TypeObject {
		t.Errorf("expected type object, got %s", schema.Type)
	}

	if len(schema.Fields) != 6 {
		t.Errorf("expected 6 fields, got %d", len(schema.Fields))
	}

	// Find the name field
	var nameField *Field
	for i := range schema.Fields {
		if schema.Fields[i].Name == "name" {
			nameField = &schema.Fields[i]
			break
		}
	}

	if nameField == nil {
		t.Fatal("name field not found")
	}
	if nameField.Type != TypeString {
		t.Errorf("expected name type string, got %s", nameField.Type)
	}

	// Find the age field
	var ageField *Field
	for i := range schema.Fields {
		if schema.Fields[i].Name == "age" {
			ageField = &schema.Fields[i]
			break
		}
	}

	if ageField == nil {
		t.Fatal("age field not found")
	}
	if ageField.Type != TypeNumber {
		t.Errorf("expected age type number, got %s", ageField.Type)
	}

	// Find the active field
	var activeField *Field
	for i := range schema.Fields {
		if schema.Fields[i].Name == "active" {
			activeField = &schema.Fields[i]
			break
		}
	}

	if activeField == nil {
		t.Fatal("active field not found")
	}
	if activeField.Type != TypeBool {
		t.Errorf("expected active type boolean, got %s", activeField.Type)
	}

	// Find the tags field
	var tagsField *Field
	for i := range schema.Fields {
		if schema.Fields[i].Name == "tags" {
			tagsField = &schema.Fields[i]
			break
		}
	}

	if tagsField == nil {
		t.Fatal("tags field not found")
	}
	if tagsField.Type != TypeArray {
		t.Errorf("expected tags type array, got %s", tagsField.Type)
	}
	if tagsField.SubType != "string" {
		t.Errorf("expected tags subType string, got %s", tagsField.SubType)
	}
}

func TestAnalyzeJSONEmpty(t *testing.T) {
	schema := Analyze(map[string]interface{}{})
	if schema.Type != TypeObject {
		t.Errorf("expected type object, got %s", schema.Type)
	}
	if len(schema.Fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(schema.Fields))
	}
}

func TestAnalyzeJSONNested(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"address": map[string]interface{}{
				"city": "NYC",
			},
		},
	}

	schema := Analyze(data)
	if len(schema.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(schema.Fields))
	}

	userField := schema.Fields[0]
	if userField.Type != TypeObject {
		t.Errorf("expected user type object, got %s", userField.Type)
	}

	// Check nested fields
	if len(userField.SubType) == 0 {
		t.Error("expected nested subType for object field")
	}
}

func TestAnalyzeJSONArray(t *testing.T) {
	data := map[string]interface{}{
		"items": []interface{}{float64(1), float64(2), float64(3)},
	}

	schema := Analyze(data)
	if len(schema.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(schema.Fields))
	}

	itemsField := schema.Fields[0]
	if itemsField.Type != TypeArray {
		t.Errorf("expected items type array, got %s", itemsField.Type)
	}
	if itemsField.SubType != "number" {
		t.Errorf("expected items subType number, got %s", itemsField.SubType)
	}
}

func TestFieldStruct(t *testing.T) {
	f := Field{
		Name:     "test",
		Type:     TypeString,
		Enum:     []string{"a", "b"},
		Required: true,
	}

	if f.Name != "test" {
		t.Errorf("expected name test, got %s", f.Name)
	}
	if f.Type != TypeString {
		t.Errorf("expected type string, got %s", f.Type)
	}
	if !reflect.DeepEqual(f.Enum, []string{"a", "b"}) {
		t.Errorf("expected enum [a b], got %v", f.Enum)
	}
	if !f.Required {
		t.Error("expected required true")
	}
}
