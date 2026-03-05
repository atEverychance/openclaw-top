package ui

import (
	"fmt"
	"github.com/ateverychance/openclaw-top/pkg/models"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar displays agent stats at the bottom
type StatusBar struct {
	width   int
	stats   *models.AppStats
	height  int
	message string
	err     error
}

var (
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("235")).
			Height(1)

	agentCountStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	refreshStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75"))
)

// NewStatusBar creates a new status bar
func NewStatusBar() *StatusBar {
	return &StatusBar{
		width:  80,
		height: 1,
	}
}

// SetDimensions sets status bar dimensions
func (s *StatusBar) SetDimensions(width int) {
	s.width = width
}

// SetStats sets the app stats to display
func (s *StatusBar) SetStats(stats *models.AppStats) {
	s.stats = stats
}

// SetMessage sets a status message
func (s *StatusBar) SetMessage(msg string) {
	s.message = msg
	s.err = nil
}

// SetError sets an error to display
func (s *StatusBar) SetError(err error) {
	s.err = err
	s.message = ""
}

// View renders the status bar
func (s *StatusBar) View() string {
	// Show error if present
	if s.err != nil {
		errStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("196")).
			Width(s.width)
		msg := fmt.Sprintf(" Error: %s ", s.err.Error())
		return errStyle.Render(msg)
	}

	// Show message if present
	if s.message != "" {
		msgStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("28")).
			Width(s.width)
		return msgStyle.Render(fmt.Sprintf(" %s ", s.message))
	}

	if s.stats == nil {
		return statusStyle.Width(s.width).Render(" Initializing... ")
	}

	agentStr := fmt.Sprintf(" %d agents ", s.stats.TotalAgents)
	refreshStr := fmt.Sprintf(" Last refresh: %s ", formatRefreshTime(s.stats.LastRefresh))

	// Build the status string
	parts := []string{agentStr, refreshStr}
	status := strings.Join(parts, "|")

	// Pad to full width
	padding := s.width - lipgloss.Width(status) - 2
	if padding > 0 {
		status += strings.Repeat(" ", padding)
	}

	return statusStyle.Width(s.width).Render(status)
}

func formatRefreshTime(t time.Time) string {
	return t.Format("15:04:05")
}
