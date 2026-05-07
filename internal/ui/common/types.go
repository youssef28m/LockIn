package common

import tea "github.com/charmbracelet/bubbletea"

type Page int

const (
    HomePage Page = iota
    SetTimerPage
    TimerPage
    BlockSitesPage
    BlockListPage
    BlockAppsPage
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

// The sub-models can return this without knowing about the Root
type NavigateMsg struct {
    Target Page
}

func NavigateTo(p Page) tea.Cmd {
    return func() tea.Msg {
        return NavigateMsg{Target: p}
    }
}