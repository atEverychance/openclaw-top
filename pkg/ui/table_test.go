package ui

import (
	"github.com/ateverychance/openclaw-top/pkg/models"
	"testing"
)

func TestTableDimensions(t *testing.T) {
	table := NewTable()
	table.SetDimensions(100, 30)
	
	if table.width != 100 {
		t.Errorf("expected width 100, got %d", table.width)
	}
	if table.height != 30 {
		t.Errorf("expected height 30, got %d", table.height)
	}
}

func TestTableSetData(t *testing.T) {
	table := NewTable()
	sessions := []models.AgentSession{
		{AgentID: "scout", Status: "RUNNING", Runtime: "14m 32s", TotalTokens: 4200, Task: "Test task"},
	}
	table.SetData(sessions)
	
	if len(table.sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(table.sessions))
	}
}

func TestTableSort(t *testing.T) {
	table := NewTable()
	sessions := []models.AgentSession{
		{AgentID: "coder", Status: "RUNNING", Runtime: "2m 10s", TotalTokens: 1800, Task: "Task C"},
		{AgentID: "scout", Status: "RUNNING", Runtime: "14m 32s", TotalTokens: 4200, Task: "Task A"},
		{AgentID: "bigBrain", Status: "RUNNING", Runtime: "8m 15s", TotalTokens: 12400, Task: "Task B"},
	}
	table.SetData(sessions)
	table.SetSort(0, false) // Sort by AgentID ascending
	
	// View triggers the sort internally
	_ = table.View()
	
	if table.sortColumn != 0 {
		t.Errorf("expected sort column 0, got %d", table.sortColumn)
	}
}

func TestTableSelected(t *testing.T) {
	table := NewTable()
	table.SetSelected(5)
	
	if table.selected != 5 {
		t.Errorf("expected selected 5, got %d", table.selected)
	}
}

func TestTableView(t *testing.T) {
	table := NewTable()
	table.SetDimensions(80, 20)
	table.SetData([]models.AgentSession{
		{AgentID: "scout", Status: "RUNNING", Runtime: "14m 32s", TotalTokens: 4200, Model: "kimi-k2.5", Task: "Twitter hunt: AI agents"},
	})
	
	view := table.View()
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}
