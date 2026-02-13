package ui

import (
	"strings"
	
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
	"github.com/youssef28m/LockIn/internal/ui/pages"
)

// ================================================================
// Root Model
// ================================================================

type RootModel struct {
	page       common.Page
	home       pages.HomeModel
	setTimer   pages.SetTimerModel
	timer      pages.TimerModel
	blockSites pages.BlockSitesModel
	blockList  pages.BlockListModel
	help       help.Model
	service    *service.AppService
	width      int
	height     int
}

func NewRootModel(service *service.AppService) *RootModel {
	return &RootModel{
		page:       common.HomePage,
		home:       pages.NewHomeModel(service),
		setTimer:   pages.NewSetTimerModel(service),
		timer:      pages.NewTimerModel(service),
		blockSites: pages.NewBlockSitesModel(service),
		blockList:  pages.NewBlockListModel(service),
		help:       help.New(),
		service:    service,
	}
}

// ================================================================
// Bubble Tea Lifecycle
// ================================================================

func (m *RootModel) Init() tea.Cmd {
	session, err := m.service.GetActiveSession()

	if err == nil && session.Remaining() > 0 {
		m.page = common.TimerPage

		return func() tea.Msg {
			return common.StartTimerMsg{
				DurationSeconds: int(session.Remaining()),
			}
		}
	}

	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	switch msg := msg.(type) {

	case common.NavigateMsg:
		m.page = msg.Target

		if m.page == common.BlockListPage {
			return m, common.FetchBlockedSitesCmd(*m.service)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		
	case common.StartTimerMsg:
		m.page = common.TimerPage

		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd

	switch m.page {
	case common.HomePage:
		m.home, cmd = m.home.Update(msg)
	case common.SetTimerPage:
		m.setTimer, cmd = m.setTimer.Update(msg)
	case common.TimerPage:
		m.timer, cmd = m.timer.Update(msg)
	case common.BlockSitesPage:
		m.blockSites, cmd = m.blockSites.Update(msg)
	case common.BlockListPage:
		m.blockList, cmd = m.blockList.Update(msg)
	}

	return m, cmd
}

func (m *RootModel) View() string {
	var pageView string

	switch m.page {
	case common.HomePage:
		pageView = m.home.View()
	case common.SetTimerPage:
		pageView = m.setTimer.View()
	case common.TimerPage:
		pageView = m.timer.View()
	case common.BlockSitesPage:
		pageView = m.blockSites.View()
	case common.BlockListPage:
		pageView = m.blockList.View()
	}

	helpView := m.help.View(m.currentPageKeys())

	// Calculate how many lines the page and help take, and add spacing to push the help to the bottom
	pageLines := strings.Count(pageView, "\n") + 1
	helpLines := strings.Count(helpView, "\n") + 1

	gap := m.height - pageLines - helpLines
	if gap < 1 {
		gap = 1
	}

	return pageView + strings.Repeat("\n", gap) + helpView
}

// ================================================================
// Page Helpers
// ================================================================

func (m *RootModel) currentPageKeys() help.KeyMap {
	var pageModel any

	switch m.page {
	case common.HomePage:
		pageModel = m.home
	case common.SetTimerPage:
		pageModel = m.setTimer
	case common.TimerPage:
		pageModel = m.timer
	case common.BlockSitesPage:
		pageModel = m.blockSites
	case common.BlockListPage:
		pageModel = m.blockList
	}

	if pk, ok := pageModel.(PageKeys); ok {
		return pk.Keys()
	}

	return gKeys
}

// ================================================================
// Global Key Bindings
// ================================================================

type PageKeys interface {
	Keys() help.KeyMap
}

type globalKeys struct {
	Help key.Binding
	Quit key.Binding
}

func (k globalKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k globalKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}

var gKeys = globalKeys{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
