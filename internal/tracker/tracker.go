package tracker

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type Tracker struct {
	FilePath     string
	Installation *Installation
}

type Installation struct {
	Timestamp string      `json:"timestamp"`
	Scope     string      `json:"scope"`
	BasePath  string      `json:"base_path"`
	Files     []FileEntry `json:"files"`
}

type FileEntry struct {
	Path   string `json:"path"`
	MD5    string `json:"md5"`
	Size   int64  `json:"size"`
	Source string `json:"source"`
}

func NewTracker(filePath, scope, basePath string) *Tracker {
	return &Tracker{
		FilePath: filePath,
		Installation: &Installation{
			Timestamp: time.Now().Format(time.RFC3339),
			Scope:     scope,
			BasePath:  basePath,
			Files:     []FileEntry{},
		},
	}
}

func (t *Tracker) Load() error {
	data, err := os.ReadFile(t.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// New tracking file - this is normal
			return nil
		}
		return fmt.Errorf("read tracking file: %w", err)
	}

	if err := json.Unmarshal(data, t.Installation); err != nil {
		return fmt.Errorf("unmarshal tracking data: %w", err)
	}

	return nil
}

func (t *Tracker) Save() error {
	data, err := json.MarshalIndent(t.Installation, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tracking data: %w", err)
	}

	if err := os.WriteFile(t.FilePath, data, 0644); err != nil {
		return fmt.Errorf("write tracking file: %w", err)
	}

	return nil
}

func (t *Tracker) RecordFile(relPath, fullPath, source string) error {
	info, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Errorf("stat file %s: %w", fullPath, err)
	}

	md5sum, err := calculateFileMD5(fullPath)
	if err != nil {
		return fmt.Errorf("calculate MD5 for %s: %w", fullPath, err)
	}

	entry := FileEntry{
		Path:   relPath,
		MD5:    md5sum,
		Size:   info.Size(),
		Source: source,
	}

	// Check if file already tracked and update it
	found := false
	for i, existing := range t.Installation.Files {
		if existing.Path == relPath {
			t.Installation.Files[i] = entry
			found = true
			break
		}
	}

	if !found {
		t.Installation.Files = append(t.Installation.Files, entry)
	}

	return nil
}

func (t *Tracker) GetTrackedFiles() []string {
	var files []string
	for _, entry := range t.Installation.Files {
		files = append(files, entry.Path)
	}
	return files
}

func (t *Tracker) RemoveFile(relPath string) {
	var filtered []FileEntry
	for _, entry := range t.Installation.Files {
		if entry.Path != relPath {
			filtered = append(filtered, entry)
		}
	}
	t.Installation.Files = filtered
}

func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
