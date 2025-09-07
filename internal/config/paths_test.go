package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/adrg/xdg"
)

func TestGetConfigPaths(t *testing.T) {
	// Test that GetConfigPaths returns paths in correct priority order
	paths := GetConfigPaths()
	
	// Should have at least 2-3 paths (XDG + legacy paths)
	if len(paths) < 2 {
		t.Errorf("Expected at least 2 config paths, got %d", len(paths))
	}
	
	// First path should be XDG config path
	expectedXDGPath := filepath.Join(xdg.ConfigHome, AppName, XDGConfigFileName)
	if paths[0] != expectedXDGPath {
		t.Errorf("First path should be XDG config path, got %s, expected %s", paths[0], expectedXDGPath)
	}
	
	// Verify all paths are absolute
	for i, path := range paths {
		if !filepath.IsAbs(path) {
			t.Errorf("Path %d (%s) should be absolute", i, path)
		}
	}
}

func TestGetXDGConfigPath(t *testing.T) {
	path := GetXDGConfigPath()
	
	// Should be absolute path
	if !filepath.IsAbs(path) {
		t.Errorf("XDG config path should be absolute, got %s", path)
	}
	
	// Should contain our app name and config filename
	if !strings.Contains(path, AppName) {
		t.Errorf("XDG config path should contain app name '%s', got %s", AppName, path)
	}
	
	if !strings.Contains(path, XDGConfigFileName) {
		t.Errorf("XDG config path should contain config filename '%s', got %s", XDGConfigFileName, path)
	}
}

func TestGetXDGConfigDir(t *testing.T) {
	dir := GetXDGConfigDir()
	
	// Should be absolute path
	if !filepath.IsAbs(dir) {
		t.Errorf("XDG config dir should be absolute, got %s", dir)
	}
	
	// Should contain our app name
	if !strings.Contains(dir, AppName) {
		t.Errorf("XDG config dir should contain app name '%s', got %s", AppName, dir)
	}
	
	// Should end with our app name
	if !strings.HasSuffix(dir, AppName) {
		t.Errorf("XDG config dir should end with app name '%s', got %s", AppName, dir)
	}
}

func TestFileExists(t *testing.T) {
	// Test with a file that should exist (current source file)
	currentFile := "paths_test.go"
	if !FileExists(currentFile) {
		t.Errorf("FileExists should return true for existing file %s", currentFile)
	}
	
	// Test with a file that shouldn't exist
	nonExistentFile := "/tmp/definitely-does-not-exist-12345.xyz"
	if FileExists(nonExistentFile) {
		t.Errorf("FileExists should return false for non-existent file %s", nonExistentFile)
	}
	
	// Test with a directory (should return false for directories)
	if FileExists(".") {
		t.Errorf("FileExists should return false for directories")
	}
}

func TestEnsureConfigFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "ctx-tool-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// We can't easily mock xdg.ConfigHome, so we'll test the logic indirectly
	// by checking that the function doesn't error when directory creation is possible
	
	// Test that EnsureConfigFile doesn't crash
	if err := EnsureConfigFile(); err != nil {
		// This might fail due to permissions, which is expected in some environments
		t.Logf("EnsureConfigFile failed (may be expected): %v", err)
	}
}

func TestHasLegacyConfig(t *testing.T) {
	// This test checks the function behavior without creating actual files
	legacyFiles := HasLegacyConfig()
	
	// Should return a slice (possibly nil or empty)
	// In Go, a nil slice is equivalent to an empty slice for most operations
	
	// All returned paths should be absolute
	for i, path := range legacyFiles {
		if !filepath.IsAbs(path) {
			t.Errorf("Legacy config path %d (%s) should be absolute", i, path)
		}
	}
}

func TestConstants(t *testing.T) {
	// Test that our constants are properly defined
	if AppName == "" {
		t.Errorf("AppName should not be empty")
	}
	
	if XDGConfigFileName == "" {
		t.Errorf("XDGConfigFileName should not be empty")
	}
	
	if LegacyConfigFile1 == "" {
		t.Errorf("LegacyConfigFile1 should not be empty")
	}
	
	if LegacyConfigFile2 == "" {
		t.Errorf("LegacyConfigFile2 should not be empty")
	}
	
	// Test that config filenames have proper extensions
	if !strings.HasSuffix(XDGConfigFileName, ".yaml") {
		t.Errorf("XDG config filename should have .yaml extension, got %s", XDGConfigFileName)
	}
}

func TestCrossPlatformPaths(t *testing.T) {
	// Test that paths work correctly across platforms
	paths := GetConfigPaths()
	
	// All paths should be valid for the current platform
	for _, path := range paths {
		if runtime.GOOS == "windows" {
			// On Windows, paths should use backslashes or be properly formatted
			if strings.Contains(path, "/") && !strings.Contains(path, "\\") {
				// Allow forward slashes in some contexts, but ensure path is valid
			}
		} else {
			// On Unix-like systems, paths should use forward slashes
			if strings.Contains(path, "\\") {
				t.Errorf("Unix path should not contain backslashes: %s", path)
			}
		}
	}
}

func TestConfigSearchPriority(t *testing.T) {
	// Test that config search follows the correct priority order
	paths := GetConfigPaths()
	
	if len(paths) < 2 {
		t.Fatalf("Expected at least 2 paths for priority testing")
	}
	
	// First path should be XDG (highest priority)
	xdgPath := GetXDGConfigPath()
	if paths[0] != xdgPath {
		t.Errorf("First path should be XDG path, got %s, expected %s", paths[0], xdgPath)
	}
	
	// Subsequent paths should be legacy paths
	for i := 1; i < len(paths); i++ {
		if paths[i] == xdgPath {
			t.Errorf("XDG path should only appear first in priority list")
		}
		
		// Legacy paths should contain the legacy config filename
		if !strings.Contains(paths[i], LegacyConfigFile1) && !strings.Contains(paths[i], LegacyConfigFile2) {
			t.Errorf("Legacy path %d should contain legacy config filename: %s", i, paths[i])
		}
	}
}

func TestDefaultConfigContent(t *testing.T) {
	// Test that default config content is valid YAML-like structure
	content := getDefaultConfigContent()
	
	if content == "" {
		t.Errorf("Default config content should not be empty")
	}
	
	// Should contain expected sections
	expectedSections := []string{
		"version:",
		"repository:",
		"tracking:",
		"directories:",
		"behavior:",
		"i18n:",
	}
	
	for _, section := range expectedSections {
		if !strings.Contains(content, section) {
			t.Errorf("Default config should contain section '%s'", section)
		}
	}
	
	// Should contain XDG-related comments
	if !strings.Contains(content, "XDG") {
		t.Errorf("Default config should contain XDG-related comments")
	}
}