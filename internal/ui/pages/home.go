package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

// ================================================================
// Model
// ================================================================

type HomeModel struct {
	cursor  int
	choice  string
	keys    homeKeys
	service *service.AppService
}

func NewHomeModel(s *service.AppService) HomeModel {
	return HomeModel{
		service: s,
	}
}

// ================================================================
// Bubble Tea Lifecycle
// ================================================================

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			m.choice = choices[m.cursor]

			switch m.choice {
			case "Add website to block list":
				return m, common.NavigateTo(common.BlockSitesPage)

			case "Set Timer":
				return m, common.NavigateTo(common.SetTimerPage)

			case "Show Block list":
				return m, common.NavigateTo(common.BlockListPage)
			}

		case "q", "Q", "ctrl+c":
			return m, tea.Quit

		case "down", "k":
			if m.cursor < len(choices)-1 {
				m.cursor++
			}

		case "up", "j":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}

	return m, nil
}

func (m HomeModel) View() string {
	var b strings.Builder

	title := figure.NewFigure("LockIn", "", true)
	b.WriteString(common.TitleStyle.Render(title.String()) + "\n")

	for i, c := range choices {
		cursor := "  "
		label := c

		if m.cursor == i {
			cursor = common.CursorStyle.Render("➜ ")
			label = common.SelectedItemStyle.Render(c)
		}

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, label))
	}

	return b.String()
}

// ================================================================
// Data
// ================================================================

var choices = []string{
	"Add website to block list",
	"Show Block list",
	"Set Timer",
}

// ================================================================
// Key Bindings
// ================================================================

type homeKeys struct {
	Help key.Binding
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
}

func (m HomeModel) Keys() help.KeyMap {
	return hKeys
}

func (k homeKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k homeKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Help, k.Quit},
	}
}

var hKeys = homeKeys{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
}
