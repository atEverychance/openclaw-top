package main

import (
	"flag"
	"fmt"
	"os"

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
}

type model struct {
	table      *ui.Table
	statusBar  *ui.StatusBar
	help       *ui.HelpOverlay
	client     Client
	app        *models.AppModel
}

func initialModel() *model {
	client := gateway.NewOpenClawClient()
	app := models.NewAppModel()

	return &model{
		table:     ui.NewTable(),
		statusBar: ui.NewStatusBar(),
		help:      ui.NewHelpOverlay(),
		client:    client,
		app:       app,
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

type DataMsg struct {
	Stats    *models.AppStats
	Sessions []models.AgentSession
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

	case TickMsg:
		// Auto-refresh triggered
		return m, m.fetchData

	case tea.KeyMsg:
		return m.handleKey(msg)

	default:
		return m, nil
	}
}

type TickMsg struct{}

func (m *model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		if m.app.View == models.ViewStateHelp {
			m.app.View = models.ViewStateTable
		} else {
			m.app.View = models.ViewStateHelp
		}

	case "r":
		return m, m.fetchData

	case "up", "k":
		if m.app.Selected > 0 {
			m.app.Selected--
			m.table.SetSelected(m.app.Selected)
		}

	case "down", "j":
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

func (m *model) updateLayout() {
	tableHeight := m.app.Height - statusBarHeight - 2 // Account for borders
	if tableHeight < 1 {
		tableHeight = 1
	}
	m.table.SetDimensions(m.app.Width, tableHeight)
	m.statusBar.SetDimensions(m.app.Width)
	m.help.SetDimensions(m.app.Width, m.app.Height)
}

func (m *model) View() string {
	if m.app.Width < minWidth || m.app.Height < minHeight {
		return fmt.Sprintf("Terminal too small (min %dx%d)", minWidth, minHeight)
	}

	if m.app.View == models.ViewStateHelp {
		helpView := m.help.View()
		helpHeight := lipgloss.Height(helpView)
		helpWidth := lipgloss.Width(helpView)
		topPadding := (m.app.Height - helpHeight) / 2
		leftPadding := (m.app.Width - helpWidth) / 2
		return lipgloss.PlaceVertical(topPadding, lipgloss.Top,
			lipgloss.PlaceHorizontal(leftPadding, lipgloss.Left, helpView))
	}

	// Render main view: table + status bar
	tableView := m.table.View()
	statusView := m.statusBar.View()

	return tableView + "\n" + borderStyle.Width(m.app.Width).Render(statusView)
}

func main() {
	// Parse flags
	showVersion := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("openclaw-top %s\n", version.FullInfo())
		os.Exit(0)
	}

	run()
}

func run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
