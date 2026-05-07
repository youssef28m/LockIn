package pages

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/youssef28m/LockIn/internal/models"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

type HistoryModel struct {
	records []models.Session
	err     error
	service *service.AppService
}

func NewHistoryModel(s *service.AppService) HistoryModel {
	return HistoryModel{
		service: s,
	}
}

func (m HistoryModel) Init() tea.Cmd { return nil }

func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.HistoryLoadedMsg:
		m.records = msg.Sessions
		m.err = nil
		return m, nil

	case common.HistoryErrorMsg:
		m.err = msg.Err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, common.NavigateTo(common.HomePage)
		}
	}

	return m, nil
}

func (m HistoryModel) View() string {
	var b strings.Builder

	b.WriteString(common.TitleStyle.Render("\nSession History\n"))

	if m.err != nil {
		b.WriteString("\n" + common.ErrorStyle.Render(" " + m.err.Error()))
		return b.String()
	}

	if len(m.records) == 0 {
		b.WriteString("\n" + common.LabelStyle.Render(" No completed sessions yet.\n Start a focus session to see your history here."))
		return b.String()
	}

	sorted := make([]models.Session, len(m.records))
	copy(sorted, m.records)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].StartTime > sorted[j].StartTime
	})

	var totalSeconds int64
	for _, r := range sorted {
		totalSeconds += r.DurationSeconds
	}

	totalStr := formatDuration(totalSeconds)
	b.WriteString(fmt.Sprintf(
		"\n%s\n\n",
		common.LabelStyle.Render(fmt.Sprintf("Total: %s across %d sessions", totalStr, len(sorted))),
	))

	headerStyle := lipgloss.NewStyle().
		Foreground(common.PrimaryColor).
		Bold(true)
	b.WriteString(headerStyle.Render(fmt.Sprintf("  %-14s %s", "Date", "Duration")) + "\n")

	limit := 50
	if len(sorted) < limit {
		limit = len(sorted)
	}

	for i, r := range sorted[:limit] {
		dateStr := time.Unix(r.StartTime, 0).Format("Jan 02  15:04")
		durStr := formatDuration(r.DurationSeconds)

		var lineStyle lipgloss.Style
		if i%2 == 0 {
			lineStyle = lipgloss.NewStyle().Foreground(common.AccentColor)
		} else {
			lineStyle = lipgloss.NewStyle().Foreground(common.PrimaryColor)
		}
		b.WriteString(lineStyle.Render(fmt.Sprintf("  %-14s %s", dateStr, durStr)) + "\n")
	}

	if len(sorted) > limit {
		b.WriteString(common.LabelStyle.Render(fmt.Sprintf("  ... and %d more\n", len(sorted)-limit)))
	}

	return b.String()
}

func formatDuration(seconds int64) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 && m > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", m)
}

// ================================================================
// Key Bindings
// ================================================================

type historyKeyMap struct {
	Help key.Binding
	Quit key.Binding
	Back key.Binding
}

func (k historyKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Help, k.Quit}
}

func (k historyKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back},
		{k.Help, k.Quit},
	}
}

var historyKeys = historyKeyMap{
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

func (m HistoryModel) Keys() help.KeyMap {
	return historyKeys
}
