package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath string

func init() {
	// Try common paths in order
	paths := []string{
		"./json2struct",              // run from project root
		"../json2struct/json2struct", // run from tests/ dir
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			binaryPath = p
			return
		}
	}
	// Fallback to hardcoded path (for CI)
	binaryPath = "/root/workspace/json2struct/json2struct"
}

func runCmd(stdin string, args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
	return cmd.Output()
}

func TestGoOutput(t *testing.T) {
	out, err := runCmd(`{"name":"John","age":30,"active":true}`, "--name", "User", "--lang", "go")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "type User struct") {
		t.Errorf("expected User struct in output, got: %s", output)
	}
	if !strings.Contains(output, "Name string") {
		t.Errorf("expected Name string field in output, got: %s", output)
	}
	if !strings.Contains(output, "json:\"name\"") {
		t.Errorf("expected json tag in output, got: %s", output)
	}
}

func TestPythonOutput(t *testing.T) {
	out, err := runCmd(`{"name":"John","age":30}`, "--name", "User", "--lang", "python")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "class User") {
		t.Errorf("expected User class in output, got: %s", output)
	}
	if !strings.Contains(output, "name: str") {
		t.Errorf("expected name: str in output, got: %s", output)
	}
}

func TestTypeScriptOutput(t *testing.T) {
	out, err := runCmd(`{"name":"John","age":30}`, "--name", "User", "--lang", "typescript")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "interface User") {
		t.Errorf("expected User interface in output, got: %s", output)
	}
	if !strings.Contains(output, "name?: string") {
		t.Errorf("expected name?: string in output, got: %s", output)
	}
}

func TestRustOutput(t *testing.T) {
	out, err := runCmd(`{"name":"John","age":30}`, "--name", "User", "--lang", "rust")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "pub struct User") {
		t.Errorf("expected User struct in output, got: %s", output)
	}
	if !strings.Contains(output, "serde(rename") {
		t.Errorf("expected serde tag in output, got: %s", output)
	}
}

func TestJavaOutput(t *testing.T) {
	out, err := runCmd(`{"name":"John","age":30}`, "--name", "User", "--lang", "java")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "public class User") {
		t.Errorf("expected User class in output, got: %s", output)
	}
	if !strings.Contains(output, "getName()") {
		t.Errorf("expected getName() in output, got: %s", output)
	}
}

func TestNestedJSON(t *testing.T) {
	out, err := runCmd(`{"address":{"city":"NYC"}}`, "--name", "User", "--lang", "go")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "Address") {
		t.Errorf("expected Address field in output, got: %s", output)
	}
}

func TestArrayJSON(t *testing.T) {
	out, err := runCmd(`{"tags":["admin","user"]}`, "--name", "User", "--lang", "go")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "Tags []string") {
		t.Errorf("expected Tags []string in output, got: %s", output)
	}
}

func TestHelpFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "json2struct") {
		t.Errorf("expected json2struct in help output, got: %s", output)
	}
}

func TestMissingLang(t *testing.T) {
	_, err := runCmd(`{"name":"John"}`, "--name", "User")
	if err == nil {
		t.Error("expected error when --lang is missing")
	}
}

func TestFileInput(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`{"name":"John","age":30}`)
	tmpFile.Close()

	cmd := exec.Command(binaryPath, "--name", "User", "--lang", "go", tmpFile.Name())
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "type User struct") {
		t.Errorf("expected User struct in output, got: %s", output)
	}
}
