package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Manager struct {
	v      *viper.Viper
	config *Config
}

func NewManager() *Manager {
	v := viper.New()
	return &Manager{
		v:      v,
		config: &Config{},
	}
}

func (m *Manager) Load(configPath string) error {
	// Set defaults
	m.setDefaults()

	// Setup environment variables (highest priority)
	m.v.SetEnvPrefix("CTX_TOOL")
	m.v.AutomaticEnv()
	m.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Configure config file paths
	var configFound bool
	var foundPath string
	
	if configPath != "" {
		// Explicit config path provided via flag
		m.v.SetConfigFile(configPath)
		foundPath = configPath
		configFound = FileExists(configPath)
	} else {
		// Search for config files in priority order
		searchPaths := GetConfigPaths()
		
		for _, path := range searchPaths {
			if FileExists(path) {
				m.v.SetConfigFile(path)
				foundPath = path
				configFound = true
				break
			}
		}
		
		// If no config file found, try to create XDG default config
		if !configFound {
			if err := EnsureConfigFile(); err != nil {
				fmt.Printf("Warning: Failed to create default config file: %v\n", err)
			} else {
				xdgPath := GetXDGConfigPath()
				if FileExists(xdgPath) {
					m.v.SetConfigFile(xdgPath)
					foundPath = xdgPath
					configFound = true
					fmt.Printf("Created default configuration file: %s\n", xdgPath)
				}
			}
		}
	}

	// Read config file if found
	if configFound {
		if err := m.v.ReadInConfig(); err != nil {
			return fmt.Errorf("config file error reading %s: %w", foundPath, err)
		}
		fmt.Printf("Using config file: %s\n", foundPath)
		
		// Check for legacy config files and show migration hint
		if !isXDGPath(foundPath) {
			m.showMigrationHint(foundPath)
		}
	} else {
		fmt.Println("No config file found, using defaults")
	}

	// Unmarshal to struct
	if err := m.v.Unmarshal(m.config); err != nil {
		return fmt.Errorf("config unmarshal error: %w", err)
	}

	return nil
}

// isXDGPath checks if the given path is the XDG-compliant config path
func isXDGPath(path string) bool {
	return path == GetXDGConfigPath()
}

// showMigrationHint shows a helpful message about migrating to XDG config
func (m *Manager) showMigrationHint(currentPath string) {
	xdgPath := GetXDGConfigPath()
	fmt.Printf("\nNote: You're using a legacy config file location: %s\n", currentPath)
	fmt.Printf("Consider migrating to the XDG-compliant location: %s\n", xdgPath)
	fmt.Printf("You can copy your current config:\n")
	fmt.Printf("  mkdir -p %s && cp %s %s\n", GetXDGConfigDir(), currentPath, xdgPath)
	fmt.Println()
}

func (m *Manager) setDefaults() {
	m.v.SetDefault("version", "1.0")
	m.v.SetDefault("repository.url", "https://github.com/Wirasm/PRPs-agentic-eng")
	m.v.SetDefault("repository.branch", "development")
	m.v.SetDefault("tracking.file", ".ctx-tool-tracking.json")
	m.v.SetDefault("directories.allowed", []string{".claude", "PRPs", "claude_md_files"})
	m.v.SetDefault("behavior.backup_on_conflict", true)
	m.v.SetDefault("behavior.verify_md5", true)
	m.v.SetDefault("behavior.clean_empty_dirs", true)
}

func (m *Manager) GetConfig() *Config {
	return m.config
}
