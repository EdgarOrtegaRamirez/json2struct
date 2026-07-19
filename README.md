# json2struct

Convert JSON data to typed structs for multiple languages. A fast, single-binary CLI tool written in Go.

## Features

- **Multi-language output** — Go, Python, TypeScript, Rust, Java
- **Nested objects** — Recursive type inference for nested structures
- **Array detection** — Infers item types from array elements
- **YAML input** — Also accepts YAML input via `--format yaml`
- **File or stdin** — Read from file path or pipe via stdin

## Install

```bash
go install github.com/EdgarOrtegaRamirez/json2struct/cmd/json2struct@latest
```

## Quick Start

### Go

```bash
echo '{"name":"John","age":30,"active":true}' | json2struct --name User --lang go
```

**Output:**
```go
package user

// User is auto-generated from JSON data
type User struct {
    Name string `json:"name"`
    Age float64 `json:"age"`
    Active bool `json:"active"`
}
```

### Python

```bash
echo '{"name":"John","age":30,"active":true}' | json2struct --name User --lang python
```

**Output:**
```python
from typing import Optional, List, Union, Any

class User:
    """Auto-generated from JSON data."""

    name: str = None
    age: float = None
    active: bool = None
```

### TypeScript

```bash
echo '{"name":"John","age":30}' | json2struct --name User --lang typescript
```

**Output:**
```typescript
export interface User {
    name?: string;
    age?: number;
}
```

### Rust

```bash
echo '{"name":"John","age":30}' | json2struct --name User --lang rust
```

### Java

```bash
echo '{"name":"John","age":30}' | json2struct --name User --lang java
```

## CLI Reference

```
json2struct <flags> [input_file]

Flags:
  -n, --name string   Name for the generated struct (default "Data")
  -l, --lang string   Target language (go, python, typescript, rust, java) [required]
  -f, --format string Input format (json or yaml) (default "json")
  -h, --help          Help
```

### Subcommand

```
json2struct output <language> <input_file_or_stdin>

Flags:
  -n, --name string   Name for the generated struct (default "Data")
  -f, --format string Input format (json or yaml) (default "json")
```

## Type Mapping

| JSON Type | Go | Python | TypeScript | Rust | Java |
|-----------|-----|--------|------------|------|------|
| string | string | str | string | String | String |
| number | float64 | float | number | f64 | Double |
| boolean | bool | bool | boolean | bool | Boolean |
| null | any | None | null | serde_json::Value | Object |
| array | []T | List[T] | Array<T> | Vec<T> | List<Object> |
| object | map[string]interface{} | dict | Record<string, any> | serde_json::Value | Map<String, Object> |

## License

MIT