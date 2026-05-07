package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

type BlockAppsModel struct {
	textInput  textinput.Model
	err        error
	successMsg string
	service    *service.AppService
}

func NewBlockAppsModel(s *service.AppService) BlockAppsModel {
	ti := textinput.New()
	ti.Placeholder = "firefox"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 30
	ti.PromptStyle = lipgloss.NewStyle().Foreground(common.PrimaryColor)
	ti.TextStyle = lipgloss.NewStyle().Foreground(common.AccentColor)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return BlockAppsModel{
		textInput: ti,
		service:   s,
	}
}

func (m BlockAppsModel) Init() tea.Cmd { return nil }

func (m BlockAppsModel) Update(msg tea.Msg) (BlockAppsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.err = nil
		m.successMsg = ""

		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(m.textInput.Value())
			if name == "" {
				return m, nil
			}
			err := m.service.AddBlockedApp(name)
			if err != nil {
				m.err = err
			} else {
				m.successMsg = "✓ " + name + " added"
			}
			m.textInput.Reset()
			return m, nil

		case "esc":
			return m, common.NavigateTo(common.HomePage)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m BlockAppsModel) View() string {
	var b strings.Builder

	b.WriteString(common.TitleStyle.Render("\nBlock Apps\n"))
	b.WriteString("\nEnter the process name to block (e.g. firefox):\n\n")

	inputBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(common.PrimaryColor).
		Padding(1, 1)

	b.WriteString(inputBorder.Render(m.textInput.View()) + "\n")

	if m.successMsg != "" {
		b.WriteString("\n " + lipgloss.NewStyle().Foreground(common.SuccessColor).Render(m.successMsg) + "\n")
	}

	if m.err != nil {
		b.WriteString("\n" + common.ErrorStyle.Render(" " + m.err.Error()))
	}

	return b.String()
}

// ================================================================
// Key Bindings
// ================================================================

type blockAppsKeyMap struct {
	Help key.Binding
	Quit key.Binding
	Add  key.Binding
	Back key.Binding
}

func (k blockAppsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Back, k.Help, k.Quit}
}

func (k blockAppsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Back},
		{k.Help, k.Quit},
	}
}

var blockAppsKeys = blockAppsKeyMap{
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Add:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "add app")),
	Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

func (m BlockAppsModel) Keys() help.KeyMap {
	return blockAppsKeys
}
