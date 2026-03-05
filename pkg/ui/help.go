package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// HelpOverlay displays keyboard shortcuts
type HelpOverlay struct {
	width  int
	height int
}

var (
	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("255")).
			Padding(1, 2)

	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250"))
)

// NewHelpOverlay creates a new help overlay
func NewHelpOverlay() *HelpOverlay {
	return &HelpOverlay{
		width:  60,
		height: 20,
	}
}

// SetDimensions sets overlay dimensions
func (h *HelpOverlay) SetDimensions(width, height int) {
	h.width = width
	h.height = height
}

// View renders the help overlay
func (h *HelpOverlay) View() string {
	content := fmt.Sprintf(`%s

Navigation
  %s  %s
  %s  %s

Actions
  %s  %s
  %s  %s
  %s  %s

Sorting
  %s  %s

Press any key to close`,
		helpTitleStyle.Render(" OpenClaw Agent Monitor "),
		helpKeyStyle.Render("↑/k"), helpDescStyle.Render("Select next/previous agent"),
		helpKeyStyle.Render("↓/j"), helpDescStyle.Render("Select next/previous agent"),
		helpKeyStyle.Render("a"), helpDescStyle.Render("Attach to selected agent"),
		helpKeyStyle.Render("x"), helpDescStyle.Render("Kill selected agent"),
		helpKeyStyle.Render("r"), helpDescStyle.Render("Refresh agent list"),
		helpKeyStyle.Render("1-4 / a,s,r,t"), helpDescStyle.Render("Sort by Agent/Status/Runtime/Tokens"))

	return helpBoxStyle.Render(content)
}

// Width returns overlay width
func (h *HelpOverlay) Width() int {
	return h.width
}

// Height returns overlay height
func (h *HelpOverlay) Height() int {
	return h.height
}
