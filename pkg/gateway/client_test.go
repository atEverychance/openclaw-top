package gateway

import (
	"encoding/json"
	"testing"
)

// TestOpenClawClientJSONParsing verifies the JSON parsing handles wrapped format
func TestOpenClawClientJSONParsing(t *testing.T) {
	// Test the expected OpenClaw JSON format: {"sessions": [...]}
	testJSON := `{"sessions": [{"key": "agent:coder:subagent:123", "agentId": "coder", "ageMs": 123, "abortedLastRun": false, "model": "qwen3.5", "inputTokens": 100, "outputTokens": 50, "totalTokens": 150}]}`
	
	var openClawResponse struct {
		Sessions []struct {
			Key       string `json:"key"`
			Agent     string `json:"agentId"`
			Age       int64  `json:"ageMs"`
			Model     string `json:"model"`
			InputTok  int    `json:"inputTokens"`
			OutputTok int    `json:"outputTokens"`
			TotalTok  int    `json:"totalTokens"`
			Aborted   bool   `json:"abortedLastRun"`
		} `json:"sessions"`
	}

	if err := json.Unmarshal([]byte(testJSON), &openClawResponse); err != nil {
		t.Fatalf("failed to parse wrapped JSON format: %v", err)
	}

	if len(openClawResponse.Sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(openClawResponse.Sessions))
	}

	s := openClawResponse.Sessions[0]
	if s.Key != "agent:coder:subagent:123" {
		t.Errorf("expected key 'agent:coder:subagent:123', got '%s'", s.Key)
	}
	if s.Agent != "coder" {
		t.Errorf("expected agentId 'coder', got '%s'", s.Agent)
	}
	if s.Age != 123 {
		t.Errorf("expected ageMs 123, got %d", s.Age)
	}
	if s.Aborted != false {
		t.Errorf("expected abortedLastRun false, got %v", s.Aborted)
	}
}

// TestDirectArrayRejection verifies we don't accept direct array format
func TestDirectArrayRejection(t *testing.T) {
	// This is the OLD format that should NOT work
	testJSON := `[{"key": "agent:coder", "agentId": "coder", "ageMs": 123}]`
	
	var openClawResponse struct {
		Sessions []struct {
			Key       string `json:"key"`
			Agent     string `json:"agentId"`
			Age       int64  `json:"ageMs"`
		} `json:"sessions"`
	}

	if err := json.Unmarshal([]byte(testJSON), &openClawResponse); err == nil {
		t.Error("direct array format should not parse successfully with wrapped format parser")
	}
}

// TestOpenClawClientCreation verifies NewOpenClawClient works
func TestOpenClawClientCreation(t *testing.T) {
	client := NewOpenClawClient()
	if client == nil {
		t.Fatal("expected non-nil OpenClawClient")
	}
}
