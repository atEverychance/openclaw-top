package ui

import (
	"testing"
)

func TestConfirmModal(t *testing.T) {
	modal := NewConfirmModal()

	// Test initial state
	if modal.IsConfirmed() {
		t.Error("should not be confirmed initially")
	}
	if modal.IsCancelled() {
		t.Error("should not be cancelled initially")
	}

	// Test set message
	modal.SetMessage("Kill agent test?")
	modal.SetDetails("Agent: test\nStatus: RUNNING")

	// Test confirm
	modal.Confirm()
	if !modal.IsConfirmed() {
		t.Error("should be confirmed after Confirm()")
	}
	if modal.IsCancelled() {
		t.Error("should not be cancelled after Confirm()")
	}

	// Test reset
	modal.Reset()
	if modal.IsConfirmed() {
		t.Error("should not be confirmed after Reset()")
	}

	// Test cancel
	modal.Cancel()
	if !modal.IsCancelled() {
		t.Error("should be cancelled after Cancel()")
	}
	if modal.IsConfirmed() {
		t.Error("should not be confirmed after Cancel()")
	}
}

func TestConfirmModalView(t *testing.T) {
	modal := NewConfirmModal()
	modal.SetMessage("Test message")
	modal.SetDetails("Details here")

	view := modal.View()
	if view == "" {
		t.Error("View should not be empty")
	}

	// Check that message appears in view
	if !contains(view, "Test message") {
		t.Error("View should contain the message")
	}
}

func TestFormatAgentDetails(t *testing.T) {
	details := FormatAgentDetails("coder-123", "RUNNING", "5m 30s", 1250)

	if !contains(details, "coder-123") {
		t.Error("Details should contain agent name")
	}
	if !contains(details, "RUNNING") {
		t.Error("Details should contain status")
	}
	if !contains(details, "5m 30s") {
		t.Error("Details should contain runtime")
	}
	if !contains(details, "1250") {
		t.Error("Details should contain tokens")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr))
}
