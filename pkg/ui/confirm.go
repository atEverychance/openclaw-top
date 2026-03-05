package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ConfirmModal represents a confirmation dialog
type ConfirmModal struct {
	width     int
	height    int
	message   string
	details   string
	confirmed bool
	cancelled bool
}

// NewConfirmModal creates a new confirmation modal
func NewConfirmModal() *ConfirmModal {
	return &ConfirmModal{
		width:   60,
		height:  10,
		message: "Confirm action?",
	}
}

// SetMessage sets the confirmation message
func (c *ConfirmModal) SetMessage(msg string) {
	c.message = msg
}

// SetDetails sets additional details to display
func (c *ConfirmModal) SetDetails(details string) {
	c.details = details
}

// SetDimensions sets the modal dimensions
func (c *ConfirmModal) SetDimensions(width, height int) {
	c.width = width
	c.height = height
}

// Confirm marks the dialog as confirmed
func (c *ConfirmModal) Confirm() {
	c.confirmed = true
	c.cancelled = false
}

// Cancel marks the dialog as cancelled
func (c *ConfirmModal) Cancel() {
	c.confirmed = false
	c.cancelled = true
}

// IsConfirmed returns true if user confirmed
func (c *ConfirmModal) IsConfirmed() bool {
	return c.confirmed
}

// IsCancelled returns true if user cancelled
func (c *ConfirmModal) IsCancelled() bool {
	return c.cancelled
}

// Reset clears the confirmation state
func (c *ConfirmModal) Reset() {
	c.confirmed = false
	c.cancelled = false
}

// View renders the confirmation modal
func (c *ConfirmModal) View() string {
	var sb strings.Builder

	// Styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2).
		Width(c.width - 4)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	detailStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))

	// Build content
	sb.WriteString(titleStyle.Render("⚠️  " + c.message))
	sb.WriteString("\n\n")

	if c.details != "" {
		sb.WriteString(detailStyle.Render(c.details))
		sb.WriteString("\n\n")
	}

	sb.WriteString(helpStyle.Render("[y] Confirm    [n] Cancel    [q] Quit"))

	return boxStyle.Render(sb.String())
}

// GetHeight returns the modal height
func (c *ConfirmModal) GetHeight() int {
	return lipgloss.Height(c.View())
}

// GetWidth returns the modal width
func (c *ConfirmModal) GetWidth() int {
	return lipgloss.Width(c.View())
}

// FormatAgentDetails formats agent details for display
func FormatAgentDetails(name, status, runtime string, tokens int) string {
	return fmt.Sprintf("Agent: %s\nStatus: %s\nRuntime: %s\nTokens: %d",
		name, status, runtime, tokens)
}
