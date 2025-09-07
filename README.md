# ctx-tool

A CLI application for managing Claude Code configurations across projects.

## Overview

ctx-tool helps you manage Claude Code configurations by syncing them from the [PRPs-agentic-eng repository](https://github.com/Wirasm/PRPs-agentic-eng) to your local projects or global configuration.

## Installation

```bash
go install github.com/doodleEsc/ctx-tool@latest
```

Or build from source:

```bash
git clone https://github.com/doodleEsc/ctx-tool
cd ctx-tool
go build -o ctx-tool
```

## Configuration

### XDG Base Directory Specification

ctx-tool follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) for configuration file management, providing better organization and cross-platform compatibility.

### Configuration File Locations

Configuration files are searched in the following priority order:

1. **XDG Config Directory** (highest priority):
   - Linux/Unix: `~/.config/ctx-tool/config.yaml`
   - macOS: `~/Library/Application Support/ctx-tool/config.yaml`
   - Windows: `%APPDATA%\ctx-tool\config.yaml`

2. **Legacy Home Directory**:
   - `~/.ctx-tool.yaml`

3. **Current Directory** (lowest priority):
   - `./.ctx-tool.yaml`

### Configuration File Format

```yaml
# ctx-tool Configuration File
version: "1.0"

# Repository configuration
repository:
  url: "https://github.com/Wirasm/PRPs-agentic-eng"
  branch: "development"  # Use "main" for stable version

# Tracking file configuration
tracking:
  file: ".ctx-tool-tracking.json"

# Allowed directories to sync
directories:
  allowed:
    - ".claude"
    - "PRPs"
    - "claude_md_files"

# Behavior configuration
behavior:
  backup_on_conflict: true  # Create .backup files when overwriting
  verify_md5: true          # Check MD5 before overwriting files
  clean_empty_dirs: true    # Remove empty directories on uninstall

# Internationalization configuration
i18n:
  language: ""              # "en" (English), "zh-Hans" (Simplified Chinese)
  locales_dir: ""          # Custom locales directory (optional)
```

### Migrating from Legacy Configuration

If you're using a legacy configuration file (`.ctx-tool.yaml` in your home directory or current directory), ctx-tool will automatically detect it and suggest migration to the XDG-compliant location.

To migrate manually:

1. **Create the XDG config directory**:
   ```bash
   # Linux/Unix
   mkdir -p ~/.config/ctx-tool
   
   # macOS
   mkdir -p ~/Library/Application\ Support/ctx-tool
   
   # Windows
   mkdir %APPDATA%\ctx-tool
   ```

2. **Copy your existing configuration**:
   ```bash
   # Linux/Unix
   cp ~/.ctx-tool.yaml ~/.config/ctx-tool/config.yaml
   
   # macOS
   cp ~/.ctx-tool.yaml ~/Library/Application\ Support/ctx-tool/config.yaml
   
   # Windows
   copy %USERPROFILE%\.ctx-tool.yaml %APPDATA%\ctx-tool\config.yaml
   ```

3. **Remove the old configuration** (optional):
   ```bash
   rm ~/.ctx-tool.yaml
   ```

### Environment Variables

Configuration values can be overridden using environment variables with the `CTX_TOOL_` prefix:

```bash
export CTX_TOOL_REPOSITORY_URL="https://github.com/your-org/your-repo"
export CTX_TOOL_REPOSITORY_BRANCH="main"
export CTX_TOOL_I18N_LANGUAGE="zh-Hans"
```

Environment variables have the highest priority and will override any configuration file settings.

## Usage

### Add Configurations

Install all available configurations:

```bash
ctx-tool add --all
```

Install specific directories:

```bash
ctx-tool add .claude PRPs
```

Install globally to your home directory:

```bash
ctx-tool add --all --global
```

### Remove Configurations

Remove previously installed configurations:

```bash
ctx-tool remove
```

Remove with detailed output:

```bash
ctx-tool remove --verbose
```

### Command Options

- `-c, --config`: Specify a custom configuration file path
- `-l, --lang`: Set language (en, zh-Hans)  
- `-v, --verbose`: Enable verbose output

### Examples

```bash
# Use custom config file
ctx-tool --config /path/to/custom/config.yaml add --all

# Set language to Simplified Chinese
ctx-tool --lang zh-Hans add .claude

# Install with verbose output
ctx-tool --verbose add --all
```

## Configuration Priority

Settings are applied in the following order (highest to lowest priority):

1. **Environment Variables** (`CTX_TOOL_*`)
2. **XDG Configuration File** (`~/.config/ctx-tool/config.yaml`)
3. **Legacy Configuration Files** (`~/.ctx-tool.yaml` or `./.ctx-tool.yaml`)
4. **Default Values**

## Tracking

ctx-tool maintains a tracking file (`.ctx-tool-tracking.json`) to keep track of installed files. This enables:

- Conflict detection and backup creation
- Clean removal of installed files
- MD5 verification to prevent unnecessary overwrites

## Cross-Platform Support

ctx-tool works on:

- Linux (follows XDG specification)
- macOS (uses Application Support directory)
- Windows (uses APPDATA directory)

All paths are automatically adjusted for the target platform.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions, issues, or feature requests, please visit the [GitHub repository](https://github.com/doodleEsc/ctx-tool).