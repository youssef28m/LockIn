package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

// ================================================================
// Types
// ================================================================

type timerFocus int

const (
	focusHours timerFocus = iota
	focusMinutes
)

// ================================================================
// Model
// ================================================================

type SetTimerModel struct {
	hours   int
	minutes int
	service *service.AppService
	focus   timerFocus
	err     string
}

func NewSetTimerModel(service *service.AppService) SetTimerModel {
	return SetTimerModel{
		minutes: 25, // Default Pomodoro length
		service: service,
	}
}

// ================================================================
// Bubble Tea Lifecycle
// ================================================================

func (m SetTimerModel) Init() tea.Cmd {
	return nil
}

func (m SetTimerModel) Update(msg tea.Msg) (SetTimerModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		m.err = ""

		switch msg.String() {

		case "left", "right", "tab":
			if m.focus == focusHours {
				m.focus = focusMinutes
			} else {
				m.focus = focusHours
			}

		case "up", "k":
			if m.focus == focusHours {
				if m.hours < 24 {
					m.hours++
				}
			} else {
				if m.minutes < 59 {
					m.minutes++
				} else {
					m.minutes = 0
				}
			}

		case "down", "j":
			if m.focus == focusHours {
				if m.hours > 0 {
					m.hours--
				}
			} else {
				if m.minutes > 0 {
					m.minutes--
				} else {
					m.minutes = 59
				}
			}

		case "esc":
			return m, common.NavigateTo(common.HomePage)

		case "enter":
			totalSeconds := (m.hours * 3600) + (m.minutes * 60)
			if totalSeconds > 0 {
				if err := m.service.CreateAndStartSession(totalSeconds); err != nil {
					m.err = err.Error()
					return m, nil
				}

				return m, func() tea.Msg {
					return common.StartTimerMsg{
						DurationSeconds: totalSeconds,
					}
				}
			}
		}
	}

	return m, nil
}

func (m SetTimerModel) View() string {
	title := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(60).
		Render(common.TitleStyle.Render("Set Focus Duration"))

	timeStr := fmt.Sprintf("%02d:%02d", m.hours, m.minutes)

	bigTimeLines := make([]string, 5)

	for i, char := range timeStr {
		digit, exists := common.DigitMap[char]
		if !exists {
			continue
		}

		var style lipgloss.Style

		if i < 2 && m.focus == focusHours {
			style = common.BigDigitStyle.Foreground(common.PrimaryColor)
		} else if i < 2 {
			style = common.BigDigitStyle.Foreground(common.AccentColor)
		} else if i == 2 {
			style = common.BigDigitStyle.Foreground(lipgloss.Color("#FFFFFF"))
		} else if m.focus == focusMinutes {
			style = common.BigDigitStyle.Foreground(common.PrimaryColor)
		} else {
			style = common.BigDigitStyle.Foreground(common.AccentColor)
		}

		for line := 0; line < 5; line++ {
			bigTimeLines[line] += style.Render(digit[line]) + " "
		}
	}

	bigTime := ""
	for _, line := range bigTimeLines {
		bigTime += line + "\n"
	}

	centeredTimer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(60).
		Render(bigTime)

	errorMsg := ""
	if m.err != "" {
		errorMsg = "\n" + lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(60).
			Render(common.ErrorStyle.Render(m.err)) + "\n"
	}

	startBtn := lipgloss.NewStyle().
		Foreground(common.SuccessColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("▶ Press Enter to Start")

	return fmt.Sprintf(
		"\n%s\n\n%s\n%s%s\n",
		title,
		centeredTimer,
		errorMsg,
		startBtn,
	)
}

// ================================================================
// Key Bindings
// ================================================================

func (m SetTimerModel) Keys() help.KeyMap {
	return timerKeyMap
}

type timerKeyBindings struct {
	Up     key.Binding
	Down   key.Binding
	Switch key.Binding
	Start  key.Binding
	Exit   key.Binding
}

func (k timerKeyBindings) ShortHelp() []key.Binding {
	return []key.Binding{k.Switch, k.Up, k.Down, k.Start, k.Exit}
}

func (k timerKeyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Switch, k.Up, k.Down, k.Start, k.Exit},
	}
}

var timerKeyMap = timerKeyBindings{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "increase")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "decrease")),
	Switch: key.NewBinding(key.WithKeys("left", "right", "tab"), key.WithHelp("tab", "switch H/M")),
	Start:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start")),
	Exit:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back to home")),
}
