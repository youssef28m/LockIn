package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

type listTab int

const (
	sitesTab listTab = iota
	appsTab
)

type BlockListModel struct {
	sites   []string
	apps    []string
	tab     listTab
	cursor  int
	err     error
	service *service.AppService
}

func NewBlockListModel(s *service.AppService) BlockListModel {
	return BlockListModel{
		service: s,
	}
}

func (b BlockListModel) Init() tea.Cmd {
	return nil
}

func (b BlockListModel) Update(msg tea.Msg) (BlockListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.BlockedListLoadedMsg:
		b.sites = msg.Sites
		b.err = nil
		return b, nil

	case common.BlockedAppsListLoadedMsg:
		b.apps = msg.Apps
		b.err = nil
		return b, nil

	case common.BlockedListErrorMsg:
		b.err = msg.Err
		return b, nil

	case common.BlockedAppsListErrorMsg:
		b.err = msg.Err
		return b, nil

	case tea.KeyMsg:
		b.err = nil

		switch msg.String() {
		case "left":
			if b.tab != sitesTab {
				b.tab = sitesTab
				b.cursor = 0
			}
			return b, nil

		case "right":
			if b.tab != appsTab {
				b.tab = appsTab
				b.cursor = 0
			}
			return b, nil

		case "up", "k":
			if b.cursor > 0 {
				b.cursor--
			}
			return b, nil

		case "down", "j":
			items := b.currentItems()
			if b.cursor < len(items)-1 {
				b.cursor++
			}
			return b, nil

		case "x", "backspace":
			if b.tab == sitesTab && len(b.sites) > 0 {
				name := b.sites[b.cursor]
				if err := b.service.RemoveBlockedSite(name); err != nil {
					b.err = err
					return b, nil
				}
				return b, common.FetchBlockedSitesCmd(*b.service)
			}
			if b.tab == appsTab && len(b.apps) > 0 {
				name := b.apps[b.cursor]
				if err := b.service.RemoveBlockedApp(name); err != nil {
					b.err = err
					return b, nil
				}
				return b, common.FetchBlockedAppsCmd(*b.service)
			}
			return b, nil

		case "esc":
			return b, common.NavigateTo(common.HomePage)
		}
	}

	return b, nil
}

func (b BlockListModel) currentItems() []string {
	if b.tab == sitesTab {
		return b.sites
	}
	return b.apps
}

func (b BlockListModel) View() string {
	var s strings.Builder

	s.WriteString(common.TitleStyle.Render("\nBlocked List\n"))

	activeTabStyle := lipgloss.NewStyle().
		Foreground(common.AccentColor).
		Bold(true).
		Padding(0, 1)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	s.WriteString("\n")
	if b.tab == sitesTab {
		s.WriteString(activeTabStyle.Render("● Websites"))
		s.WriteString(inactiveTabStyle.Render("Apps"))
	} else {
		s.WriteString(inactiveTabStyle.Render("Websites"))
		s.WriteString(activeTabStyle.Render("● Apps"))
	}
	s.WriteString("\n\n")

	items := b.currentItems()
	if len(items) > 0 {
		for i, item := range items {
			var itemStyle lipgloss.Style
			if i == b.cursor {
				itemStyle = lipgloss.NewStyle().
					Foreground(common.AccentColor).
					Bold(true).
					PaddingLeft(1).
					BorderLeft(true).
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(common.SuccessColor)
			} else {
				itemStyle = lipgloss.NewStyle().
					Foreground(common.PrimaryColor).
					PaddingLeft(2)
			}
			s.WriteString(itemStyle.Render(item) + "\n")
		}
	} else {
		s.WriteString(common.LabelStyle.Render(" (none)\n"))
	}

	if b.err != nil {
		s.WriteString("\n" + common.ErrorStyle.Render(" " + b.err.Error()))
	}

	return s.String()
}

// ================================================================
// Key Bindings
// ================================================================

type BlockListKeyMap struct {
	Help   key.Binding
	Quit   key.Binding
	Switch key.Binding
	Delete key.Binding
	Back   key.Binding
}

func (k BlockListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Switch, k.Delete, k.Back, k.Help, k.Quit}
}

func (k BlockListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Switch, k.Delete},
		{k.Back, k.Help, k.Quit},
	}
}

var blockListKeys = BlockListKeyMap{
	Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Switch: key.NewBinding(key.WithKeys("left", "right"), key.WithHelp("←/→", "switch tab")),
	Delete: key.NewBinding(key.WithKeys("x", "backspace"), key.WithHelp("x", "delete")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

func (m BlockListModel) Keys() help.KeyMap {
	return blockListKeys
}
