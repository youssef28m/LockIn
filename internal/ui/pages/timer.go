package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type TimerModel struct{ elapsed int }

func NewTimerModel() TimerModel { return TimerModel{} }

func (m TimerModel) Init() tea.Cmd { return nil }

func (m TimerModel) Update(msg tea.Msg) (TimerModel, tea.Cmd) {
	// placeholder for timer updates
	return m, nil
}

func (m TimerModel) View() string {
	return fmt.Sprintf("⏱ Timer Page\n\nElapsed: %d\n\nPress q → Quit", m.elapsed)
}
