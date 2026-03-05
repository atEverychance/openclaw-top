package ui

import (
	"strings"
	"testing"
)

// TestHelpOverlayNoProcessMetrics verifies help text doesn't contain process metrics
func TestHelpOverlayNoProcessMetrics(t *testing.T) {
	help := NewHelpOverlay()
	view := help.View()
	
	// Convert to uppercase for case-insensitive check
	viewUpper := strings.ToUpper(view)
	
	forbidden := []string{"PID", "CPU", "MEMORY", "PROCESS"}
	for _, word := range forbidden {
		if strings.Contains(viewUpper, word) {
			t.Errorf("help overlay should not contain '%s', but found it in output", word)
		}
	}
}

// TestHelpOverlayContainsRequiredText verifies required phrases are present
func TestHelpOverlayContainsRequiredText(t *testing.T) {
	help := NewHelpOverlay()
	view := help.View()
	
	required := []string{
		"Select next/previous agent",
		"Sort by column",
		"Attach to selected agent",
		"View logs snapshot",
		"Kill selected agent",
	}
	
	for _, phrase := range required {
		if !strings.Contains(view, phrase) {
			t.Errorf("help overlay must contain '%s', but it was not found", phrase)
		}
	}
}

// TestHelpOverlayDimensions verifies dimensions work
func TestHelpOverlayDimensions(t *testing.T) {
	help := NewHelpOverlay()
	help.SetDimensions(80, 24)
	
	if help.Width() != 80 {
		t.Errorf("expected width 80, got %d", help.Width())
	}
	if help.Height() != 24 {
		t.Errorf("expected height 24, got %d", help.Height())
	}
}

// TestHelpOverlayViewNotEmpty verifies view renders content
func TestHelpOverlayViewNotEmpty(t *testing.T) {
	help := NewHelpOverlay()
	view := help.View()
	
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}
