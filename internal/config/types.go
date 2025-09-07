package config

type Config struct {
	Version     string            `mapstructure:"version"`
	Repository  RepositoryConfig  `mapstructure:"repository"`
	Tracking    TrackingConfig    `mapstructure:"tracking"`
	Directories DirectoriesConfig `mapstructure:"directories"`
	Behavior    BehaviorConfig    `mapstructure:"behavior"`
	I18n        I18nConfig        `mapstructure:"i18n"`
}

type RepositoryConfig struct {
	URL    string `mapstructure:"url"`
	Branch string `mapstructure:"branch"`
}

type TrackingConfig struct {
	File string `mapstructure:"file"`
}

type DirectoriesConfig struct {
	Allowed []string `mapstructure:"allowed"`
}

type BehaviorConfig struct {
	BackupOnConflict bool `mapstructure:"backup_on_conflict"`
	VerifyMD5        bool `mapstructure:"verify_md5"`
	CleanEmptyDirs   bool `mapstructure:"clean_empty_dirs"`
}

type I18nConfig struct {
	Language   string `mapstructure:"language"`
	LocalesDir string `mapstructure:"locales_dir"`
}
