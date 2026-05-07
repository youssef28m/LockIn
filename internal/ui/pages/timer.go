package pages

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

func bellCmd() tea.Cmd {
	return func() tea.Msg {
		for i := 0; i < 3; i++ {
			os.Stderr.WriteString("\a")
			time.Sleep(200 * time.Millisecond)
		}
		return nil
	}
}

// ================================================================
// Key Bindings
// ================================================================

type timerRunningKeyMap struct {
	Help key.Binding
	Quit key.Binding
}

func (k timerRunningKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k timerRunningKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Help, k.Quit}}
}

var timerRunningKeys = timerRunningKeyMap{
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

type timerDoneKeyMap struct {
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
}

func (k timerDoneKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Help, k.Quit}
}

func (k timerDoneKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Enter}, {k.Help, k.Quit}}
}

var timerDoneKeys = timerDoneKeyMap{
	Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Enter: key.NewBinding(key.WithKeys("enter", "esc"), key.WithHelp("enter", "back to home")),
}

func (m TimerModel) Keys() help.KeyMap {
	if m.done {
		return timerDoneKeys
	}
	return timerRunningKeys
}

// TickMsg is sent every second to update the countdown
type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TimerModel struct {
	TotalSeconds int
	Remaining    int
	Running      bool
	done         bool
	service      *service.AppService
}

func NewTimerModel(s *service.AppService) TimerModel {
	return TimerModel{
		service: s,
	}
}

// Init doesn't start the timer immediately; 
// it starts when the RootModel receives StartTimerMsg
func (m TimerModel) Init() tea.Cmd {
	return nil
}

func (m TimerModel) Update(msg tea.Msg) (TimerModel, tea.Cmd) {
	switch msg := msg.(type) {

	case common.StartTimerMsg:
		m.TotalSeconds = msg.DurationSeconds
		m.Remaining = msg.DurationSeconds
		m.Running = true
		m.done = false
		return m, tick()

	case TickMsg:
		if !m.Running || m.done {
			return m, nil
		}

		session, err := m.service.GetActiveSession()
		if err != nil {
			m.Running = false
			return m, nil
		}

		m.Remaining = int(session.Remaining())
		if m.Remaining <= 0 {
			m.Running = false
			m.Remaining = 0
			if !m.done {
				m.done = true
				return m, bellCmd()
			}
		}

		return m, tick()

	case tea.KeyMsg:
		if m.done {
			switch msg.String() {
			case "enter", "esc":
				return m, common.NavigateTo(common.HomePage)
			}
		}
	}
	return m, nil
}

func (m TimerModel) View() string {
	if m.done {
		completion := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(60).
			Foreground(common.SuccessColor).
			Bold(true).
			Render("Time's Up! Focus Session Complete")

		instruction := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(60).
			Foreground(common.PrimaryColor).
			Render("Press Enter to continue")

		return "\n\n" + completion + "\n\n" + instruction + "\n"
	}

	hours := m.Remaining / 3600
	minutes := (m.Remaining % 3600) / 60
	seconds := m.Remaining % 60

	timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	title := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(60).
		Render(common.TitleStyle.Render("Focus Session"))

	bigTimer := common.RenderBigTime(timeStr)

	centeredTimer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(60).
		Render(bigTimer)

	return fmt.Sprintf(
		"\n%s\n\n%s\n",
		title,
		centeredTimer,
	)
}