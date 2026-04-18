package fs

import (
	"fmt"
	"os"
)

// WriteFile writes the given content to a file at the specified path.
func WriteFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", path, err)
	}
	return nil
}

// ReadFile reads the content of a file at the specified path and returns it as a string.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return string(data), nil
}

// NextSprintNumber returns the next sprint number based on existing sprint folders in the given base directory.
func NextSprintNumber(sprintBase string) (int, error) {
	entries, err := os.ReadDir(sprintBase)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil // No sprints exist yet, start with sprint 1
		}
		return 0, fmt.Errorf("failed to read sprint base directory %s: %w", sprintBase, err)
	}
	return len(entries) + 1, nil
}
