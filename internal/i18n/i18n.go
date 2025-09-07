package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle     *i18n.Bundle
	localizer  *i18n.Localizer
	currentLang string
)

//go:embed locales/*
var localesFS embed.FS

// Init initializes the i18n system
func Init(lang string) error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load embedded locale files
	if err := loadEmbeddedLocales(); err != nil {
		return fmt.Errorf("failed to load embedded locales: %w", err)
	}

	// Set language preference
	SetLanguage(lang)
	
	return nil
}

// loadEmbeddedLocales loads all locale files from embedded filesystem
func loadEmbeddedLocales() error {
	entries, err := localesFS.ReadDir("locales")
	if err != nil {
		return fmt.Errorf("read locales directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			langDir := filepath.Join("locales", entry.Name())
			files, err := localesFS.ReadDir(langDir)
			if err != nil {
				fmt.Printf("Warning: failed to read directory %s: %v\n", langDir, err)
				continue
			}

			for _, file := range files {
				if !file.IsDir() {
					filePath := filepath.Join(langDir, file.Name())
					data, err := localesFS.ReadFile(filePath)
					if err != nil {
						fmt.Printf("Warning: failed to read file %s: %v\n", filePath, err)
						continue
					}

					// Loading messages silently
					_, err = bundle.ParseMessageFileBytes(data, filePath)
					if err != nil {
						return fmt.Errorf("failed to parse %s: %w", filePath, err)
					}
				}
			}
		}
	}

	return nil
}

// SetLanguage sets the current language
func SetLanguage(lang string) {
	if lang == "" {
		// Try to get language from environment
		lang = os.Getenv("CTX_TOOL_LANG")
		if lang == "" {
			lang = os.Getenv("LANG")
			if len(lang) > 2 {
				lang = lang[:2]
			}
		}
		if lang == "" {
			lang = "en"
		}
	}

	// Map language codes
	switch lang {
	case "zh", "zh-CN", "zh-Hans", "zh_CN":
		lang = "zh-Hans"
	case "en", "en-US", "en_US":
		lang = "en"
	default:
		lang = "en"
	}

	currentLang = lang
	localizer = i18n.NewLocalizer(bundle, lang)
}

// GetLanguage returns the current language
func GetLanguage() string {
	return currentLang
}

// T translates a message with the given ID
func T(messageID string) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID // Fallback to message ID if translation not found
	}
	return msg
}

// Tf translates a message with format parameters
func Tf(messageID string, data map[string]interface{}) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID // Fallback to message ID if translation not found
	}
	return msg
}

// Tn translates a message with pluralization
func Tn(messageID string, count int, data map[string]interface{}) string {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["Count"] = count

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		PluralCount:  count,
	})
	if err != nil {
		return fmt.Sprintf("%s (%d)", messageID, count)
	}
	return msg
}

// Error creates a localized error
func Error(messageID string) error {
	return fmt.Errorf(T(messageID))
}

// Errorf creates a localized error with format parameters
func Errorf(messageID string, data map[string]interface{}) error {
	return fmt.Errorf(Tf(messageID, data))
}

// ListEmbeddedFiles lists all embedded files (for testing)
func ListEmbeddedFiles() ([]string, error) {
	entries, err := localesFS.ReadDir("locales")
	if err != nil {
		return nil, err
	}
	
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			langDir := filepath.Join("locales", entry.Name())
			subFiles, err := localesFS.ReadDir(langDir)
			if err != nil {
				continue
			}
			for _, file := range subFiles {
				if !file.IsDir() {
					files = append(files, filepath.Join(langDir, file.Name()))
				}
			}
		}
	}
	return files, nil
}