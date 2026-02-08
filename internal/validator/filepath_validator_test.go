package validator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewDefaultFilePathValidator(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	if validator == nil {
		t.Fatal("NewDefaultFilePathValidator returned nil")
	}
}

func TestDefaultFilePathValidator_ValidatePath_EmptyPath(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	err := validator.ValidatePath("", false)
	if err == nil {
		t.Error("Empty path should return error")
	}

	expected := "文件路径不能为空"
	if err.Error() != expected {
		t.Errorf("Error message should be %q, got %q", expected, err.Error())
	}
}

func TestDefaultFilePathValidator_ValidatePath_InvalidCharacters(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	invalidPaths := []string{
		"path<with>invalid",
		"path|with|invalid",
		"path:with:invalid",
		"path?with?invalid",
		"path*with*invalid",
	}

	for _, path := range invalidPaths {
		t.Run(path, func(t *testing.T) {
			err := validator.ValidatePath(path, false)
			if err == nil {
				t.Errorf("Path with invalid characters should return error: %s", path)
			}

			expected := "文件路径包含无效字符"
			if err.Error() != expected {
				t.Errorf("Error message should be %q, got %q", expected, err.Error())
			}
		})
	}
}

func TestDefaultFilePathValidator_ValidatePath_NonExistentDirectory(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	// Create a temp directory for testing
	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, "nonexistent", "file.txt")

	err := validator.ValidatePath(nonExistentPath, false)
	if err == nil {
		t.Error("Non-existent directory should return error when createIfMissing is false")
	}

	expected := "目录不存在"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Error message should contain %q, got %q", expected, err.Error())
	}
}

func TestDefaultFilePathValidator_ValidatePath_CreateDirectory(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	// Create a temp directory for testing
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir", "subdir")
	newPath := filepath.Join(newDir, "file.txt")

	err := validator.ValidatePath(newPath, true)
	if err != nil {
		t.Errorf("Creating directory should succeed: %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(newDir)
	if err != nil {
		t.Errorf("Created directory should exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("Created path should be a directory")
	}
}

func TestDefaultFilePathValidator_ValidatePath_NotADirectory(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	// Create a temp file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "notadir")
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Try to validate a path where the parent is a file
	invalidPath := filepath.Join(filePath, "subfile.txt")

	err = validator.ValidatePath(invalidPath, false)
	if err == nil {
		t.Error("Path through a file should return error")
	}

	expected := "路径不是目录"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Error message should contain %q, got %q", expected, err.Error())
	}
}

func TestDefaultFilePathValidator_ValidatePath_NoWritePermission(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	validator := NewDefaultFilePathValidator()

	// Create a read-only directory
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	err := os.Mkdir(readOnlyDir, 0444)
	if err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Cleanup: restore permissions

	testPath := filepath.Join(readOnlyDir, "file.txt")

	err = validator.ValidatePath(testPath, false)
	if err == nil {
		t.Error("Read-only directory should return error")
	}

	// Note: The exact error message may vary by OS, so just check that we got an error
	if err == nil {
		t.Error("Expected error for read-only directory")
	}
}

func TestDefaultFilePathValidator_ValidatePath_ValidPath(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	// Create a temp directory
	tmpDir := t.TempDir()
	validPath := filepath.Join(tmpDir, "file.txt")

	err := validator.ValidatePath(validPath, false)
	if err != nil {
		t.Errorf("Valid path should not return error: %v", err)
	}
}

func TestDefaultFilePathValidator_ValidatePath_ValidPathWithSubdirectory(t *testing.T) {
	validator := NewDefaultFilePathValidator()

	// Create a temp directory with subdirectory
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	validPath := filepath.Join(subDir, "file.txt")

	err = validator.ValidatePath(validPath, false)
	if err != nil {
		t.Errorf("Valid path with subdirectory should not return error: %v", err)
	}
}
