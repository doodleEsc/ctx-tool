# PRP Agentic Coding Tool

## Goal

当前项目实现一个命令行工具，将github仓库`https://github.com/Wirasm/PRPs-agentic-eng`中的claude code自定义命令、subagents以及相关文件下载并拷贝到当前项目的根目录中的工具。

## 远程仓库文件树

```
➜ PRPs-agentic-eng (development) ✓ tree -a -I ".git|.python-version"
.
├── .claude
│   ├── agents
│   │   ├── codebase-analyst.md
│   │   └── library-researcher.md
│   └── commands
│       ├── code-quality
│       │   ├── refactor-simple.md
│       │   ├── review-general.md
│       │   └── review-staged-unstaged.md
│       ├── development
│       │   ├── create-pr.md
│       │   ├── debug-RCA.md
│       │   ├── new-dev-branch.md
│       │   ├── onboarding.md
│       │   ├── prime-core.md
│       │   └── smart-commit.md
│       ├── git-operations
│       │   ├── conflict-resolver-general.md
│       │   ├── conflict-resolver-specific.md
│       │   └── smart-resolver.md
│       ├── prp-commands
│       │   ├── api-contract-define.md
│       │   ├── prp-base-create.md
│       │   ├── prp-base-execute.md
│       │   ├── prp-planning-create.md
│       │   ├── prp-poc-create-parallel.md
│       │   ├── prp-poc-execute-parallel.md
│       │   ├── prp-spec-create.md
│       │   ├── prp-spec-execute.md
│       │   ├── prp-story-create.md
│       │   ├── prp-story-execute.md
│       │   ├── prp-task-create.md
│       │   ├── prp-task-execute.md
│       │   ├── prp-ts-create.md
│       │   ├── prp-ts-execute.md
│       │   └── task-list-init.md
│       ├── rapid-development
│       │   └── experimental
│       │       ├── create-base-prp-parallel.md
│       │       ├── create-planning-parallel.md
│       │       ├── hackathon-prp-parallel.md
│       │       ├── hackathon-research.md
│       │       ├── parallel-prp-creation.md
│       │       ├── prp-analyze-run.md
│       │       ├── prp-validate.md
│       │       └── user-story-rapid.md
│       └── typescript
│           ├── TS-create-base-prp.md
│           ├── TS-execute-base-prp.md
│           ├── TS-review-general.md
│           └── TS-review-staged-unstaged.md
├── .gitignore
├── claude_md_files
│   ├── CLAUDE-ASTRO.md
│   ├── CLAUDE-JAVA-GRADLE.md
│   ├── CLAUDE-JAVA-MAVEN.md
│   ├── CLAUDE-NEXTJS-15.md
│   ├── CLAUDE-NODE.md
│   ├── CLAUDE-PYTHON-BASIC.md
│   ├── CLAUDE-REACT.md
│   └── CLAUDE-RUST.md
├── CLAUDE.md
├── PRPs
│   ├── ai_docs
│   │   ├── build_with_claude_code.md
│   │   ├── cc_administration.md
│   │   ├── cc_cli.md
│   │   ├── cc_commands.md
│   │   ├── cc_containers.md
│   │   ├── cc_deployment.md
│   │   ├── cc_hooks.md
│   │   ├── cc_mcp.md
│   │   ├── cc_monitoring.md
│   │   ├── cc_settings.md
│   │   ├── cc_troubleshoot.md
│   │   ├── getting_started.md
│   │   ├── github_actions.md
│   │   ├── hooks.md
│   │   └── subagents.md
│   ├── example-from-workshop-mcp-crawl4ai-refactor-1.md
│   ├── pydantic-ai-prp-creation-agent-parallel.md
│   ├── README.md
│   ├── scripts
│   │   └── prp_runner.py
│   ├── STORY_WORKFLOW_GUIDE.md
│   └── templates
│       ├── prp_base_typescript.md
│       ├── prp_base.md
│       ├── prp_planning.md
│       ├── prp_poc_react.md
│       ├── prp_spec.md
│       ├── prp_story_task.md
│       └── prp_task.md
├── pyproject.toml
├── README-for-DUMMIES.md
└── README.md

16 directories, 81 files
```

## Requirements

- 使用spf13/cobra工具构建命令行工具，工具名称为'ctx-tool'，提供`add`以及`remove`子命令
- 提供配置文件，文件格式为yaml，可用选项为`repo`，表示从哪个远程仓库拷贝<https://github.com/Wirasm/PRPs-agentic-eng>
- `add`子命令使用说明如下：
  - ctx-tool add [--global | --project] [--all | <DIR_NAME>]：将`.cluade`目录以及`PRPs`目录或制定目录，拷贝到本机
    - `--local`选项表示将文件以及文件夹拷贝到当前项目目录中对应目录；默认为`--local`
    - `--global`表示将文件拷贝到$HOME/.claude目录。如果$HOME/.claude目录下存在agents或commands目录，**不能直接覆盖目录**，而是添加。
    - `--all`表示拷贝`.claude`目录和`PRPs`目录；如果指定目录，则选择拷贝指定目录
    - `add`拷贝文件前需要检目标目录中是否存在对应文件；不存在则直接拷贝；如果存在，则检查md5值，只拷贝md5值不同的文件。
    - `add`检测或拷贝文件之后，需要记录文件名称，用于后续`remove`指令删除文件
    - `add`指令检测文件时，不能将repo仓库中不存在的文件记录。否则可能会导致额外文件被记录或删除。
- `remove`指令根据`add`创建的记录文件删除对应文件.
  - ctx-tool remove
