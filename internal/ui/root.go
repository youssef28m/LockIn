package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/ui/pages"
)

type Page int

const (
    HomePage Page = iota
    SetTimerPage
    TimerPage
    BlockSitesPage
)

type NavigateMsg Page

// globalKeys holds keybindings that work on every page.
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

// PageKeys is an interface each page model can optionally implement
// to expose its own keybindings to the root help footer.
type PageKeys interface {
	Keys() help.KeyMap
}

type RootModel struct {
    page       Page
    home       pages.HomeModel
    setTimer   pages.SetTimerModel
    timer      pages.TimerModel
    blockSites pages.BlockSitesModel
    help       help.Model
    width      int
	height     int
}

func NewRootModel() *RootModel {
    return &RootModel{
        page:       HomePage,
        home:       pages.NewHomeModel(),
        setTimer:   pages.NewSetTimerModel(),
        timer:      pages.NewTimerModel(),
        blockSites: pages.NewBlockSitesModel(),
        help:       help.New(),
    }
}

func (m *RootModel) Init() tea.Cmd { return nil }

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
 
        case NavigateMsg:
            m.page = Page(msg)
            return m, nil

        case tea.WindowSizeMsg:
		    m.width = msg.Width
            m.height = msg.Height
            m.help.Width = msg.Width
    }

    var cmd tea.Cmd
    switch m.page {
    case HomePage:
        m.home, cmd = m.home.Update(msg)
    case SetTimerPage:
        m.setTimer, cmd = m.setTimer.Update(msg)
    case TimerPage:
        m.timer, cmd = m.timer.Update(msg)
    case BlockSitesPage:
        m.blockSites, cmd = m.blockSites.Update(msg)
    }

    return m, cmd
}


// currentPageKeys returns the active page's key.Map if it implements PageKeys,
// otherwise falls back to just the global keys.
func (m *RootModel) currentPageKeys() help.KeyMap {
	var pageModel any
	switch m.page {
	case HomePage:
		pageModel = m.home
	case SetTimerPage:
		pageModel = m.setTimer
	case TimerPage:
		pageModel = m.timer
	case BlockSitesPage:
		pageModel = m.blockSites
	}

	if pk, ok := pageModel.(PageKeys); ok {
		return pk.Keys()
	}
	return gKeys
}

func (m *RootModel) View() string {

     var pageView string

    switch m.page {
    case HomePage:
        pageView = m.home.View()
    case SetTimerPage:
        pageView = m.setTimer.View()
    case TimerPage:
        pageView = m.timer.View()
    case BlockSitesPage:
        pageView = m.blockSites.View()
    }

    
    helpView := m.help.View(m.currentPageKeys())
	
    // Pin help to the bottom by filling the gap with newlines
	pageLines := strings.Count(pageView, "\n") + 1
	helpLines := strings.Count(helpView, "\n") + 1
	gap := m.height - pageLines - helpLines
	if gap < 1 {
		gap = 1
	}

    return pageView + strings.Repeat("\n", gap) + helpView
}

