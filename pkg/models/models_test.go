package models

import (
	"testing"
	"time"
)

func TestNewAppModel(t *testing.T) {
	app := NewAppModel()
	
	if app.Width != 80 {
		t.Errorf("expected default width 80, got %d", app.Width)
	}
	if app.Height != 24 {
		t.Errorf("expected default height 24, got %d", app.Height)
	}
	if app.View != ViewStateTable {
		t.Errorf("expected default view ViewStateTable, got %v", app.View)
	}
	if app.SortColumn != 0 {
		t.Errorf("expected default sort column 0, got %d", app.SortColumn)
	}
	if app.SortDesc {
		t.Error("expected default sort ascending")
	}
}

func TestAgentSessionStruct(t *testing.T) {
	s := AgentSession{
		AgentID:     "scout",
		Status:      "RUNNING",
		Runtime:     "14m 32s",
		TotalTokens: 4200,
		Model:       "kimi-k2.5",
		Task:        "Twitter hunt: AI agents",
	}
	
	if s.AgentID != "scout" {
		t.Errorf("expected AgentID scout, got %s", s.AgentID)
	}
	if s.Status != "RUNNING" {
		t.Errorf("expected Status RUNNING, got %s", s.Status)
	}
	if s.TotalTokens != 4200 {
		t.Errorf("expected TotalTokens 4200, got %d", s.TotalTokens)
	}
}

func TestAppStats(t *testing.T) {
	stats := AppStats{
		TotalAgents: 5,
		LastRefresh: time.Now(),
	}
	
	if stats.TotalAgents != 5 {
		t.Errorf("expected TotalAgents 5, got %d", stats.TotalAgents)
	}
}

func TestViewState(t *testing.T) {
	if ViewStateTable != 0 {
		t.Errorf("expected ViewStateTable = 0")
	}
	if ViewStateHelp != 1 {
		t.Errorf("expected ViewStateHelp = 1")
	}
}
