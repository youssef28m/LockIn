package pages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)



// implement help keys
type homeKeys struct {
	Help key.Binding
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
}

func (k homeKeys) ShortHelp() []key.Binding {
    return []key.Binding{k.Help, k.Quit}
}

func (k homeKeys) FullHelp() [][]key.Binding {
    return [][]key.Binding{
		{k.Up, k.Down},
        {k.Help ,k.Quit},
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

type HomeModel struct {
	cursor int
	choice string
	keys   homeKeys
}

func (m HomeModel) Keys() help.KeyMap {
    return hKeys
}


func NewHomeModel() HomeModel { return HomeModel{} }

func (m HomeModel) Init() tea.Cmd { return nil }

var choices = []string{"Add website to block list", "Set Timer"}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.choice = choices[m.cursor]
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

	b.WriteString("\n🔒  LockIn\n")
	b.WriteString("====================\n\n")

	
	for i, c := range choices {
		cursor := "   "
		if m.cursor == i {
			cursor = "➜  "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, c))
	}


	return b.String()
}
