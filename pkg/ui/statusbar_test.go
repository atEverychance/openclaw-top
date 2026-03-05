package ui

import (
	"github.com/ateverychance/openclaw-top/pkg/models"
	"testing"
	"time"
)

func TestStatusBarDimensions(t *testing.T) {
	sb := NewStatusBar()
	sb.SetDimensions(100)
	
	if sb.width != 100 {
		t.Errorf("expected width 100, got %d", sb.width)
	}
}

func TestStatusBarSetStats(t *testing.T) {
	sb := NewStatusBar()
	stats := &models.AppStats{
		TotalAgents: 3,
		LastRefresh: time.Now(),
	}
	sb.SetStats(stats)
	
	if sb.stats.TotalAgents != 3 {
		t.Errorf("expected TotalAgents 3, got %d", sb.stats.TotalAgents)
	}
}

func TestStatusBarView(t *testing.T) {
	sb := NewStatusBar()
	sb.SetDimensions(80)
	sb.SetStats(&models.AppStats{
		TotalAgents: 3,
		LastRefresh: time.Now(),
	})
	
	view := sb.View()
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}

func TestStatusBarNilStats(t *testing.T) {
	sb := NewStatusBar()
	sb.SetDimensions(80)
	// Don't set stats - should handle nil gracefully
	
	view := sb.View()
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}
