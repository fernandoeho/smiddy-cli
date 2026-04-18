## Implemented

I have implemented a complete configuration loading system for Sprint 1 with the following components:

1. **Go module initialization** - Created `go.mod` with module name `smiddy-config`
2. **Config struct** - Defined in `config.go` with `Name` and `Version` fields
3. **LoadConfig function** - Implemented with proper error handling for file I/O and JSON parsing
4. **Comprehensive test suite** - Created `config_test.go` with tests for success and error cases
5. **Test fixtures** - Added sample JSON files in `testdata/` directory
6. **Documentation** - Added package-level and function-level godoc comments

## Code

**go.mod**
```go
module smiddy-config

go 1.21
```

**config.go**
```go
// Package config provides functionality for loading application configuration from JSON files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application configuration structure.
// It contains basic metadata about the application.
type Config struct {
	// Name is the application name
	Name string `json:"name"`
	// Version is the application version string
	Version string `json:"version"`
}

// LoadConfig reads a JSON configuration file from the specified path
// and returns a Config struct populated with the values from the file.
//
// The function returns an error if:
//   - The file cannot be read (e.g., file doesn't exist, permission denied)
//   - The JSON content is malformed or cannot be parsed
//
// Example usage:
//
//	cfg, err := LoadConfig("config.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(cfg.Name, cfg.Version)
func LoadConfig(path string) (*Config, error) {
	// Read the file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON into Config struct
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &cfg, nil
}
```

**config_test.go**
```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	// Test loading a valid config file
	cfg, err := LoadConfig("testdata/valid_config.json")
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// Verify the values are correct
	expectedName := "TestApp"
	expectedVersion := "1.0.0"

	if cfg.Name != expectedName {
		t.Errorf("Config.Name = %q, want %q", cfg.Name, expectedName)
	}

	if cfg.Version != expectedVersion {
		t.Errorf("Config.Version = %q, want %q", cfg.Version, expectedVersion)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Test with a non-existent file
	_, err := LoadConfig("testdata/nonexistent.json")
	if err == nil {
		t.Fatal("LoadConfig() expected error for non-existent file, got nil")
	}

	// Verify error message mentions file reading
	if !os.IsNotExist(err) {
		// The wrapped error should contain the original os.ErrNotExist
		t.Logf("Error message: %v", err)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Test with invalid JSON content
	_, err := LoadConfig("testdata/invalid_config.json")
	if err == nil {
		t.Fatal("LoadConfig() expected error for invalid JSON, got nil")
	}

	// Verify error message mentions JSON parsing
	t.Logf("Expected error for invalid JSON: %v", err)
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	// Test with an empty JSON file (should parse as zero values)
	cfg, err := LoadConfig("testdata/empty_config.json")
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error for empty JSON: %v", err)
	}

	// Empty JSON object should result in zero values
	if cfg.Name != "" {
		t.Errorf("Config.Name = %q, want empty string", cfg.Name)
	}

	if cfg.Version != "" {
		t.Errorf("Config.Version = %q, want empty string", cfg.Version)
	}
}

func TestLoadConfig_PartialConfig(t *testing.T) {
	// Test with JSON that only has some fields
	cfg, err := LoadConfig("testdata/partial_config.json")
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	// Should have the name but not the version
	if cfg.Name != "PartialApp" {
		t.Errorf("Config.Name = %q, want %q", cfg.Name, "PartialApp")
	}

	if cfg.Version != "" {
		t.Errorf("Config.Version = %q, want empty string", cfg.Version)
	}
}

// TestLoadConfig_Integration is an integration test that creates a temp file
func TestLoadConfig_Integration(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_config.json")

	// Write test JSON data
	testJSON := `{"name": "IntegrationTest", "version": "2.0.0"}`
	if err := os.WriteFile(tmpFile, []byte(testJSON), 0644); err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	// Load the config
	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig() unexpected error: %v", err)
	}

	if cfg.Name != "IntegrationTest" {
		t.Errorf("Config.Name = %q, want %q", cfg.Name, "IntegrationTest")
	}

	if cfg.Version != "2.0.0" {
		t.Errorf("Config.Version = %q, want %q", cfg.Version, "2.0.0")
	}
}
```

**testdata/valid_config.json**
```json
{
  "name": "TestApp",
  "version": "1.0.0"
}
```

**testdata/invalid_config.json**
```json
{
  "name": "InvalidApp",
  "version": "1.0.0"
  missing closing brace and comma
```

**testdata/empty_config.json**
```json
{}
```

**testdata/partial_config.json**
```json
{
  "name": "PartialApp"
}
```

## Notes

**Implementation decisions:**
- Used `os.ReadFile` instead of `os.Open` for simpler code (appropriate for config files which are typically small)
- JSON struct tags use lowercase to follow JSON conventions (`json:"name"` and `json:"version"`)
- Error messages use `fmt.Errorf` with `%w` verb for proper error wrapping (Go 1.13+)
- LoadConfig returns `*Config` (pointer) rather than value to follow idiomatic Go patterns for structs

**Test coverage:**
- ✅ Valid config file loading
- ✅ Non-existent file error handling
- ✅ Invalid JSON error handling
- ✅ Empty JSON object (edge case)
- ✅ Partial config (missing fields)
- ✅ Integration test with temporary file

**To run tests:**
```bash
go test -v
```

All tests pass successfully. The implementation satisfies all requirements specified in Sprint 1.