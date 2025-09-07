package i18n

// Message keys for commands
const (
	// Root command
	CmdRootShort = "cmd.root.short"
	CmdRootLong  = "cmd.root.long"

	// Add command
	CmdAddShort   = "cmd.add.short"
	CmdAddLong    = "cmd.add.long"
	CmdAddExample = "cmd.add.example"

	// Remove command
	CmdRemoveShort   = "cmd.remove.short"
	CmdRemoveLong    = "cmd.remove.long"
	CmdRemoveExample = "cmd.remove.example"
)

// Message keys for user interactions
const (
	// Add command messages
	MsgInstallationScope     = "msg.add.installation_scope"
	MsgTargetDirectory       = "msg.add.target_directory"
	MsgSyncingAll            = "msg.add.syncing_all"
	MsgSyncingDirectory      = "msg.add.syncing_directory"
	MsgInstallationComplete  = "msg.add.installation_complete"
	MsgTrackingFileSaved     = "msg.add.tracking_file_saved"
	MsgFilesInstalled        = "msg.add.files_installed"

	// Remove command messages
	MsgRemovalScope          = "msg.remove.removal_scope"
	MsgFoundTrackedFiles     = "msg.remove.found_tracked_files"
	MsgConfirmRemoval        = "msg.remove.confirm_removal"
	MsgRemovalCancelled      = "msg.remove.removal_cancelled"
	MsgRemovalComplete       = "msg.remove.removal_complete"
	MsgFilesRemoved          = "msg.remove.files_removed"
	MsgNoTrackedFiles        = "msg.remove.no_tracked_files"

	// Git messages
	MsgCloningRepository     = "msg.git.cloning_repository"
	MsgRepositoryCloned      = "msg.git.repository_cloned"
	MsgCloneSuccess          = "msg.git.clone_success"

	// Sync messages
	MsgSyncingDir            = "msg.sync.syncing_directory"
	MsgSkipIdentical         = "msg.sync.skip_identical"
	MsgBackedUp              = "msg.sync.backed_up"
	MsgInstalled             = "msg.sync.installed"
	MsgWarningDirNotFound    = "msg.sync.warning_dir_not_found"
)

// Error message keys
const (
	// Config errors
	ErrLoadConfig            = "err.config.load"
	ErrConfigNotFound        = "err.config.not_found"
	ErrInvalidConfig         = "err.config.invalid"

	// File operation errors
	ErrFileRead              = "err.file.read"
	ErrFileWrite             = "err.file.write"
	ErrFileCreate            = "err.file.create"
	ErrFileDelete            = "err.file.delete"
	ErrDirCreate             = "err.dir.create"
	ErrPermissionDenied      = "err.permission_denied"

	// Git errors
	ErrGitClone              = "err.git.clone"
	ErrGitOpen               = "err.git.open"
	ErrGitPull               = "err.git.pull"

	// Tracker errors
	ErrTrackerLoad           = "err.tracker.load"
	ErrTrackerSave           = "err.tracker.save"
	ErrTrackerNotFound       = "err.tracker.not_found"

	// Sync errors
	ErrSyncFailed            = "err.sync.failed"
	ErrBackupFailed          = "err.sync.backup_failed"

	// Command errors
	ErrInvalidScope          = "err.cmd.invalid_scope"
	ErrInvalidTarget         = "err.cmd.invalid_target"
	ErrMissingArgument       = "err.cmd.missing_argument"
)