package ui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AttachView represents a live attach session view
type AttachView struct {
	viewport    viewport.Model
	width       int
	height      int
	sessionID   string
	agentID     string
	content     string
	ready       bool
	logStream   chan string
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewAttachView creates a new attach view
func NewAttachView() *AttachView {
	return &AttachView{
		width:     80,
		height:    20,
		content:   "",
		logStream: make(chan string, 100),
	}
}

// SetSession sets the session to attach to
func (a *AttachView) SetSession(sessionID, agentID string) {
	a.sessionID = sessionID
	a.agentID = agentID
}

// Start begins streaming logs
func (a *AttachView) Start() tea.Cmd {
	a.ctx, a.cancel = context.WithCancel(context.Background())
	return a.streamLogs()
}

// Stop cancels the log streaming
func (a *AttachView) Stop() {
	if a.cancel != nil {
		a.cancel()
	}
}

// SetDimensions sets the view dimensions
func (a *AttachView) SetDimensions(width, height int) {
	a.width = width
	a.height = height
	if a.ready {
		a.viewport.Width = width - 4
		a.viewport.Height = height - 6
	}
}

// streamLogs streams log data from the session
func (a *AttachView) streamLogs() tea.Cmd {
	return func() tea.Msg {
		// Safety check - if context is nil, we're not properly initialized
		if a == nil || a.ctx == nil {
			return nil
		}

		// For now, poll logs every second
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-a.ctx.Done():
				return nil
			case <-ticker.C:
				// Safety check before fetching logs
				if a == nil || a.ctx == nil || a.ctx.Err() != nil {
					return nil
				}
				// Try to get new logs
				logs, err := a.fetchLogs()
				if err == nil && logs != "" {
					return LogUpdateMsg{Content: logs}
				}
			}
		}
	}
}

// fetchLogs retrieves logs from the session
func (a *AttachView) fetchLogs() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get logs via openclaw CLI
	// Using --follow style if available, otherwise just get last 50 lines
	cmd := exec.CommandContext(ctx, "openclaw", "sessions", "logs", a.sessionID, "--lines", "50")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// LogUpdateMsg carries new log content
type LogUpdateMsg struct {
	Content string
}

// Init initializes the view
func (a *AttachView) Init() tea.Cmd {
	if !a.ready {
		a.viewport = viewport.New(a.width-4, a.height-6)
		a.viewport.SetContent(a.content)
		a.ready = true
	}
	return a.streamLogs()
}

// Update handles messages
func (a *AttachView) Update(msg tea.Msg) (*AttachView, tea.Cmd) {
	if !a.ready {
		a.viewport = viewport.New(a.width-4, a.height-6)
		a.viewport.SetContent(a.content)
		a.ready = true
	}

	switch m := msg.(type) {
	case LogUpdateMsg:
		// Append new content
		if m.Content != "" {
			a.content += m.Content
			a.viewport.SetContent(a.content)
			a.viewport.GotoBottom()
		}
		// Continue streaming
		return a, a.streamLogs()
	default:
		var cmd tea.Cmd
		a.viewport, cmd = a.viewport.Update(msg)
		return a, cmd
	}
}

// View renders the attach view
func (a *AttachView) View() string {
	if !a.ready {
		return "Initializing attach view..."
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("236")).
		Width(a.width - 2).
		Padding(0, 1)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Width(a.width - 2).
		Padding(0, 1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("235")).
		Width(a.width).
		Height(a.height)

	var sb strings.Builder

	// Title
	title := fmt.Sprintf(" Attach: %s (Live) ", a.agentID)
	sb.WriteString(titleStyle.Render(title))
	sb.WriteString("\n")

	// Viewport
	sb.WriteString(a.viewport.View())
	sb.WriteString("\n")

	// Footer
	footer := footerStyle.Render("Press q to exit | ↑/↓ or PgUp/PgDn to scroll")
	sb.WriteString(footer)

	return boxStyle.Render(sb.String())
}

// ScrollDown scrolls down
func (a *AttachView) ScrollDown() {
	if a.ready {
		a.viewport.LineDown(1)
	}
}

// ScrollUp scrolls up
func (a *AttachView) ScrollUp() {
	if a.ready {
		a.viewport.LineUp(1)
	}
}

// PageDown pages down
func (a *AttachView) PageDown() {
	if a.ready {
		a.viewport.ViewDown()
	}
}

// PageUp pages up
func (a *AttachView) PageUp() {
	if a.ready {
		a.viewport.ViewUp()
	}
}

// GotoBottom jumps to bottom
func (a *AttachView) GotoBottom() {
	if a.ready {
		a.viewport.GotoBottom()
	}
}
