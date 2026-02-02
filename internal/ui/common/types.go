package common

import tea "github.com/charmbracelet/bubbletea"

type Page int

const (
    HomePage Page = iota
    SetTimerPage
    TimerPage
    BlockSitesPage
)

type StartTimerMsg struct {
    DurationSeconds int
}

// The sub-models can return this without knowing about the Root
type NavigateMsg struct {
    Target Page
}

func Goto(p Page) tea.Cmd {
    return func() tea.Msg {
        return NavigateMsg{Target: p}
    }
}