package pages

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

// 1. Create a custom item type for the list
type siteItem string

func (s siteItem) FilterValue() string { return string(s) }
func (s siteItem) Title() string       { return string(s) }
func (s siteItem) Description() string { return "" }

type BlockListModel struct {
	list    list.Model
	service *service.AppService
	err     error
}

func NewBlockListModel(s *service.AppService) BlockListModel {

	delegate := list.NewDefaultDelegate()

	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	// Style the normal (unselected) items
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(common.PrimaryColor).
		PaddingLeft(2).
		MarginLeft(0).
        MarginBottom(1)

	// Style the selected item with accent color and left border indicator
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(common.AccentColor).
		Bold(true).
		PaddingLeft(1).
		MarginLeft(0).
        MarginBottom(1).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(common.SuccessColor)

	// Style the dimmed (when not focused) title
	delegate.Styles.DimmedTitle = delegate.Styles.DimmedTitle.
		Foreground(common.PrimaryColor).
		Faint(true).
		PaddingLeft(2)

	l := list.New([]list.Item{}, delegate, 20, 12)

	l.Title = "🔒 Blocked Sites"

	l.Styles.Title = common.TitleStyle
	l.Styles.PaginationStyle = common.LabelStyle
	l.Styles.HelpStyle = common.LabelStyle

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return BlockListModel{
		list:    l,
		service: s,
	}
}

func (b BlockListModel) Init() tea.Cmd {
	return nil
}

func (b BlockListModel) Update(msg tea.Msg) (BlockListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.BlockedListLoadedMsg:
		items := make([]list.Item, len(msg.Sites))
		for i, site := range msg.Sites {
			items[i] = siteItem(site)
		}
		// Update the list's data
		b.list.SetItems(items)

	case tea.KeyMsg:
		switch msg.String() {
		case "x", "backspace":
			if i, ok := b.list.SelectedItem().(siteItem); ok {
				// Remove the site from the service
				siteToDelete := string(i)
				err := b.service.RemoveBlockedSite(siteToDelete)
				if err != nil {
					b.err = err
					return b, nil
				}

				b.list.RemoveItem(b.list.Index())

				// refresh the list from the service to ensure it's up to date
				return b, common.FetchBlockedSitesCmd(*b.service)
			}
		case "esc":
			return b, common.NavigateTo(common.HomePage)
		}
	}

	var cmd tea.Cmd
	b.list, cmd = b.list.Update(msg)

	return b, cmd
}

func (b BlockListModel) View() string {
	return b.list.View()
}

// ================================================================
// Key Bindings
// ================================================================

type BlockListKeyMap struct {
	Help   key.Binding
	Quit   key.Binding
	Delete key.Binding
	Back   key.Binding
}

// ShortHelp returns bindings to show in the small help view
func (k BlockListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Delete, k.Back, k.Help, k.Quit}
}

// FullHelp returns bindings for the expanded help view (when '?' is pressed)
func (k BlockListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Delete, k.Back},
		{k.Help, k.Quit},
	}
}

var blockListKeys = BlockListKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x", "backspace"),
		key.WithHelp("x", "delete site"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

func (m BlockListModel) Keys() help.KeyMap {
	return blockListKeys
}