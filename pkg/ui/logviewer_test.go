package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestLogViewer(t *testing.T) {
	viewer := NewLogViewer()

	// Test initial state
	if viewer.ready {
		t.Error("should not be ready initially")
	}

	// Test set content
	content := "Line 1\nLine 2\nLine 3"
	viewer.SetContent(content)

	// Test set title
	viewer.SetTitle("Test Logs")

	// Test dimensions
	viewer.SetDimensions(100, 30)
}

func TestLogViewerView(t *testing.T) {
	viewer := NewLogViewer()
	viewer.SetContent("Test log content")
	viewer.SetTitle("Test")
	viewer.SetDimensions(80, 20)

	// Before init, should show loading
	view := viewer.View()
	if view == "" {
		t.Error("View should not be empty")
	}
}

func TestLogViewerScrolling(t *testing.T) {
	viewer := NewLogViewer()
	viewer.SetContent("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")
	viewer.SetDimensions(80, 10)

	// Test that scrolling methods don't panic before ready
	viewer.ScrollUp()
	viewer.ScrollDown()
	viewer.PageUp()
	viewer.PageDown()
	viewer.GotoBottom()
}

func TestLogViewerUpdate(t *testing.T) {
	viewer := NewLogViewer()
	viewer.SetDimensions(80, 20)

	// First update initializes the viewport
	updated, cmd := viewer.Update(tea.KeyMsg{})
	if updated == nil {
		t.Error("Update should return the viewer")
	}
	if cmd != nil {
		t.Error("First update should return nil cmd")
	}

	// After first update, should be ready
	if !updated.ready {
		t.Error("Viewer should be ready after first update")
	}
}
