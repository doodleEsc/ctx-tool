package sync

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CalculateFileMD5 calculates the MD5 checksum of a file
func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CalculateDirectoryChecksums calculates MD5 checksums for all files in a directory
func CalculateDirectoryChecksums(dir string) (map[string]string, error) {
	checksums := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files and git files
		name := filepath.Base(path)
		if strings.HasPrefix(name, ".") {
			return nil
		}

		// Skip .git directory contents
		if strings.Contains(path, ".git/") || strings.Contains(path, ".git\\") {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("calculate relative path: %w", err)
		}

		// Calculate MD5
		checksum, err := CalculateFileMD5(path)
		if err != nil {
			return fmt.Errorf("calculate MD5 for %s: %w", relPath, err)
		}

		checksums[relPath] = checksum
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walk directory: %w", err)
	}

	return checksums, nil
}

// CompareFiles compares two files using MD5 checksums
func CompareFiles(file1, file2 string) (bool, error) {
	// Check if both files exist
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", file1)
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", file2)
	}

	// Calculate checksums
	md5_1, err := CalculateFileMD5(file1)
	if err != nil {
		return false, fmt.Errorf("calculate MD5 for %s: %w", file1, err)
	}

	md5_2, err := CalculateFileMD5(file2)
	if err != nil {
		return false, fmt.Errorf("calculate MD5 for %s: %w", file2, err)
	}

	return md5_1 == md5_2, nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
