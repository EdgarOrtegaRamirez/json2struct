# json2struct — Agent Notes

## Overview
Go CLI tool that converts JSON data samples to typed structs for Go, Python, TypeScript, Rust, and Java.

## Key Files
- `cmd/json2struct/main.go` — CLI entry point (cobra)
- `internal/analyzer/analyzer.go` — JSON structure analysis
- `internal/generator/generator.go` — Code generation for 5 languages
- `tests/integration_test.go` — CLI integration tests

## Build & Test
```bash
go build -o json2struct ./cmd/json2struct/
go test ./...
./json2struct --name User --lang go < input.json
```

## Adding a Language
1. Add Language constant to `generator.go`
2. Add `generate<Language>()` method to Generator
3. Add type converter function (`jsonTypeTo<Language>`)
4. Add to CLI flag validation in `main.go`