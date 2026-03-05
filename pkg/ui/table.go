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
				Bold(true).
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

	// Status icons
	runningIcon = "🟢"
	errorIcon   = "🔴"
	idleIcon    = "🟡"
	doneIcon    = "🔵"

	// Sparkline characters (low to high)
	sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
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
	progressBar := formatProgressBar(s.Runtime, t.sessions)
	sparkline := formatSparkline(s.TotalTokens, t.sessions)
	statusIcon := getStatusIcon(s.Status)

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

	// Format runtime with progress bar
	runtimeDisplay := fmt.Sprintf("%s %s", progressBar, s.Runtime)

	// Format tokens with sparkline
	tokensDisplay := fmt.Sprintf("%s %s", sparkline, tokensStr)

	return fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │",
		style.Width(agentW-1).Render(statusIcon+" "+truncate(s.AgentID, agentW-3)),
		statusStyle.Width(statusW-1).Render(s.Status),
		style.Width(runtimeW-1).Render(runtimeDisplay),
		style.Width(tokensW-1).Render(tokensDisplay),
		style.Width(taskW-1).Render(truncate(s.Task, taskW-1)))
}

// getStatusIcon returns the appropriate icon for a status
func getStatusIcon(status string) string {
	switch status {
	case "RUNNING":
		return runningIcon
	case "ERROR":
		return errorIcon
	case "IDLE":
		return idleIcon
	case "DONE":
		return doneIcon
	default:
		return "⚪"
	}
}

// formatProgressBar creates a mini progress bar for runtime
func formatProgressBar(runtime string, allSessions []models.AgentSession) string {
	// Parse runtime to get relative value
	currentMinutes := parseRuntimeMinutes(runtime)
	maxMinutes := getMaxRuntime(allSessions)

	if maxMinutes == 0 {
		return "░░░░░░░░░░"
	}

	// Calculate progress (0-10 blocks)
	progress := (currentMinutes * 10) / maxMinutes
	if progress > 10 {
		progress = 10
	}

	filled := strings.Repeat("█", progress)
	empty := strings.Repeat("░", 10-progress)
	return filled + empty
}

// parseRuntimeMinutes converts runtime string to minutes
func parseRuntimeMinutes(runtime string) int {
	var minutes, seconds int
	fmt.Sscanf(runtime, "%dm %ds", &minutes, &seconds)
	return minutes
}

// getMaxRuntime finds the maximum runtime in minutes
func getMaxRuntime(sessions []models.AgentSession) int {
	max := 0
	for _, s := range sessions {
		mins := parseRuntimeMinutes(s.Runtime)
		if mins > max {
			max = mins
		}
	}
	return max
}

// formatSparkline creates a mini sparkline for token usage
func formatSparkline(tokens int, allSessions []models.AgentSession) string {
	if len(allSessions) == 0 {
		return string(sparkChars[0])
	}

	// Find min and max tokens
	minTokens, maxTokens := allSessions[0].TotalTokens, allSessions[0].TotalTokens
	for _, s := range allSessions {
		if s.TotalTokens < minTokens {
			minTokens = s.TotalTokens
		}
		if s.TotalTokens > maxTokens {
			maxTokens = s.TotalTokens
		}
	}

	if maxTokens == minTokens {
		return string(sparkChars[len(sparkChars)/2])
	}

	// Normalize tokens to sparkline range
	rangeTokens := maxTokens - minTokens
	if rangeTokens == 0 {
		rangeTokens = 1
	}

	normalized := (tokens - minTokens) * (len(sparkChars) - 1) / rangeTokens
	if normalized < 0 {
		normalized = 0
	}
	if normalized >= len(sparkChars) {
		normalized = len(sparkChars) - 1
	}

	return string(sparkChars[normalized])
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
