# openclaw-top 🎯

[![Go Report Card](https://goreportcard.com/badge/github.com/ateverychance/openclaw-top)](https://goreportcard.com/report/github.com/ateverychance/openclaw-top)
[![GoDoc](https://godoc.org/github.com/ateverychance/openclaw-top?status.svg)](https://godoc.org/github.com/ateverychance/openclaw-top)

> **htop for OpenClaw agents** — Real-time terminal UI for monitoring your AI agent fleet

![Demo](demo.svg)

## ✨ Features

- 🚀 **Real-time updates** — Auto-refresh every 2 seconds
- 🎨 **Beautiful terminal UI** — Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- ⌨️ **Vim-inspired navigation** — `j/k` to move, familiar keybindings
- 🟢 **Color-coded status** — Instantly see RUNNING, IDLE, ERROR states
- 📊 **Live metrics** — Token usage, runtime, model info at a glance
- 🎯 **Interactive** — Attach, kill, inspect agents directly from the UI

## 📦 Installation

### Using `go install` (recommended)

```bash
# Install the binary directly
go install github.com/ateverychance/openclaw-top/cmd/openclaw-top@latest
```

This will install `openclaw-top` to your `$GOPATH/bin` directory.

### Using the alias `octop`

To install with the shorter `octop` binary name:

```bash
# Build and install as octop
git clone https://github.com/ateverychance/openclaw-top.git
cd openclaw-top
make install-alias
```

Or manually:

```bash
go build -o octop ./cmd/openclaw-top
sudo mv octop /usr/local/bin/
```

### From source

```bash
git clone https://github.com/ateverychance/openclaw-top.git
cd openclaw-top
make build
```

### Using Makefile

```bash
make build          # Build binary with version info
make install        # Install to GOPATH/bin
make install-alias   # Install as 'octop'
make test           # Run tests
make clean          # Clean build artifacts
```

## 🚀 Usage

### Basic

```bash
# Run the TUI
openclaw-top

# Show version info
openclaw-top --version

# Or use the short alias (if installed)
octop
```

### Coming Soon (WIP)

```bash
# Auto-attach to a specific agent
openclaw-top --attach coder

# Kill a zombie agent
openclaw-top --kill scout

# Watch mode (stream updates)
openclaw-top --watch
```

## ⌨️ Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate up/down |
| `r` | Refresh data |
| `1-4` | Sort by column (Agent, Status, Runtime, Tokens) |
| `?` | Toggle help |
| `q` or `Ctrl+C` | Quit |

### Planned Keybindings

| Key | Action | Status |
|-----|--------|--------|
| `a` | Attach to selected agent | 🚧 WIP |
| `k` | Kill selected agent | 🚧 WIP |
| `l` | View logs | 🚧 WIP |
| `Enter` | Agent detail view | 🚧 WIP |

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│  Gateway Client (OpenClaw CLI)         │
│  └─→ Calls `openclaw sessions --json`  │
├─────────────────────────────────────────┤
│  Bubble Tea TUI                        │
│  ├─ Table (agent list)                 │
│  ├─ StatusBar (stats)                  │
│  └─ HelpOverlay (keybindings)          │
└─────────────────────────────────────────┘
```

## 🤝 Contributing

This project is part of the [OpenClaw](https://github.com/openclaw) ecosystem. PRs welcome!

## 📄 License

MIT License — see [LICENSE](LICENSE)

## 🙏 Ackowledgments

- Built with [Charm](https://charm.sh/)'s amazing TUI libraries
- Inspired by [htop](https://htop.dev/) and the TUI renaissance
- For the agent-driven development movement 🚀
