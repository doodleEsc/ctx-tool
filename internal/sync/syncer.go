package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/doodleEsc/ctx-tool/internal/config"
	"github.com/doodleEsc/ctx-tool/internal/i18n"
	"github.com/doodleEsc/ctx-tool/internal/tracker"
)

type Syncer struct {
	sourceDir string
	targetDir string
	tracker   *tracker.Tracker
	config    *config.Config
}

func NewSyncer(sourceDir, targetDir string, tracker *tracker.Tracker, config *config.Config) *Syncer {
	return &Syncer{
		sourceDir: sourceDir,
		targetDir: targetDir,
		tracker:   tracker,
		config:    config,
	}
}

// SyncDirectory syncs an entire directory from source to target
func (s *Syncer) SyncDirectory(dirName string) error {
	sourcePath := filepath.Join(s.sourceDir, dirName)

	// Check if source directory exists
	info, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %s not found in repository", dirName)
		}
		return fmt.Errorf("stat source directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dirName)
	}

	// Check if this directory is allowed
	if !s.isAllowedDirectory(dirName) {
		return fmt.Errorf("directory %s is not in allowed list", dirName)
	}

	fmt.Printf("%s\n", i18n.Tf(i18n.MsgSyncingDir, map[string]interface{}{"Dir": dirName}))

	// Walk through the directory and sync files
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files and .git
		name := filepath.Base(path)
		if strings.HasPrefix(name, ".") {
			return nil
		}

		// Calculate relative path from source directory
		relPath, err := filepath.Rel(s.sourceDir, path)
		if err != nil {
			return fmt.Errorf("calculate relative path: %w", err)
		}

		// Sync the file
		return s.SyncFile(relPath)
	})
}

// SyncFile syncs a single file from source to target
func (s *Syncer) SyncFile(relPath string) error {
	sourcePath := filepath.Join(s.sourceDir, relPath)
	targetPath := filepath.Join(s.targetDir, relPath)

	// Check if target file exists
	if FileExists(targetPath) && s.config.Behavior.VerifyMD5 {
		// Compare MD5 checksums
		sourceMD5, err := CalculateFileMD5(sourcePath)
		if err != nil {
			return fmt.Errorf("calculate source MD5: %w", err)
		}

		targetMD5, err := CalculateFileMD5(targetPath)
		if err != nil {
			return fmt.Errorf("calculate target MD5: %w", err)
		}

		if sourceMD5 == targetMD5 {
			fmt.Printf("  %s\n", i18n.Tf(i18n.MsgSkipIdentical, map[string]interface{}{"File": relPath}))
			// Still track the file even if skipped
			s.tracker.RecordFile(relPath, targetPath, s.sourceDir)
			return nil
		}

		// Files are different, backup if configured
		if s.config.Behavior.BackupOnConflict {
			backupPath := targetPath + ".backup"
			if err := s.copyFile(targetPath, backupPath); err != nil {
				return fmt.Errorf("backup file: %w", err)
			}
			fmt.Printf("  %s\n", i18n.Tf(i18n.MsgBackedUp, map[string]interface{}{"Original": relPath, "Backup": filepath.Base(backupPath)}))
		}
	}

	// Create target directory if needed
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Copy the file
	if err := s.copyFile(sourcePath, targetPath); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	// Track the installed file
	if err := s.tracker.RecordFile(relPath, targetPath, s.sourceDir); err != nil {
		return fmt.Errorf("track file: %w", err)
	}

	fmt.Printf("  %s\n", i18n.Tf(i18n.MsgInstalled, map[string]interface{}{"File": relPath}))
	return nil
}

// SyncAll syncs all allowed directories
func (s *Syncer) SyncAll() error {
	for _, dir := range s.config.Directories.Allowed {
		// Check if directory exists in source
		sourcePath := filepath.Join(s.sourceDir, dir)
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("%s\n", i18n.Tf(i18n.MsgWarningDirNotFound, map[string]interface{}{"Dir": dir}))
			continue
		}

		if err := s.SyncDirectory(dir); err != nil {
			return fmt.Errorf("sync directory %s: %w", dir, err)
		}
	}
	return nil
}

// copyFile copies a file from source to destination
func (s *Syncer) copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	// Preserve file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("set permissions: %w", err)
	}

	return nil
}

// isAllowedDirectory checks if a directory is in the allowed list
func (s *Syncer) isAllowedDirectory(dir string) bool {
	for _, allowed := range s.config.Directories.Allowed {
		if allowed == dir {
			return true
		}
	}
	return false
}
