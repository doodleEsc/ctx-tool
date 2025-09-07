package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	AppName            = "ctx-tool"
	XDGConfigFileName  = "config.yaml"
	LegacyConfigFile1  = ".ctx-tool.yaml"  // Home directory
	LegacyConfigFile2  = ".ctx-tool.yaml"  // Current directory
)

// GetConfigPaths returns configuration file paths in priority order:
// 1. XDG config directory (~/.config/ctx-tool/config.yaml)
// 2. Legacy home directory (~/.ctx-tool.yaml)  
// 3. Legacy current directory (./.ctx-tool.yaml)
func GetConfigPaths() []string {
	paths := make([]string, 0, 3)
	
	// XDG config path (highest priority)
	xdgConfigPath := filepath.Join(xdg.ConfigHome, AppName, XDGConfigFileName)
	paths = append(paths, xdgConfigPath)
	
	// Legacy home directory config
	if homeDir, err := os.UserHomeDir(); err == nil {
		legacyHomePath := filepath.Join(homeDir, LegacyConfigFile1)
		paths = append(paths, legacyHomePath)
	}
	
	// Legacy current directory config (lowest priority)
	if cwd, err := os.Getwd(); err == nil {
		legacyCwdPath := filepath.Join(cwd, LegacyConfigFile2)
		paths = append(paths, legacyCwdPath)
	}
	
	return paths
}

// GetXDGConfigPath returns the XDG-compliant configuration file path
func GetXDGConfigPath() string {
	return filepath.Join(xdg.ConfigHome, AppName, XDGConfigFileName)
}

// GetXDGConfigDir returns the XDG-compliant configuration directory
func GetXDGConfigDir() string {
	return filepath.Join(xdg.ConfigHome, AppName)
}

// EnsureConfigFile creates the XDG configuration file with default content if it doesn't exist
func EnsureConfigFile() error {
	configPath := GetXDGConfigPath()
	configDir := GetXDGConfigDir()
	
	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // File exists, nothing to do
	}
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}
	
	// Create default configuration content
	defaultConfig := getDefaultConfigContent()
	
	// Write default configuration file
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to create config file %s: %w", configPath, err)
	}
	
	return nil
}

// FileExists checks if a file exists and is readable
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// HasLegacyConfig checks if any legacy configuration files exist
func HasLegacyConfig() []string {
	var legacyFiles []string
	
	// Check home directory legacy config
	if homeDir, err := os.UserHomeDir(); err == nil {
		legacyHomePath := filepath.Join(homeDir, LegacyConfigFile1)
		if FileExists(legacyHomePath) {
			legacyFiles = append(legacyFiles, legacyHomePath)
		}
	}
	
	// Check current directory legacy config
	if cwd, err := os.Getwd(); err == nil {
		legacyCwdPath := filepath.Join(cwd, LegacyConfigFile2)
		if FileExists(legacyCwdPath) {
			legacyFiles = append(legacyFiles, legacyCwdPath)
		}
	}
	
	return legacyFiles
}

// getDefaultConfigContent returns the default configuration file content
func getDefaultConfigContent() string {
	return `# ctx-tool Configuration File
# XDG Base Directory compliant configuration
# Location: ~/.config/ctx-tool/config.yaml

version: "1.0"

# Repository configuration
repository:
  url: "https://github.com/Wirasm/PRPs-agentic-eng"
  branch: "development"  # Use "main" for stable version

# Tracking file configuration  
tracking:
  file: ".ctx-tool-tracking.json"

# Allowed directories to sync
directories:
  allowed:
    - ".claude"
    - "PRPs"
    - "claude_md_files"

# Behavior configuration
behavior:
  backup_on_conflict: true  # Create .backup files when overwriting
  verify_md5: true          # Check MD5 before overwriting files
  clean_empty_dirs: true    # Remove empty directories on uninstall

# Internationalization configuration
i18n:
  language: ""              # Language preference: "en" (English), "zh-Hans" (Simplified Chinese)
                           # Leave empty to use system language (CTX_TOOL_LANG or LANG env var)
  locales_dir: ""          # Custom locales directory path (leave empty to use embedded files)
                           # Use this to override built-in translations with custom ones
`
}