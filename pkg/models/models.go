package models

import "time"

// AgentSession represents an OpenClaw agent session
type AgentSession struct {
	AgentID     string
	Status      string // RUNNING, IDLE, ERROR, DONE
	Runtime     string // e.g., "14m 32s"
	TotalTokens int
	Model       string
	Task        string
}

// AppStats represents application-level statistics
type AppStats struct {
	TotalAgents  int
	LastRefresh  time.Time
}

// ViewState represents the current view mode
type ViewState int

const (
	ViewStateTable ViewState = iota
	ViewStateHelp
	ViewStateConfirm
	ViewStateLogs
	ViewStateAttach
)

// AppModel is the main application model
type AppModel struct {
	Stats      *AppStats
	Sessions   []AgentSession
	Width      int
	Height     int
	View       ViewState
	SortColumn int
	SortDesc   bool
	Selected   int
}

// NewAppModel creates a new application model
func NewAppModel() *AppModel {
	return &AppModel{
		Stats:      &AppStats{},
		Sessions:   []AgentSession{},
		Width:      80,
		Height:     24,
		View:       ViewStateTable,
		SortColumn: 0,
		SortDesc:   false,
		Selected:   0,
	}
}
