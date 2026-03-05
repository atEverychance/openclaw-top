package ui

import (
	"fmt"
	"github.com/ateverychance/openclaw-top/pkg/models"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Table renders the agent sessions table
type Table struct {
	width       int
	height      int
	sessions    []models.AgentSession
	sortColumn  int
	sortDesc    bool
	selected    int
}

// Column widths (percentages)
const (
	colAgent   = 15
	colStatus  = 10
	colRuntime = 12
	colTokens  = 10
	colTask    = 53
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("236"))

	rowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	selectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("235"))

	alternateRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("250"))

	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")) // Green

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")) // Red

	idleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")) // Gray

	doneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")) // Blue
)

// NewTable creates a new table component
func NewTable() *Table {
	return &Table{
		width:       80,
		height:      20,
		sortColumn:  0,
		sortDesc:    false,
		selected:    0,
	}
}

// SetDimensions sets the table dimensions
func (t *Table) SetDimensions(width, height int) {
	t.width = width
	t.height = height
}

// SetData sets the agent session data
func (t *Table) SetData(sessions []models.AgentSession) {
	t.sessions = sessions
}

// SetSort sets the sort column and direction
func (t *Table) SetSort(column int, desc bool) {
	t.sortColumn = column
	t.sortDesc = desc
}

// SetSelected sets the selected row
func (t *Table) SetSelected(selected int) {
	t.selected = selected
}

// View renders the table
func (t *Table) View() string {
	var sb strings.Builder

	// Calculate column widths based on available width
	usableWidth := t.width - 2 // Account for borders
	agentWidth := (usableWidth * colAgent) / 100
	statusWidth := (usableWidth * colStatus) / 100
	runtimeWidth := (usableWidth * colRuntime) / 100
	tokensWidth := (usableWidth * colTokens) / 100
	taskWidth := usableWidth - agentWidth - statusWidth - runtimeWidth - tokensWidth - 6

	// Header
	header := fmt.Sprintf("│ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │",
		agentWidth-1, "AGENT",
		statusWidth-1, "STATUS",
		runtimeWidth-1, "RUNTIME",
		tokensWidth-1, "TOKENS",
		taskWidth-1, "TASK")
	sb.WriteString(headerStyle.Render(header) + "\n")

	// Sort sessions
	sorted := make([]models.AgentSession, len(t.sessions))
	copy(sorted, t.sessions)
	sort.Slice(sorted, func(i, j int) bool {
		var less bool
		switch t.sortColumn {
		case 0:
			less = sorted[i].AgentID < sorted[j].AgentID
		case 1:
			less = sorted[i].Status < sorted[j].Status
		case 2:
			less = sorted[i].Runtime < sorted[j].Runtime
		case 3:
			less = sorted[i].TotalTokens < sorted[j].TotalTokens
		case 4:
			less = sorted[i].Task < sorted[j].Task
		default:
			less = sorted[i].AgentID < sorted[j].AgentID
		}
		if t.sortDesc {
			return !less
		}
		return less
	})

	// Rows
	maxRows := t.height - 2
	if len(sorted) < maxRows {
		maxRows = len(sorted)
	}

	for i := 0; i < maxRows; i++ {
		s := sorted[i]
		rowStr := t.formatRow(s, i, agentWidth, statusWidth, runtimeWidth, tokensWidth, taskWidth)
		sb.WriteString(rowStr + "\n")
	}

	return sb.String()
}

func (t *Table) formatRow(s models.AgentSession, index, agentW, statusW, runtimeW, tokensW, taskW int) string {
	tokensStr := formatTokens(s.TotalTokens)

	// Apply color based on status
	statusStyle := rowStyle
	switch s.Status {
	case "RUNNING":
		statusStyle = runningStyle
	case "ERROR":
		statusStyle = errorStyle
	case "IDLE":
		statusStyle = idleStyle
	case "DONE":
		statusStyle = doneStyle
	}

	// Select style based on row state
	style := rowStyle
	if index == t.selected {
		style = selectedRowStyle
	} else if index%2 == 1 {
		style = alternateRowStyle
	}

	return fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
		style.Width(agentW-1).Render(truncate(s.AgentID, agentW-1)),
		statusStyle.Width(statusW-1).Render(s.Status),
		style.Width(runtimeW-1).Render(s.Runtime),
		style.Width(tokensW-1).Render(tokensStr),
		style.Width(taskW-1).Render(truncate(s.Task, taskW-1)))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + ".."
}

func formatTokens(tokens int) string {
	if tokens >= 1000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}
