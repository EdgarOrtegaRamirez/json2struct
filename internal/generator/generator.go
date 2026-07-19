package generator

import (
	"fmt"
	"strings"

	"github.com/EdgarOrtegaRamirez/json2struct/internal/analyzer"
)

// Language represents a target programming language
type Language string

const (
	GoLang     Language = "go"
	Python     Language = "python"
	Typescript Language = "typescript"
	Rust       Language = "rust"
	Java       Language = "java"
)

// Generator converts JSON schemas to typed code for different languages
type Generator struct {
	StructName string
}

// New creates a new Generator with the given struct name
func New(structName string) *Generator {
	return &Generator{StructName: structName}
}

// Generate generates code for the specified language
func (g *Generator) Generate(schema *analyzer.Schema, lang Language) (string, error) {
	switch lang {
	case GoLang:
		return g.generateGo(schema), nil
	case Python:
		return g.generatePython(schema), nil
	case Typescript:
		return g.generateTypeScript(schema), nil
	case Rust:
		return g.generateRust(schema), nil
	case Java:
		return g.generateJava(schema), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}
}

func (g *Generator) generateGo(schema *analyzer.Schema) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(strings.ReplaceAll(g.StructName, " ", ""))))
	sb.WriteString(fmt.Sprintf("// %s is auto-generated from JSON data\n", g.StructName))
	sb.WriteString(fmt.Sprintf("type %s struct {\n", g.StructName))
	for _, field := range schema.Fields {
		goType := jsonTypeToGo(string(field.Type), field.SubType)
		requiredTag := ""
		if field.Required {
			requiredTag = " required"
		}
		sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"%s`\n",
			pascalCase(field.Name), goType, field.Name, requiredTag))
	}
	sb.WriteString("}\n")
	return sb.String()
}

func (g *Generator) generatePython(schema *analyzer.Schema) string {
	var sb strings.Builder
	sb.WriteString("from typing import Optional, List, Union, Any\n\n")
	sb.WriteString(fmt.Sprintf("class %s:\n", pascalCase(g.StructName)))
	sb.WriteString("    \"\"\"Auto-generated from JSON data.\"\"\"\n\n")
	for _, field := range schema.Fields {
		typeHint := jsonTypeToPython(string(field.Type), field.SubType)
		optional := ""
		if !field.Required {
			optional = " = None"
		}
		sb.WriteString(fmt.Sprintf("    %s: %s%s\n", snakeCase(field.Name), typeHint, optional))
	}
	return sb.String()
}

func (g *Generator) generateTypeScript(schema *analyzer.Schema) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("export interface %s {\n", pascalCase(g.StructName)))
	for _, field := range schema.Fields {
		tsType := jsonTypeToTypeScript(string(field.Type), field.SubType)
		optional := ""
		if !field.Required {
			optional = "?"
		}
		sb.WriteString(fmt.Sprintf("    %s%s: %s;\n", snakeCase(field.Name), optional, tsType))
	}
	sb.WriteString("}\n")
	return sb.String()
}

func (g *Generator) generateRust(schema *analyzer.Schema) string {
	var sb strings.Builder
	sb.WriteString("#![allow(dead_code)]\n\n")
	sb.WriteString("#[derive(Debug, Clone, Serialize, Deserialize)]\n")
	sb.WriteString(fmt.Sprintf("pub struct %s {\n", pascalCase(g.StructName)))
	for _, field := range schema.Fields {
		rustType := jsonTypeToRust(string(field.Type), field.SubType)
		sb.WriteString(fmt.Sprintf("\t#[serde(rename = \"%s\")]\n", field.Name))
		sb.WriteString(fmt.Sprintf("\tpub %s: %s,\n", snakeCase(field.Name), rustType))
	}
	sb.WriteString("}\n")
	return sb.String()
}

func (g *Generator) generateJava(schema *analyzer.Schema) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("public class %s {\n\n", pascalCase(g.StructName)))
	for _, field := range schema.Fields {
		javaType := jsonTypeToJava(string(field.Type), field.SubType)
		sb.WriteString(fmt.Sprintf("    private %s %s;\n\n", javaType, snakeCase(field.Name)))
	}
	for _, field := range schema.Fields {
		javaType := jsonTypeToJava(string(field.Type), field.SubType)
		sb.WriteString(fmt.Sprintf("    public %s get%s() {\n", javaType, pascalCase(field.Name)))
		sb.WriteString(fmt.Sprintf("        return %s;\n    }\n\n", snakeCase(field.Name)))
		sb.WriteString(fmt.Sprintf("    public void set%s(%s %s) {\n", pascalCase(field.Name), javaType, snakeCase(field.Name)))
		sb.WriteString(fmt.Sprintf("        this.%s = %s;\n    }\n\n", snakeCase(field.Name), snakeCase(field.Name)))
	}
	sb.WriteString("}\n")
	return sb.String()
}

// jsonTypeToGo converts a JSON type to a Go type
func jsonTypeToGo(jsonType, subType string) string {
	switch jsonType {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "null":
		return "any"
	case "array":
		if subType != "" {
			return "[]" + jsonTypeToGo(subType, "")
		}
		return "[]"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// jsonTypeToPython converts a JSON type to a Python type hint
func jsonTypeToPython(jsonType, subType string) string {
	switch jsonType {
	case "string":
		return "str"
	case "number":
		return "float"
	case "boolean":
		return "bool"
	case "null":
		return "None"
	case "array":
		if subType != "" {
			return fmt.Sprintf("List[%s]", jsonTypeToPython(subType, ""))
		}
		return "List"
	case "object":
		return "dict"
	default:
		return "Any"
	}
}

// jsonTypeToTypeScript converts a JSON type to a TypeScript type
func jsonTypeToTypeScript(jsonType, subType string) string {
	switch jsonType {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "null":
		return "null"
	case "array":
		if subType != "" {
			return fmt.Sprintf("Array<%s>", jsonTypeToTypeScript(subType, ""))
		}
		return "Array<any>"
	case "object":
		return "Record<string, any>"
	default:
		return "any"
	}
}

// jsonTypeToRust converts a JSON type to a Rust type
func jsonTypeToRust(jsonType, subType string) string {
	switch jsonType {
	case "string":
		return "String"
	case "number":
		return "f64"
	case "boolean":
		return "bool"
	case "null":
		return "serde_json::Value"
	case "array":
		if subType != "" {
			return fmt.Sprintf("Vec<%s>", jsonTypeToRust(subType, ""))
		}
		return "Vec<serde_json::Value>"
	case "object":
		return "serde_json::Value"
	default:
		return "serde_json::Value"
	}
}

// jsonTypeToJava converts a JSON type to a Java type
func jsonTypeToJava(jsonType, subType string) string {
	switch jsonType {
	case "string":
		return "String"
	case "number":
		return "Double"
	case "boolean":
		return "Boolean"
	case "null":
		return "Object"
	case "array":
		return "List<Object>"
	case "object":
		return "Map<String, Object>"
	default:
		return "Object"
	}
}

func pascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "")
}

func snakeCase(s string) string {
	result := ""
	for i, r := range s {
		if r >= 'A' && r <= 'Z' && i > 0 {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToLower(result)
}