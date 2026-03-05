package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LogViewer represents a log viewing component
type LogViewer struct {
	viewport viewport.Model
	width    int
	height   int
	content  string
	title    string
	ready    bool
}

// NewLogViewer creates a new log viewer
func NewLogViewer() *LogViewer {
	return &LogViewer{
		width:   80,
		height:  20,
		content: "",
		title:   "Log Viewer",
	}
}

// SetContent sets the log content
func (l *LogViewer) SetContent(content string) {
	l.content = content
	if l.ready {
		l.viewport.SetContent(content)
	}
}

// SetTitle sets the viewer title
func (l *LogViewer) SetTitle(title string) {
	l.title = title
}

// SetDimensions sets the viewer dimensions
func (l *LogViewer) SetDimensions(width, height int) {
	l.width = width
	l.height = height
	if l.ready {
		l.viewport.Width = width - 4
		l.viewport.Height = height - 6
	}
}

// Init initializes the viewport
func (l *LogViewer) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (l *LogViewer) Update(msg tea.Msg) (*LogViewer, tea.Cmd) {
	if !l.ready {
		l.viewport = viewport.New(l.width-4, l.height-6)
		l.viewport.SetContent(l.content)
		l.ready = true
		return l, nil
	}

	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

// View renders the log viewer
func (l *LogViewer) View() string {
	if !l.ready {
		return "Loading logs..."
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("236")).
		Width(l.width - 2).
		Padding(0, 1)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Width(l.width - 2).
		Padding(0, 1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("235")).
		Width(l.width).
		Height(l.height)

	var sb strings.Builder

	// Title
	sb.WriteString(titleStyle.Render(l.title))
	sb.WriteString("\n")

	// Viewport
	sb.WriteString(l.viewport.View())
	sb.WriteString("\n")

	// Footer
	footer := footerStyle.Render("Press q to exit | ↑/↓ or PgUp/PgDn to scroll")
	sb.WriteString(footer)

	return boxStyle.Render(sb.String())
}

// ScrollDown scrolls down
func (l *LogViewer) ScrollDown() {
	if l.ready {
		l.viewport.LineDown(1)
	}
}

// ScrollUp scrolls up
func (l *LogViewer) ScrollUp() {
	if l.ready {
		l.viewport.LineUp(1)
	}
}

// PageDown pages down
func (l *LogViewer) PageDown() {
	if l.ready {
		l.viewport.ViewDown()
	}
}

// PageUp pages up
func (l *LogViewer) PageUp() {
	if l.ready {
		l.viewport.ViewUp()
	}
}

// GotoBottom jumps to bottom
func (l *LogViewer) GotoBottom() {
	if l.ready {
		l.viewport.GotoBottom()
	}
}
