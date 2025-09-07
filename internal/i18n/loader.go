package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// LoadFromDirectory loads locale files from a directory (for development/testing)
func LoadFromDirectory(dir string) error {
	if bundle == nil {
		return fmt.Errorf("i18n not initialized")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			langDir := filepath.Join(dir, entry.Name())
			files, err := os.ReadDir(langDir)
			if err != nil {
				continue
			}

			for _, file := range files {
				if !file.IsDir() {
					filePath := filepath.Join(langDir, file.Name())
					if _, err := bundle.LoadMessageFile(filePath); err != nil {
						return fmt.Errorf("failed to load %s: %w", filePath, err)
					}
				}
			}
		}
	}

	return nil
}

// LoadMessage loads a single message for the current language
func LoadMessage(msg *i18n.Message) {
	if bundle != nil {
		bundle.AddMessages(language.MustParse(currentLang), msg)
	}
}

// LoadMessages loads multiple messages for the current language
func LoadMessages(messages []*i18n.Message) {
	if bundle != nil {
		for _, msg := range messages {
			bundle.AddMessages(language.MustParse(currentLang), msg)
		}
	}
}