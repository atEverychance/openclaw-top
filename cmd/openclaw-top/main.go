package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ateverychance/openclaw-top/internal/version"
	"github.com/ateverychance/openclaw-top/pkg/gateway"
	"github.com/ateverychance/openclaw-top/pkg/models"
	"github.com/ateverychance/openclaw-top/pkg/ui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	statusBarHeight = 1
	minWidth        = 60
	minHeight       = 10
)

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("235"))
)

type Client interface {
	FetchAll() (*models.AppStats, []models.AgentSession, error)
	KillSession(sessionID string) error
	GetLogs(sessionID string, lines int) (string, error)
}

type model struct {
	table       *ui.Table
	statusBar   *ui.StatusBar
	help        *ui.HelpOverlay
	confirm     *ui.ConfirmModal
	logViewer   *ui.LogViewer
	attachView  *ui.AttachView
	client      Client
	app         *models.AppModel
	logContent  string
}

func initialModel() *model {
	client := gateway.NewOpenClawClient()
	app := models.NewAppModel()

	return &model{
		table:      ui.NewTable(),
		statusBar:  ui.NewStatusBar(),
		help:       ui.NewHelpOverlay(),
		confirm:    ui.NewConfirmModal(),
		logViewer:  ui.NewLogViewer(),
		attachView: ui.NewAttachView(),
		client:     client,
		app:        app,
	}
}

func (m *model) Init() tea.Cmd {
	return m.fetchData
}

func (m *model) fetchData() tea.Msg {
	stats, sessions, err := m.client.FetchAll()
	if err != nil {
		return fmt.Errorf("fetch error: %v", err)
	}
	return DataMsg{Stats: stats, Sessions: sessions}
}

func (m *model) killAgent() tea.Msg {
	if m.app.Selected < 0 || m.app.Selected >= len(m.app.Sessions) {
		return fmt.Errorf("no agent selected")
	}
	
	session := m.app.Sessions[m.app.Selected]
	err := m.client.KillSession(session.AgentID)
	if err != nil {
		return fmt.Errorf("kill failed: %v", err)
	}
	return KillSuccessMsg{AgentID: session.AgentID}
}

func (m *model) fetchLogs() tea.Msg {
	if m.app.Selected < 0 || m.app.Selected >= len(m.app.Sessions) {
		return fmt.Errorf("no agent selected")
	}
	
	session := m.app.Sessions[m.app.Selected]
	logs, err := m.client.GetLogs(session.AgentID, 100)
	if err != nil {
		return fmt.Errorf("fetch logs failed: %v", err)
	}
	return LogsMsg{Content: logs, AgentID: session.AgentID}
}

type DataMsg struct {
	Stats    *models.AppStats
	Sessions []models.AgentSession
}

type KillSuccessMsg struct {
	AgentID string
}

type LogsMsg struct {
	Content string
	AgentID string
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.app.Width = msg.Width
		m.app.Height = msg.Height
		m.updateLayout()
		return m, nil

	case DataMsg:
		m.app.Stats = msg.Stats
		m.app.Sessions = msg.Sessions
		m.table.SetData(msg.Sessions)
		m.statusBar.SetStats(msg.Stats)
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return TickMsg{}
		})

	case KillSuccessMsg:
		m.confirm.Reset()
		m.app.View = models.ViewStateTable
		m.statusBar.SetMessage(fmt.Sprintf("✓ Killed agent %s", msg.AgentID))
		return m, m.fetchData

	case LogsMsg:
		m.logContent = msg.Content
		m.logViewer.SetContent(msg.Content)
		m.logViewer.SetTitle(fmt.Sprintf("Logs: %s", msg.AgentID))
		m.app.View = models.ViewStateLogs
		return m, nil

	case ui.LogUpdateMsg:
		// Handle live log updates in attach mode
		_, cmd := m.attachView.Update(msg)
		return m, cmd

	case error:
		m.statusBar.SetError(msg)
		return m, nil

	case TickMsg:
		// Auto-refresh triggered
		return m, m.fetchData

	case tea.KeyMsg:
		return m.handleKey(msg)

	default:
		// Handle log viewer updates
		if m.app.View == models.ViewStateLogs {
			_, cmd := m.logViewer.Update(msg)
			return m, cmd
		}
		// Handle attach view updates
		if m.app.View == models.ViewStateAttach {
			_, cmd := m.attachView.Update(msg)
			return m, cmd
		}
		return m, nil
	}
}

type TickMsg struct{}

func (m *model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle view-specific keys first
	switch m.app.View {
	case models.ViewStateConfirm:
		return m.handleConfirmKeys(key)
	case models.ViewStateLogs:
		return m.handleLogKeys(key)
	case models.ViewStateAttach:
		return m.handleAttachKeys(key)
	case models.ViewStateHelp:
		// Any key exits help view
		m.app.View = models.ViewStateTable
		return m, nil
	}

	// Table view keys
	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		m.app.View = models.ViewStateHelp

	case "r":
		return m, m.fetchData

	case "k":
		if m.app.Selected >= 0 && m.app.Selected < len(m.app.Sessions) {
			session := m.app.Sessions[m.app.Selected]
			if session.Status == "RUNNING" {
				m.confirm.SetMessage(fmt.Sprintf("Kill agent %s?", session.AgentID))
				m.confirm.SetDetails(ui.FormatAgentDetails(
					session.AgentID,
					session.Status,
					session.Runtime,
					session.TotalTokens,
				))
				m.app.View = models.ViewStateConfirm
			} else {
				m.statusBar.SetMessage("Cannot kill: agent is not running")
			}
		}

	case "l":
		if m.app.Selected >= 0 && m.app.Selected < len(m.app.Sessions) {
			return m, m.fetchLogs
		}

	case "a":
		if m.app.Selected >= 0 && m.app.Selected < len(m.app.Sessions) {
			session := m.app.Sessions[m.app.Selected]
			// Only allow attach on RUNNING sessions (not IDLE or ERROR)
			if session.Status != "RUNNING" {
				m.statusBar.SetMessage(fmt.Sprintf("Cannot attach to %s agent (not running)", session.Status))
				return m, nil
			}
			m.attachView.SetSession(session.AgentID, session.AgentID)
			m.app.View = models.ViewStateAttach
			return m, m.attachView.Init()
		}

	case "up":
		if m.app.Selected > 0 {
			m.app.Selected--
			m.table.SetSelected(m.app.Selected)
		}

	case "down":
		if m.app.Selected < len(m.app.Sessions)-1 {
			m.app.Selected++
			m.table.SetSelected(m.app.Selected)
		}

	case "1":
		m.app.SortColumn = 0
		m.app.SortDesc = !m.app.SortDesc
		m.table.SetSort(m.app.SortColumn, m.app.SortDesc)

	case "2":
		m.app.SortColumn = 1
		m.app.SortDesc = !m.app.SortDesc
		m.table.SetSort(m.app.SortColumn, m.app.SortDesc)

	case "3":
		m.app.SortColumn = 2
		m.app.SortDesc = !m.app.SortDesc
		m.table.SetSort(m.app.SortColumn, m.app.SortDesc)

	case "4":
		m.app.SortColumn = 3
		m.app.SortDesc = !m.app.SortDesc
		m.table.SetSort(m.app.SortColumn, m.app.SortDesc)
	}

	return m, nil
}

func (m *model) handleConfirmKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "y", "Y":
		m.confirm.Confirm()
		return m, m.killAgent
	case "n", "N", "q", "esc":
		m.confirm.Cancel()
		m.confirm.Reset()
		m.app.View = models.ViewStateTable
		return m, nil
	}
	return m, nil
}

func (m *model) handleLogKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc":
		m.app.View = models.ViewStateTable
		return m, nil
	case "up", "k":
		m.logViewer.ScrollUp()
	case "down", "j":
		m.logViewer.ScrollDown()
	case "pgup":
		m.logViewer.PageUp()
	case "pgdown":
		m.logViewer.PageDown()
	case "end", "G":
		m.logViewer.GotoBottom()
	}
	return m, nil
}

func (m *model) handleAttachKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "esc":
		m.attachView.Stop()
		m.app.View = models.ViewStateTable
		return m, nil
	case "up":
		m.attachView.ScrollUp()
	case "down":
		m.attachView.ScrollDown()
	case "pgup":
		m.attachView.PageUp()
	case "pgdown":
		m.attachView.PageDown()
	case "end", "G":
		m.attachView.GotoBottom()
	}
	return m, nil
}

func (m *model) updateLayout() {
	tableHeight := m.app.Height - statusBarHeight - 2 // Account for borders
	if tableHeight < 1 {
		tableHeight = 1
	}
	m.table.SetDimensions(m.app.Width, tableHeight)
	m.statusBar.SetDimensions(m.app.Width)
	m.help.SetDimensions(m.app.Width, m.app.Height)
	m.confirm.SetDimensions(m.app.Width, m.app.Height)
	m.logViewer.SetDimensions(m.app.Width, m.app.Height)
	m.attachView.SetDimensions(m.app.Width, m.app.Height)
}

func (m *model) View() string {
	if m.app.Width < minWidth || m.app.Height < minHeight {
		return fmt.Sprintf("Terminal too small (min %dx%d)", minWidth, minHeight)
	}

	switch m.app.View {
	case models.ViewStateHelp:
		return m.renderHelpView()
	case models.ViewStateConfirm:
		return m.renderConfirmView()
	case models.ViewStateLogs:
		return m.renderLogView()
	case models.ViewStateAttach:
		return m.renderAttachView()
	default:
		return m.renderTableView()
	}
}

func (m *model) renderHelpView() string {
	helpView := m.help.View()
	helpHeight := lipgloss.Height(helpView)
	helpWidth := lipgloss.Width(helpView)
	topPadding := (m.app.Height - helpHeight) / 2
	leftPadding := (m.app.Width - helpWidth) / 2
	return lipgloss.PlaceVertical(topPadding, lipgloss.Top,
		lipgloss.PlaceHorizontal(leftPadding, lipgloss.Left, helpView))
}

func (m *model) renderConfirmView() string {
	confirmView := m.confirm.View()
	confirmHeight := lipgloss.Height(confirmView)
	confirmWidth := lipgloss.Width(confirmView)
	topPadding := (m.app.Height - confirmHeight) / 2
	leftPadding := (m.app.Width - confirmWidth) / 2
	return lipgloss.PlaceVertical(topPadding, lipgloss.Top,
		lipgloss.PlaceHorizontal(leftPadding, lipgloss.Left, confirmView))
}

func (m *model) renderLogView() string {
	return m.logViewer.View()
}

func (m *model) renderAttachView() string {
	return m.attachView.View()
}

func (m *model) renderTableView() string {
	tableView := m.table.View()
	statusView := m.statusBar.View()
	return tableView + "\n" + borderStyle.Width(m.app.Width).Render(statusView)
}

func main() {
	// Parse flags
	showVersion := flag.Bool("version", false, "Print version information and exit")
	attachAgent := flag.String("attach", "", "Attach to specific agent immediately (non-interactive)")
	killAgent := flag.String("kill", "", "Kill specific agent and exit (non-interactive)")
	watchMode := flag.Bool("watch", false, "Watch mode: stream updates without interactive TUI")
	refreshRate := flag.Int("refresh", 2, "Refresh interval in seconds (default: 2)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("openclaw-top %s\n", version.FullInfo())
		os.Exit(0)
	}

	// Handle non-interactive modes
	if *attachAgent != "" {
		os.Exit(attachToAgent(*attachAgent, *refreshRate))
	}

	if *killAgent != "" {
		os.Exit(killAgentByID(*killAgent))
	}

	if *watchMode {
		os.Exit(watchModeRun(*refreshRate))
	}

	run()
}

// attachToAgent attaches to a specific agent and streams logs
func attachToAgent(agentID string, refreshRate int) int {
	client := gateway.NewOpenClawClient()

	// Verify agent exists
	_, sessions, err := client.FetchAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching sessions: %v\n", err)
		return 1
	}

	var targetSession *models.AgentSession
	for i := range sessions {
		if sessions[i].AgentID == agentID {
			targetSession = &sessions[i]
			break
		}
	}

	if targetSession == nil {
		fmt.Fprintf(os.Stderr, "Agent not found: %s\n", agentID)
		return 2
	}

	if targetSession.Status != "RUNNING" {
		fmt.Fprintf(os.Stderr, "Agent is not running: %s (status: %s)\n", agentID, targetSession.Status)
		return 2
	}

	fmt.Printf("Attaching to agent %s...\n", agentID)
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Simple polling loop for logs
	ticker := time.NewTicker(time.Duration(refreshRate) * time.Second)
	defer ticker.Stop()

	var lastContent string
	for range ticker.C {
		logs, err := client.GetLogs(agentID, 50)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
			return 3
		}

		// Only print new content
		if logs != lastContent {
			fmt.Print(logs)
			lastContent = logs
		}
	}

	return 0
}

// killAgentByID kills a specific agent by ID
func killAgentByID(agentID string) int {
	client := gateway.NewOpenClawClient()

	err := client.KillSession(agentID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to kill agent %s: %v\n", agentID, err)
		return 3
	}

	fmt.Printf("✓ Killed agent %s\n", agentID)
	return 0
}

// watchModeRun runs in watch mode (streaming updates)
func watchModeRun(refreshRate int) int {
	client := gateway.NewOpenClawClient()

	fmt.Println("openclaw-top -- Watch Mode")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	ticker := time.NewTicker(time.Duration(refreshRate) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats, sessions, err := client.FetchAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		// Clear screen (simple approach)
		fmt.Print("\033[2J\033[H")

		fmt.Printf("Agents: %d | Refresh: %s\n", stats.TotalAgents, time.Now().Format("15:04:05"))
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-20s %-10s %-12s %-10s\n", "AGENT", "STATUS", "RUNTIME", "TOKENS")
		fmt.Println(strings.Repeat("-", 80))

		for _, s := range sessions {
			fmt.Printf("%-20s %-10s %-12s %-10d\n",
				truncate(s.AgentID, 20),
				s.Status,
				s.Runtime,
				s.TotalTokens)
		}
	}

	return 0
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
