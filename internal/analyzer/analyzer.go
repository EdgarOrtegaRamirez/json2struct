package analyzer

import "fmt"

// Type represents a JSON value type
type Type string

const (
	TypeString Type = "string"
	TypeNumber Type = "number"
	TypeBool   Type = "boolean"
	TypeNull   Type = "null"
	TypeArray  Type = "array"
	TypeObject Type = "object"
)

// Field represents a single field in a JSON object
type Field struct {
	Name    string
	Type    Type
	SubType string // For arrays: the item type, or for objects: nested field types
	Enum    []string
	Required bool
}

// Schema represents a parsed JSON schema from data samples
type Schema struct {
	Type   Type
	Fields []Field
	Enum   []string
}

// Analyze analyzes a JSON value and returns a schema
func Analyze(data map[string]interface{}) *Schema {
	schema := &Schema{Type: TypeObject}
	for key, value := range data {
		field := analyzeValue(key, value)
		schema.Fields = append(schema.Fields, field)
	}
	return schema
}

func analyzeValue(key string, value interface{}) Field {
	f := Field{Name: key}

	switch v := value.(type) {
	case string:
		f.Type = TypeString
		// Try to detect enums (short strings that look like fixed values)
		if len(v) < 50 {
			f.Enum = []string{v}
		}
	case float64:
		f.Type = TypeNumber
	case bool:
		f.Type = TypeBool
	case nil:
		f.Type = TypeNull
	case []interface{}:
		f.Type = TypeArray
		if len(v) > 0 {
			item := analyzeValue("item", v[0])
			f.SubType = string(item.Type)
			if len(item.Enum) > 0 {
				f.Enum = item.Enum
			}
		}
	case map[string]interface{}:
		f.Type = TypeObject
		schema := Analyze(v)
		var subFields []string
		for _, field := range schema.Fields {
			subFields = append(subFields, fmt.Sprintf("%s:%s", field.Name, field.Type))
		}
		f.SubType = fmt.Sprintf("{%s}", joinStrings(subFields, ","))
	}
	return f
}

func joinStrings(s []string, sep string) string {
	result := ""
	for i, s := range s {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}