package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/models"
)

type Page int

const (
	HomePage Page = iota
	SetTimerPage
	TimerPage
	BlockSitesPage
	BlockListPage
	BlockAppsPage
	HistoryPage
)

type StartTimerMsg struct {
	DurationSeconds int
}

type BlockedListLoadedMsg struct {
	Sites []string
}

type BlockedListErrorMsg struct {
	Err error
}

type BlockedAppsListLoadedMsg struct {
	Apps []string
}

type BlockedAppsListErrorMsg struct {
	Err error
}

type HistoryLoadedMsg struct {
	Sessions []models.Session
}

type HistoryErrorMsg struct {
	Err error
}

type NavigateMsg struct {
	Target Page
}

func NavigateTo(p Page) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{Target: p}
	}
}