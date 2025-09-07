package cmd

import (
	"fmt"
	"os"

	"github.com/doodleEsc/ctx-tool/internal/config"
	"github.com/doodleEsc/ctx-tool/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	lang          string
	configManager *config.Manager
	cfg           *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "ctx-tool",
	Short: "Manage Claude Code configurations",
	Long:  "ctx-tool is a CLI application for managing Claude Code configurations across projects.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initI18n(); err != nil {
			return err
		}
		if err := initConfig(); err != nil {
			return err
		}
		// Update command descriptions after i18n is initialized
		updateCommandDescriptions(cmd)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Initialize i18n early
	if err := i18n.Init(""); err != nil {
		// If i18n fails, continue with English defaults
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize i18n: %v\n", err)
	}

	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ~/.config/ctx-tool/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&lang, "lang", "l", "", "language (en, zh-Hans)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func initI18n() error {
	return i18n.Init(lang)
}

func initConfig() error {
	configManager = config.NewManager()
	if err := configManager.Load(cfgFile); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	cfg = configManager.GetConfig()

	// Update i18n language from config if not set via flag
	if lang == "" && cfg.I18n.Language != "" {
		i18n.SetLanguage(cfg.I18n.Language)
	}

	return nil
}

// updateCommandDescriptions updates all command descriptions after i18n is initialized
func updateCommandDescriptions(rootCmd *cobra.Command) {
	// Update root command
	rootCmd.Short = i18n.T(i18n.CmdRootShort)
	rootCmd.Long = i18n.T(i18n.CmdRootLong)

	// Update all subcommands
	for _, cmd := range rootCmd.Commands() {
		switch cmd.Name() {
		case "add":
			cmd.Short = i18n.T(i18n.CmdAddShort)
			cmd.Long = i18n.T(i18n.CmdAddLong)
			cmd.Example = i18n.T(i18n.CmdAddExample)
		case "remove":
			cmd.Short = i18n.T(i18n.CmdRemoveShort)
			cmd.Long = i18n.T(i18n.CmdRemoveLong)
			cmd.Example = i18n.T(i18n.CmdRemoveExample)
		}
	}
}
