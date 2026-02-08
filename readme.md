# mcc - Multiple Claude Code Accounts

> *"One subscription is never enough."* — Every power user, probably

A CLI tool to **run multiple Claude Code instances with different accounts simultaneously**.

[中文说明](#中文说明)

## What This Is (and Isn't)

**This is NOT a "switch" tool.** You don't switch your whole environment to another account.

**This IS a "run" tool.** Each terminal can run a Claude Code instance with a different account.

```bash
# Terminal 1
mcc run work      # Claude instance using work account

# Terminal 2
mcc run personal  # Claude instance using personal account

# Terminal 3
mcc               # Claude instance using default account
```

Three terminals. Three accounts. Running simultaneously. No conflicts.

## The Problem

You have multiple Claude subscriptions. But Claude Code ties to one account per config directory. Want to use your work account? Log out, log in, wait, configure...

## The Solution

```bash
mcc run work     # Launches claude with work account
mcc run personal # Launches claude with personal account
mcc              # Launches claude with default account
```

Each command spawns a Claude instance using that account's config. Open as many as you want.

## How it Works

```
~/.mcc/
├── profiles/
│   ├── default/    ← Account A's config
│   ├── work/       ← Account B's config
│   └── personal/   ← Account C's config
└── current → ...   ← Points to last used profile
```

When you run `mcc run work`, it:
1. Points the symlink to `profiles/work`
2. Launches Claude with `CLAUDE_CONFIG_DIR=~/.mcc/current`

Each terminal gets its own Claude process with the right account.

## Installation

### Download Pre-built Binary

Download the latest binary for your platform from [GitHub Releases](https://github.com/lulucatdev/mcc/releases), then:

```bash
chmod +x mcc-*
sudo mv mcc-* /usr/local/bin/mcc
```

### Build from Source

```bash
git clone https://github.com/lulucatdev/mcc.git
cd mcc
make setup
source ~/.zshrc
```

## Usage

```bash
mcc                  # Switch to default and launch claude
mcc run <name>       # Switch to profile and launch claude
mcc new <name>       # Create a new profile
mcc sync [name]      # Sync settings from ~/.claude (excludes credentials)
mcc status           # Show current status and profiles
mcc list             # List all profiles
mcc delete <name>    # Delete a profile
mcc help             # Show help
```

**Aliases:** `multicc` and `multi-claude-code` also work, if you're feeling verbose.

## Quick Start

```bash
# 1. Install
make setup && source ~/.zshrc

# 2. Check your setup
mcc status

# 3. Create a profile for your second account
mcc new work

# 4. Switch to it (will prompt for login)
mcc run work

# 5. Done. Switch back anytime:
mcc
```

## Roadmap

### v1.0 - Current
- [x] Multiple Claude Code account management
- [x] Profile switching with auto-launch
- [x] Settings sync (without credentials)
- [x] Works on macOS and Linux

### v2.0 - The Multiverse
> *What if `mcc` wasn't just for Claude?*

Imagine:
```bash
mcc run claude      # Claude Code
mcc run kimi        # Kimi (Moonshot AI)
mcc run copilot     # GitHub Copilot
mcc run cursor      # Cursor
mcc run gemini      # Google Gemini
```

One tool to rule them all. One tool to find them. One tool to bring them all, and in the terminal bind them.

**Coming eventually™** — or sooner if you open a PR.

### v3.0 - ???
- World domination
- Make coffee
- Achieve AGI
- Fix that one CSS bug

## Requirements

- macOS, Linux, or Windows
- Claude Code CLI installed
- Multiple subscriptions (optional but... why else are you here?)

## Release Workflow

Releases are automated via GitHub Actions. When a version tag is pushed:

```bash
git tag v0.2.0
git push origin v0.2.0
```

The workflow automatically:
1. Cross-compiles binaries for 6 platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64, windows/arm64)
2. Creates a GitHub Release with all binaries attached

## Contributing

PRs welcome. Issues welcome. Star this repo if it saved you from subscription juggling hell.

## License

MIT License - Lulucat Innovations

*Because managing AI accounts shouldn't require AI.*

---

# 中文说明

> *"一个订阅永远不够。"* — 每个重度用户

一个让你**同时运行多个不同账号的 Claude Code 实例**的命令行工具。

## 这是什么（和不是什么）

**这不是"切换"工具。** 不是把整个环境切换到另一个账号。

**这是"运行"工具。** 每个终端可以运行一个使用不同账号的 Claude Code 实例。

```bash
# 终端 1
mcc run work      # 使用工作账号的 Claude 实例

# 终端 2
mcc run personal  # 使用个人账号的 Claude 实例

# 终端 3
mcc               # 使用默认账号的 Claude 实例
```

三个终端。三个账号。同时运行。互不干扰。

## 问题

你有多个 Claude 订阅。但 Claude Code 的配置目录绑定一个账号。想用工作账号？登出、登录、等待、配置...

## 解决方案

```bash
mcc run work     # 启动使用工作账号的 claude
mcc run personal # 启动使用个人账号的 claude
mcc              # 启动使用默认账号的 claude
```

每个命令启动一个使用对应账号配置的 Claude 实例。想开几个开几个。

## 工作原理

```
~/.mcc/
├── profiles/
│   ├── default/    ← 账号 A 的配置
│   ├── work/       ← 账号 B 的配置
│   └── personal/   ← 账号 C 的配置
└── current → ...   ← 指向最后使用的配置
```

当你运行 `mcc run work` 时：
1. 把软链接指向 `profiles/work`
2. 用 `CLAUDE_CONFIG_DIR=~/.mcc/current` 启动 Claude

每个终端获得自己的 Claude 进程，使用正确的账号。

## 安装

### 下载预编译二进制

从 [GitHub Releases](https://github.com/lulucatdev/mcc/releases) 下载适合你平台的最新版本，然后：

```bash
chmod +x mcc-*
sudo mv mcc-* /usr/local/bin/mcc
```

### 从源码构建

```bash
git clone https://github.com/lulucatdev/mcc.git
cd mcc
make setup
source ~/.zshrc
```

## 使用方法

```bash
mcc                  # 切换到 default 并启动 claude
mcc run <名称>       # 切换到指定配置并启动 claude
mcc new <名称>       # 创建新配置
mcc sync [名称]      # 从 ~/.claude 同步设置（不包括登录凭证）
mcc status           # 显示当前状态和所有配置
mcc list             # 列出所有配置
mcc delete <名称>    # 删除配置
mcc help             # 显示帮助
```

**别名：** `multicc` 和 `multi-claude-code` 也可以用，如果你喜欢打字的话。

## 快速开始

```bash
# 1. 安装
make setup && source ~/.zshrc

# 2. 查看状态
mcc status

# 3. 为第二个账号创建配置
mcc new work

# 4. 切换过去（会提示登录）
mcc run work

# 5. 搞定。随时切回来：
mcc
```

## 路线图

### v1.0 - 当前版本
- [x] 多 Claude Code 账号管理
- [x] 配置切换 + 自动启动
- [x] 设置同步（不含凭证）
- [x] 支持 macOS 和 Linux

### v2.0 - 多元宇宙
> *如果 `mcc` 不只是给 Claude 用呢？*

想象一下：
```bash
mcc run claude      # Claude Code
mcc run kimi        # Kimi (月之暗面)
mcc run copilot     # GitHub Copilot
mcc run cursor      # Cursor
mcc run gemini      # Google Gemini
```

一个工具统治所有。一个工具找到所有。一个工具召唤所有，在终端中绑定所有。

**即将推出™** — 或者你来提 PR 会更快。

### v3.0 - ???
- 统治世界
- 自动泡咖啡
- 实现 AGI
- 修复那个 CSS bug

## 环境要求

- macOS、Linux 或 Windows
- 已安装 Claude Code CLI
- 多个订阅（可选，但... 不然你来这干嘛？）

## 发布工作流

发布通过 GitHub Actions 自动化。当推送版本标签时：

```bash
git tag v0.2.0
git push origin v0.2.0
```

工作流自动：
1. 交叉编译 6 个平台的二进制文件（darwin/amd64、darwin/arm64、linux/amd64、linux/arm64、windows/amd64、windows/arm64）
2. 创建 GitHub Release 并附带所有二进制文件

## 贡献

欢迎 PR。欢迎 Issue。如果这个工具拯救了你，给个 Star。

## 许可证

MIT License - Lulucat Innovations

*因为管理 AI 账号不应该需要 AI。*
