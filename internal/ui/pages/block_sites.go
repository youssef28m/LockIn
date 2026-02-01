package pages

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

//================= BlockSites Keys =================//

type blockKeys struct {
	Help key.Binding
	Quit key.Binding
	Return key.Binding
	Enter key.Binding
}

func (k blockKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k blockKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Return, k.Enter},
		{k.Help, k.Quit},
	}
}

var bKeys = blockKeys{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "return to home"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "add site"),
	),
}

func (m BlockSitesModel) Keys() help.KeyMap {
    return bKeys
}

//================= BlockSites Model =================//

type BlockSitesModel struct {
    textInput textinput.Model
    err       error
    service   *service.AppService
}


func NewBlockSitesModel(s *service.AppService) BlockSitesModel {
    ti := textinput.New()
    ti.Placeholder = "example.com"
    ti.Focus() // Start with the input active
    ti.CharLimit = 64
    ti.Width = 30
    
    // Style the input components
    ti.PromptStyle = lipgloss.NewStyle().Foreground(common.PrimaryColor)
    ti.TextStyle = lipgloss.NewStyle().Foreground(common.AccentColor)
    ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

    return BlockSitesModel{
        textInput: ti,
        service: s,
    }
}

func (m BlockSitesModel) Init() tea.Cmd { return nil }

func (m BlockSitesModel) Update(msg tea.Msg) (BlockSitesModel, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            domain := m.textInput.Value()
            err := m.service.AddBlockedSite(domain)
            if err != nil {
                m.err = err
            } else {
                m.err = nil
            }
            
			m.textInput.Reset()
            return m, nil
            
        case "esc":
            return m, common.Goto(common.HomePage)
        }
    }

    // This line is crucial: it updates the cursor and text state
    m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd
}

var (
    inputContainerStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(common.PrimaryColor).
        Padding(1, 1).
        MarginTop(1)
)

func (m BlockSitesModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        common.TitleStyle.Render("\nAdd Website to block list\n"),
        "\nEnter the domain you want to restrict:",
        inputContainerStyle.Render(m.textInput.View()),
    )
}
