package pages

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

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
	err 		error
	service	  *service.AppService
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
		return m, tick()

	case TickMsg:
		if !m.Running || m.Remaining <= 0 {
			m.Running = false
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
		}

		return m, tick()
	}		
	return m, nil
}

func (m TimerModel) View() string {
	if m.Remaining <= 0 && m.TotalSeconds > 0 {
		completion := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(60).
			Foreground(common.SuccessColor).
			Bold(true).
			Render("🎉 Time's Up! Focus Session Complete 🎉")
		
		return "\n\n" + completion + "\n\n"
	}

	hours := m.Remaining / 3600
	minutes := (m.Remaining % 3600) / 60
	seconds := m.Remaining % 60

	// Format the time string
	timeStr := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	
	// Title
	title := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(60).
		Render(common.TitleStyle.Render("Focus Session"))
	
	// Create big ASCII art timer
	bigTimer := common.RenderBigTime(timeStr)
	
	// Center the timer
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