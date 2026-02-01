package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type SetTimerModel struct{ minutes int }

func NewSetTimerModel() SetTimerModel { return SetTimerModel{} }

func (m SetTimerModel) Init() tea.Cmd { return nil }

func (m SetTimerModel) Update(msg tea.Msg) (SetTimerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.minutes++
		case "down":
			if m.minutes > 0 {
				m.minutes--
			}
		}
	}
	return m, nil
}

func (m SetTimerModel) View() string {
	return fmt.Sprintf("⏲ Set Timer Page\n\nMinutes: %d\n\nPress q → Quit", m.minutes)
}
