package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilePathValidator defines the interface for validating file paths.
type FilePathValidator interface {
	// ValidatePath checks if a file path is valid and writable.
	//
	// Parameters:
	//   path - The file path to validate
	//   createIfMissing - Create parent directories if they don't exist
	//
	// Returns:
	//   error - Err if path is invalid or not writable
	ValidatePath(path string, createIfMissing bool) error
}

// DefaultFilePathValidator implements FilePathValidator.
type DefaultFilePathValidator struct{}

// NewDefaultFilePathValidator creates a new file path validator.
func NewDefaultFilePathValidator() *DefaultFilePathValidator {
	return &DefaultFilePathValidator{}
}

// ValidatePath checks if a file path is valid and writable.
func (v *DefaultFilePathValidator) ValidatePath(path string, createIfMissing bool) error {
	if path == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// Check for invalid characters
	if strings.ContainsAny(path, "<>:|?*") {
		return fmt.Errorf("文件路径包含无效字符")
	}

	// Get the directory part
	dir := filepath.Dir(path)

	// Check if directory exists
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if createIfMissing {
				// Create directory
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("无法创建目录 %s: %w", dir, err)
				}
				return nil
			}
			return fmt.Errorf("目录不存在: %s", dir)
		}
		return fmt.Errorf("无法访问目录 %s: %w", dir, err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("路径不是目录: %s", dir)
	}

	// Check write permission by creating a temporary file
	testFile := filepath.Join(dir, ".write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("目录没有写入权限: %s", dir)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}
