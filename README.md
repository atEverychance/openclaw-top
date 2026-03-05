# openclaw-top 🎯

[![Go Report Card](https://goreportcard.com/badge/github.com/ateverychance/openclaw-top)](https://goreportcard.com/report/github.com/ateverychance/openclaw-top)
[![GoDoc](https://godoc.org/github.com/ateverychance/openclaw-top?status.svg)](https://godoc.org/github.com/ateverychance/openclaw-top)

> **htop for OpenClaw agents** — Real-time terminal UI for monitoring your AI agent fleet

![Demo](demo.svg)

## ✨ Features

- 🚀 **Real-time updates** — Auto-refresh every 2 seconds
- 🎨 **Beautiful terminal UI** — Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- ⌨️ **Vim-inspired navigation** — `j/k` to move, familiar keybindings
- 🟢 **Visual status indicators** — Icons + progress bars + sparklines
- 📊 **Live metrics** — Token usage, runtime, model info at a glance
- 🎯 **Interactive control** — Attach, kill, view logs directly from the UI
- 🔧 **CLI modes** — Non-interactive modes for scripting

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
make install-alias  # Install as 'octop'
make test           # Run tests
make clean          # Clean build artifacts
```

## 🚀 Usage

### Interactive TUI (default)

```bash
# Run the TUI
openclaw-top

# Show version info
openclaw-top --version

# Or use the short alias (if installed)
octop
```

### CLI Modes

```bash
# Attach to a specific agent immediately
openclaw-top --attach coder

# Kill a zombie agent
openclaw-top --kill scout

# Watch mode (stream updates, no TUI)
openclaw-top --watch

# Custom refresh rate (seconds)
openclaw-top --refresh 5
```

## ⌨️ Keybindings

### Navigation

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate up/down |
| `r` | Refresh data |
| `1-4` | Sort by column (Agent, Status, Runtime, Tokens) |
| `?` | Toggle help |
| `q` or `Ctrl+C` | Quit |

### Actions

| Key | Action |
|-----|--------|
| `a` | **Attach** to selected agent (live log streaming) |
| `l` | **View logs** snapshot for selected agent |
| `k` | **Kill** selected agent (with confirmation) |

### In Attach/Log View

| Key | Action |
|-----|--------|
| `q` or `Esc` | Exit to table view |
| `↑/↓` or `j/k` | Scroll logs |
| `PgUp/PgDn` | Page up/down |
| `End` or `G` | Jump to bottom |

## 🎨 Visual Features

- **Status icons**: 🟢 RUNNING, 🔴 ERROR, 🟡 IDLE, 🔵 DONE
- **Progress bars**: Visual runtime indicators (█░░░░░░░░░)
- **Sparklines**: Token usage visualization (▁▃▅▇▁)
- **Color coding**: Instantly identify agent states
- **Bold selection**: Clear visual indicator for selected row

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│  Gateway Client (OpenClaw CLI)         │
│  ├─ FetchAll() - Get agent sessions    │
│  ├─ KillSession() - Kill an agent      │
│  └─ GetLogs() - Retrieve session logs  │
├─────────────────────────────────────────┤
│  Bubble Tea TUI                        │
│  ├─ Table (agent list with visuals)    │
│  ├─ StatusBar (stats & messages)       │
│  ├─ HelpOverlay (keybindings)          │
│  ├─ ConfirmModal (kill confirmation)   │
│  ├─ LogViewer (static logs)            │
│  └─ AttachView (live log streaming)    │
└─────────────────────────────────────────┘
```

## 📋 All GitHub Issues

| Issue | Feature | Status |
|-------|---------|--------|
| #1 | Attach mode ('a' key) | ✅ Complete |
| #2 | Kill agent ('k' key) | ✅ Complete |
| #3 | View logs ('l' key) | ✅ Complete |
| #4 | CLI args (--attach, --kill, --watch) | ✅ Complete |
| #5 | Visual polish (progress bars, sparklines) | ✅ Complete |
| #6 | go install distribution | ✅ Complete |

**MVP Status: 6/6 Complete (100%)** 🎉

## 🤝 Contributing

This project is part of the [OpenClaw](https://github.com/openclaw) ecosystem. PRs welcome!

## 📄 License

MIT License — see [LICENSE](LICENSE)

## 🙏 Acknowledgments

- Built with [Charm](https://charm.sh/)'s amazing TUI libraries
- Inspired by [htop](https://htop.dev/) and the TUI renaissance
- For the agent-driven development movement 🚀

---

*Built with agents, for agents.* 🤖
