package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateFileMD5(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("Hello, World!")
	
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate MD5
	md5sum, err := CalculateFileMD5(testFile)
	if err != nil {
		t.Fatalf("CalculateFileMD5 failed: %v", err)
	}

	// Expected MD5 for "Hello, World!"
	expected := "65a8e27d8879283831b664bd8b7f0ad4"
	if md5sum != expected {
		t.Errorf("MD5 mismatch: got %s, want %s", md5sum, expected)
	}
}

func TestCalculateFileMD5_NonExistentFile(t *testing.T) {
	_, err := CalculateFileMD5("/non/existent/file")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	
	// File doesn't exist yet
	if FileExists(testFile) {
		t.Error("FileExists returned true for non-existent file")
	}

	// Create the file
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// File exists now
	if !FileExists(testFile) {
		t.Error("FileExists returned false for existing file")
	}
}

func TestCompareFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create two identical files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	content := []byte("Same content")
	
	if err := os.WriteFile(file1, content, 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, content, 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Compare identical files
	same, err := CompareFiles(file1, file2)
	if err != nil {
		t.Fatalf("CompareFiles failed: %v", err)
	}
	if !same {
		t.Error("CompareFiles returned false for identical files")
	}

	// Create a different file
	file3 := filepath.Join(tempDir, "file3.txt")
	if err := os.WriteFile(file3, []byte("Different content"), 0644); err != nil {
		t.Fatalf("Failed to create file3: %v", err)
	}

	// Compare different files
	same, err = CompareFiles(file1, file3)
	if err != nil {
		t.Fatalf("CompareFiles failed: %v", err)
	}
	if same {
		t.Error("CompareFiles returned true for different files")
	}
}