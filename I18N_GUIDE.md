# 国际化使用指南 / Internationalization Guide

## 中文 / Chinese

### 支持的语言
- **英文** (en) - 默认语言
- **简体中文** (zh-Hans) - 完整支持

### 使用方法

#### 1. 命令行参数
```bash
# 使用中文界面
ctx-tool --lang zh-Hans command

# 示例
ctx-tool --lang zh-Hans add --help
```

#### 2. 环境变量
```bash
# 设置环境变量
export CTX_TOOL_LANG=zh-Hans
ctx-tool command

# 或者临时使用
CTX_TOOL_LANG=zh-Hans ctx-tool add --all
```

#### 3. 配置文件
在 `.ctx-tool.yaml` 中设置：
```yaml
i18n:
  language: "zh-Hans"  # 强制使用简体中文
```

### 语言检测优先级
1. `--lang` 命令行参数（最高优先级）
2. 配置文件中的 `i18n.language` 设置
3. `CTX_TOOL_LANG` 环境变量
4. `LANG` 环境变量（取前两位）
5. 英文默认值

---

## English

### Supported Languages
- **English** (en) - Default language
- **Simplified Chinese** (zh-Hans) - Full support

### Usage

#### 1. Command Line Flag
```bash
# Use Chinese interface
ctx-tool --lang zh-Hans command

# Example
ctx-tool --lang zh-Hans add --help
```

#### 2. Environment Variable
```bash
# Set environment variable
export CTX_TOOL_LANG=zh-Hans
ctx-tool command

# Or use temporarily
CTX_TOOL_LANG=zh-Hans ctx-tool add --all
```

#### 3. Configuration File
Set in `.ctx-tool.yaml`:
```yaml
i18n:
  language: "zh-Hans"  # Force Simplified Chinese
```

### Language Detection Priority
1. `--lang` command line flag (highest priority)
2. `i18n.language` setting in config file
3. `CTX_TOOL_LANG` environment variable
4. `LANG` environment variable (first two characters)
5. English default

## Configuration Examples

### Default Configuration (Auto-detect)
```yaml
# .ctx-tool.yaml
i18n:
  language: ""          # Uses system language
  locales_dir: ""       # Uses embedded files
```

### Chinese Configuration
```yaml
# .ctx-tool-zh.yaml
i18n:
  language: "zh-Hans"   # Force Simplified Chinese
  locales_dir: ""       # Uses embedded files
```

### Custom Locales
```yaml
# Custom translation files
i18n:
  language: "zh-Hans"
  locales_dir: "./custom-locales"  # Use custom translation directory
```

## Adding New Languages

To add support for new languages:

1. Create a new directory in `internal/i18n/locales/`
2. Add a `<language>.toml` file with translated messages
3. Follow the existing message key structure
4. Update the language mapping in `SetLanguage()` function

Example for Japanese (ja):
```bash
mkdir internal/i18n/locales/ja
cp internal/i18n/locales/en/en.toml internal/i18n/locales/ja/ja.toml
# Edit ja.toml with Japanese translations
```