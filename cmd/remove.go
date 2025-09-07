package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doodleEsc/ctx-tool/internal/i18n"
	"github.com/doodleEsc/ctx-tool/internal/tracker"
	"github.com/spf13/cobra"
)

var (
	forceFlag        bool
	removeGlobalFlag bool
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove tracked configurations",
	Long:    "Remove previously installed configurations based on the tracking file.",
	Example: "  ctx-tool remove\n  ctx-tool remove --global",
	RunE:    runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Local flags for remove command
	removeCmd.Flags().BoolVar(&forceFlag, "force", false, "Skip confirmation prompt")
	removeCmd.Flags().BoolVar(&removeGlobalFlag, "global", false, "Remove from global location")
}

func runRemove(cmd *cobra.Command, args []string) error {
	// Determine scope
	scope := "project"
	trackingFile := cfg.Tracking.File

	if removeGlobalFlag {
		scope = "global"
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}
		trackingFile = filepath.Join(homeDir, ".ctx-tool-tracking.json")
	}

	fmt.Printf("%s\n", i18n.Tf(i18n.MsgRemovalScope, map[string]interface{}{"Scope": scope}))
	fmt.Printf("Tracking file: %s\n", trackingFile)

	// Check if tracking file exists
	if _, err := os.Stat(trackingFile); os.IsNotExist(err) {
		return fmt.Errorf("%s", i18n.Tf(i18n.MsgNoTrackedFiles, map[string]interface{}{"Path": trackingFile}))
	}

	// Load tracker
	trackerInstance := tracker.NewTracker(trackingFile, scope, "")
	if err := trackerInstance.Load(); err != nil {
		return fmt.Errorf("load tracking data: %w", err)
	}

	// Get list of tracked files
	trackedFiles := trackerInstance.GetTrackedFiles()
	if len(trackedFiles) == 0 {
		fmt.Println("No tracked files found - nothing to remove")
		return nil
	}

	fmt.Printf("\n%s\n", i18n.Tn(i18n.MsgFoundTrackedFiles, len(trackedFiles), map[string]interface{}{"Count": len(trackedFiles)}))
	for _, file := range trackedFiles {
		fmt.Printf("  - %s\n", file)
	}

	// Ask for confirmation unless --force is used
	if !forceFlag {
		fmt.Printf("\n%s", i18n.T(i18n.MsgConfirmRemoval))
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read user input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" && response != "y" {
			fmt.Println(i18n.T(i18n.MsgRemovalCancelled))
			return nil
		}
	}

	// Remove files
	basePath := trackerInstance.Installation.BasePath
	removedCount := 0
	failedCount := 0
	directories := make(map[string]bool)

	for _, relPath := range trackedFiles {
		fullPath := filepath.Join(basePath, relPath)

		// Track parent directories for cleanup
		dir := filepath.Dir(fullPath)
		directories[dir] = true

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Printf("  Skip %s (already removed)\n", relPath)
			continue
		}

		// Remove the file
		if err := os.Remove(fullPath); err != nil {
			fmt.Printf("  ❌ Failed to remove %s: %v\n", relPath, err)
			failedCount++
			continue
		}

		fmt.Printf("  Removed %s\n", relPath)
		removedCount++

		// Update tracker
		trackerInstance.RemoveFile(relPath)
	}

	// Clean up empty directories if configured
	if cfg.Behavior.CleanEmptyDirs {
		fmt.Println("\nCleaning up empty directories...")
		for dir := range directories {
			// Don't remove base directories like $HOME/.claude
			if dir == basePath {
				continue
			}

			// Check if directory is empty
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}

			if len(entries) == 0 {
				if err := os.Remove(dir); err == nil {
					fmt.Printf("  Removed empty directory: %s\n", dir)
				}
			}
		}
	}

	// Save updated tracking file or remove it if empty
	if len(trackerInstance.GetTrackedFiles()) == 0 {
		// No more tracked files, remove the tracking file
		if err := os.Remove(trackingFile); err != nil {
			fmt.Printf("Warning: Failed to remove tracking file: %v\n", err)
		} else {
			fmt.Printf("\nRemoved tracking file (no files left to track)\n")
		}
	} else {
		// Save updated tracking data
		if err := trackerInstance.Save(); err != nil {
			return fmt.Errorf("save tracking data: %w", err)
		}
	}

	// Summary
	fmt.Printf("\n%s\n", i18n.T(i18n.MsgRemovalComplete))
	fmt.Printf("%s\n", i18n.Tn(i18n.MsgFilesRemoved, removedCount, map[string]interface{}{"Count": removedCount}))
	if failedCount > 0 {
		fmt.Printf("Files failed: %d\n", failedCount)
	}

	return nil
}
