package tracker

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestTrackerSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	trackingFile := filepath.Join(tempDir, "tracking.json")
	
	// Create a new tracker
	tracker := NewTracker(trackingFile, "project", tempDir)
	
	// Create a test file to track
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Record the file
	if err := tracker.RecordFile("test.txt", testFile, "/source/repo"); err != nil {
		t.Fatalf("Failed to record file: %v", err)
	}
	
	// Save the tracker
	if err := tracker.Save(); err != nil {
		t.Fatalf("Failed to save tracker: %v", err)
	}
	
	// Create a new tracker and load the data
	newTracker := NewTracker(trackingFile, "", "")
	if err := newTracker.Load(); err != nil {
		t.Fatalf("Failed to load tracker: %v", err)
	}
	
	// Verify the loaded data
	if newTracker.Installation.Scope != "project" {
		t.Errorf("Scope mismatch: got %s, want project", newTracker.Installation.Scope)
	}
	
	files := newTracker.GetTrackedFiles()
	if len(files) != 1 {
		t.Errorf("File count mismatch: got %d, want 1", len(files))
	}
	
	if files[0] != "test.txt" {
		t.Errorf("File name mismatch: got %s, want test.txt", files[0])
	}
}

func TestTrackerRemoveFile(t *testing.T) {
	tempDir := t.TempDir()
	trackingFile := filepath.Join(tempDir, "tracking.json")
	
	tracker := NewTracker(trackingFile, "project", tempDir)
	
	// Create test files
	for i := 1; i <= 3; i++ {
		fileName := fmt.Sprintf("test%d.txt", i)
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := tracker.RecordFile(fileName, filePath, "/source"); err != nil {
			t.Fatalf("Failed to record file: %v", err)
		}
	}
	
	// Should have 3 files
	if len(tracker.GetTrackedFiles()) != 3 {
		t.Errorf("Expected 3 tracked files, got %d", len(tracker.GetTrackedFiles()))
	}
	
	// Remove one file
	tracker.RemoveFile("test2.txt")
	
	// Should have 2 files now
	files := tracker.GetTrackedFiles()
	if len(files) != 2 {
		t.Errorf("Expected 2 tracked files after removal, got %d", len(files))
	}
	
	// Verify the correct file was removed
	for _, file := range files {
		if file == "test2.txt" {
			t.Error("test2.txt should have been removed")
		}
	}
}

func TestTrackerLoadNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	trackingFile := filepath.Join(tempDir, "nonexistent.json")
	
	tracker := NewTracker(trackingFile, "project", tempDir)
	
	// Loading non-existent file should not error (treated as new)
	if err := tracker.Load(); err != nil {
		t.Errorf("Load should not error for non-existent file: %v", err)
	}
	
	// Files should be empty
	if len(tracker.GetTrackedFiles()) != 0 {
		t.Error("Expected empty file list for new tracker")
	}
}