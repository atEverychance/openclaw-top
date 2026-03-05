package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ateverychance/openclaw-top/pkg/models"
)

// OpenClawClient fetches real data from OpenClaw CLI
type OpenClawClient struct{}

// NewOpenClawClient creates a client that calls the real openclaw CLI
func NewOpenClawClient() *OpenClawClient {
	return &OpenClawClient{}
}

// FetchStats returns app stats based on actual session count
func (c *OpenClawClient) FetchStats() (*models.AppStats, error) {
	sessions, err := c.FetchSessions()
	if err != nil {
		return nil, err
	}
	return &models.AppStats{
		TotalAgents: len(sessions),
		LastRefresh: time.Now(),
	}, nil
}

// FetchSessions calls `openclaw sessions list --json` and parses the real output
func (c *OpenClawClient) FetchSessions() ([]models.AgentSession, error) {
	// Execute the real OpenClaw CLI with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "openclaw", "sessions", "--json")
	output, err := cmd.Output()
	if err != nil {
		// If command fails, return error so UI can show it
		return nil, fmt.Errorf("openclaw CLI failed: %w", err)
	}

	// Parse OpenClaw's JSON output (wrapped in {sessions: [...]})
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

	if err := json.Unmarshal(output, &openClawResponse); err != nil {
		return nil, fmt.Errorf("failed to parse openclaw output: %w", err)
	}

	// Convert to our AgentSession format
	sessions := make([]models.AgentSession, 0, len(openClawResponse.Sessions))
	for _, s := range openClawResponse.Sessions {
		status := deriveStatus(!s.Aborted, s.Age)
		sessions = append(sessions, models.AgentSession{
			AgentID:     deriveAgentID(s.Key, s.Agent),
			Status:      status,
			Runtime:     formatRuntime(s.Age),
			TotalTokens: s.TotalTok,
			Model:       s.Model,
			Task:        deriveTask(s.Key),
		})
	}

	return sessions, nil
}

// FetchAll fetches both stats and sessions
func (c *OpenClawClient) FetchAll() (*models.AppStats, []models.AgentSession, error) {
	sessions, err := c.FetchSessions()
	if err != nil {
		return nil, nil, err
	}
	stats := &models.AppStats{
		TotalAgents: len(sessions),
		LastRefresh: time.Now(),
	}
	return stats, sessions, nil
}

// KillSession kills a session by ID
func (c *OpenClawClient) KillSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First try to find the session key from the sessionID
	sessions, err := c.FetchSessions()
	if err != nil {
		return fmt.Errorf("failed to fetch sessions: %w", err)
	}

	var targetKey string
	for _, s := range sessions {
		if s.AgentID == sessionID {
			targetKey = s.AgentID
			break
		}
	}

	if targetKey == "" {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Try to kill via openclaw CLI
	cmd := exec.CommandContext(ctx, "openclaw", "sessions", "kill", targetKey)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to kill session: %w (output: %s)", err, string(output))
	}

	return nil
}

// GetLogs retrieves logs for a session
func (c *OpenClawClient) GetLogs(sessionID string, lines int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if lines <= 0 {
		lines = 100
	}

	// Try to get logs via openclaw CLI
	cmd := exec.CommandContext(ctx, "openclaw", "sessions", "logs", sessionID, "--lines", fmt.Sprintf("%d", lines))
	output, err := cmd.Output()
	if err != nil {
		// Fallback: try to read from log files directly
		return c.getLogsFromFile(sessionID, lines)
	}

	return string(output), nil
}

// getLogsFromFile attempts to read logs from file system
func (c *OpenClawClient) getLogsFromFile(sessionID string, lines int) (string, error) {
	// This is a fallback that looks for log files in common locations
	// In a real implementation, this would look in ~/.openclaw/sessions/ etc.
	return fmt.Sprintf("Logs not available via CLI for session %s\n(Fallback file reading not implemented)", sessionID), nil
}

// deriveStatus determines status from OpenClaw session data
func deriveStatus(active bool, ageMs int64) string {
	if !active {
		return "IDLE"
	}
	// If active and recent, it's running
	if ageMs < 300000 { // less than 5 minutes
		return "RUNNING"
	}
	return "IDLE"
}

// deriveAgentID extracts a readable agent name from session key
func deriveAgentID(key, agent string) string {
	// If agent field is set, use it
	if agent != "" {
		return agent
	}
	// Otherwise parse from key like "agent:coder:subagent:xxx"
	parts := strings.Split(key, ":")
	if len(parts) >= 2 {
		return parts[1] // e.g., "coder" from "agent:coder:..."
	}
	return key[:min(len(key), 20)] // truncate if needed
}

// deriveTask creates a task description from session key
func deriveTask(key string) string {
	if strings.Contains(key, ":scout:") {
		return "Twitter monitoring / signal hunting"
	}
	if strings.Contains(key, ":coder:") {
		return "Code implementation / bug fixes"
	}
	if strings.Contains(key, ":bigbrain:") || strings.Contains(key, ":senior-coder:") {
		return "Architecture review / complex debugging"
	}
	if strings.Contains(key, ":researcher:") {
		return "Web research / fact checking"
	}
	if strings.Contains(key, ":git-manager:") {
		return "GitHub issue tracking / PR management"
	}
	// Generic description based on key pattern
	parts := strings.Split(key, ":")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s task execution", parts[1])
	}
	return "Agent task"
}

// formatRuntime converts milliseconds to readable format
func formatRuntime(ageMs int64) string {
	age := time.Duration(ageMs) * time.Millisecond
	minutes := int(age.Minutes())
	seconds := int(age.Seconds()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
