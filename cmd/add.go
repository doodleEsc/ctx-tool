package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/doodleEsc/ctx-tool/internal/git"
	"github.com/doodleEsc/ctx-tool/internal/i18n"
	"github.com/doodleEsc/ctx-tool/internal/sync"
	"github.com/doodleEsc/ctx-tool/internal/tracker"
	"github.com/spf13/cobra"
)

var (
	globalFlag  bool
	projectFlag bool
	allFlag     bool
)

var addCmd = &cobra.Command{
	Use:     "add [directories...]",
	Short:   "Add configurations from repository",
	Long:    "Add configurations from the PRPs-agentic-eng repository to your system.",
	Example: "  ctx-tool add --all\n  ctx-tool add prompts tools",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !allFlag {
			return errors.New("requires at least one directory or --all flag")
		}
		return nil
	},
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Local flags for add command
	addCmd.Flags().BoolVar(&globalFlag, "global", false, "Install globally to $HOME/.claude")
	addCmd.Flags().BoolVar(&projectFlag, "project", false, "Install to current project (default)")
	addCmd.Flags().BoolVar(&allFlag, "all", false, "Install all directories")

	// Mark flags as mutually exclusive
	addCmd.MarkFlagsMutuallyExclusive("global", "project")
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Determine scope and base path
	scope := "project"
	basePath := "."

	if globalFlag {
		scope = "global"
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}
		basePath = filepath.Join(homeDir, ".claude")

		// Create global .claude directory if it doesn't exist
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("create global .claude directory: %w", err)
		}
	}

	// If neither flag is set, default to project
	if !globalFlag && !projectFlag {
		projectFlag = true
	}

	fmt.Printf("%s\n", i18n.Tf(i18n.MsgInstallationScope, map[string]interface{}{"Scope": scope}))
	fmt.Printf("%s\n", i18n.Tf(i18n.MsgTargetDirectory, map[string]interface{}{"Target": basePath}))

	// Clone repository to temp directory
	gitClient := git.NewClient(cfg.Repository.URL, cfg.Repository.Branch)
	tempDir, err := gitClient.CloneToTemp()
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}
	defer func() {
		// Always clean up temp directory
		os.RemoveAll(tempDir)
		fmt.Printf("Cleaned up temporary directory\n")
	}()

	// Initialize tracker
	trackingFile := cfg.Tracking.File
	if scope == "global" {
		// Use global tracking file
		homeDir, _ := os.UserHomeDir()
		trackingFile = filepath.Join(homeDir, ".ctx-tool-tracking.json")
	}

	trackerInstance := tracker.NewTracker(trackingFile, scope, basePath)

	// Load existing tracking data
	if err := trackerInstance.Load(); err != nil {
		return fmt.Errorf("load tracking data: %w", err)
	}

	// Initialize syncer
	syncer := sync.NewSyncer(tempDir, basePath, trackerInstance, cfg)

	// Determine what to sync
	if allFlag {
		// Sync all allowed directories
		fmt.Println(i18n.T(i18n.MsgSyncingAll))
		if err := syncer.SyncAll(); err != nil {
			return fmt.Errorf("sync all directories: %w", err)
		}
	} else {
		// Sync specified directories
		for _, dir := range args {
			fmt.Printf("%s\n", i18n.Tf(i18n.MsgSyncingDirectory, map[string]interface{}{"Dir": dir}))
			if err := syncer.SyncDirectory(dir); err != nil {
				return fmt.Errorf("sync directory %s: %w", dir, err)
			}
		}
	}

	// Save tracking data
	if err := trackerInstance.Save(); err != nil {
		return fmt.Errorf("save tracking data: %w", err)
	}

	fmt.Printf("\n%s\n", i18n.T(i18n.MsgInstallationComplete))
	fmt.Printf("%s\n", i18n.Tf(i18n.MsgTrackingFileSaved, map[string]interface{}{"Path": trackingFile}))
	fmt.Printf("%s\n", i18n.Tn(i18n.MsgFilesInstalled, len(trackerInstance.GetTrackedFiles()), map[string]interface{}{"Count": len(trackerInstance.GetTrackedFiles())}))

	return nil
}
