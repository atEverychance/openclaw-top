package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/ateverychance/openclaw-top/pkg/models"
	"time"
)

// Client for OpenClaw Gateway
type Client struct {
	Endpoint string
	Timeout  time.Duration
}

// NewClient creates a new gateway client
func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
		Timeout:  5 * time.Second,
	}
}

// FetchStats fetches app stats from gateway (mock)
func (c *Client) FetchStats() (*models.AppStats, error) {
	// Mock implementation - returns agent statistics
	return &models.AppStats{
		TotalAgents: 3,
		LastRefresh: time.Now(),
	}, nil
}

// FetchSessions fetches agent sessions from gateway (mock)
func (c *Client) FetchSessions() ([]models.AgentSession, error) {
	// Mock implementation - returns realistic OpenClaw agent data
	// This matches the format from `openclaw sessions list --json`
	return []models.AgentSession{
		{
			AgentID:     "scout",
			Status:      "RUNNING",
			Runtime:     "14m 32s",
			TotalTokens: 4200,
			Model:       "kimi-k2.5",
			Task:        "Twitter hunt: AI agents",
		},
		{
			AgentID:     "coder",
			Status:      "ERROR",
			Runtime:     "2m 10s",
			TotalTokens: 1800,
			Model:       "kimi-k2.5",
			Task:        "Fix: auth token refresh",
		},
		{
			AgentID:     "bigBrain",
			Status:      "RUNNING",
			Runtime:     "8m 15s",
			TotalTokens: 12400,
			Model:       "kimi-k2.5",
			Task:        "Arch review: observatory",
		},
	}, nil
}

// FetchAll fetches both stats and sessions
func (c *Client) FetchAll() (*models.AppStats, []models.AgentSession, error) {
	stats, err := c.FetchStats()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	sessions, err := c.FetchSessions()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}
	return stats, sessions, nil
}

func init() {
	// Suppress unused import warning
	_ = json.Marshal
}
