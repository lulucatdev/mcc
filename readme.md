# mcc - Multiple Claude Code Accounts

> *"One subscription is never enough."* — Every power user, probably

A CLI tool to manage multiple Claude Code accounts. Switch between your personal, work, and "totally-not-my-alt" accounts with a single command.

[中文说明](#中文说明)

## The Problem

You have 2 Claude subscriptions. Maybe 3. You're not addicted, you just have... *needs*.

But Claude Code only lets you stay logged into one account at a time. Every switch means:
1. Log out
2. Log in
3. Wait
4. Forget which account you're on
5. Accidentally use your work quota for personal projects
6. Regret

## The Solution

```bash
mcc              # Back to default, like nothing happened
mcc run work     # Work mode activated
mcc run personal # Time for side projects
```

That's it. No logout. No login. Just vibes.

## How it Works

```
~/.mcc/
├── profiles/
│   ├── default/    ← Your main account lives here
│   ├── work/       ← Work account
│   └── chaotic/    ← We don't talk about this one
└── current → profiles/default  ← Magic symlink
```

Claude Code reads from `CLAUDE_CONFIG_DIR`. We point it to `~/.mcc/current`. We control what `current` points to. Simple.

## Installation

```bash
git clone https://github.com/anthropics/mcc.git
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

- Go 1.19+
- macOS or Linux
- Claude Code CLI installed
- Multiple subscriptions (optional but... why else are you here?)

## Contributing

PRs welcome. Issues welcome. Star this repo if it saved you from subscription juggling hell.

## License

MIT License - Lulucat Innovations

*Because managing AI accounts shouldn't require AI.*

---

# 中文说明

> *"一个订阅永远不够。"* — 每个重度用户

一个管理多个 Claude Code 账号的命令行工具。在个人账号、工作账号之间一键切换。

## 问题

你有 2 个 Claude 订阅。也许 3 个。你没有上瘾，你只是有... *需求*。

但 Claude Code 一次只能登录一个账号。每次切换都意味着：
1. 登出
2. 登录
3. 等待
4. 忘记当前是哪个账号
5. 不小心用工作额度做私人项目
6. 后悔

## 解决方案

```bash
mcc              # 回到默认账号
mcc run work     # 工作模式启动
mcc run personal # 摸鱼时间到
```

就这样。不用登出。不用登录。优雅。

## 工作原理

```
~/.mcc/
├── profiles/
│   ├── default/    ← 主账号
│   ├── work/       ← 工作账号
│   └── chaotic/    ← 不可描述
└── current → profiles/default  ← 神奇的软链接
```

Claude Code 读取 `CLAUDE_CONFIG_DIR`。我们把它指向 `~/.mcc/current`。我们控制 `current` 指向谁。就这么简单。

## 安装

```bash
git clone https://github.com/anthropics/mcc.git
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

- Go 1.19+
- macOS 或 Linux
- 已安装 Claude Code CLI
- 多个订阅（可选，但... 不然你来这干嘛？）

## 贡献

欢迎 PR。欢迎 Issue。如果这个工具拯救了你，给个 Star。

## 许可证

MIT License - Lulucat Innovations

*因为管理 AI 账号不应该需要 AI。*
