package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/service"
)

func FetchBlockedSitesCmd(s service.AppService) tea.Cmd {
	return func() tea.Msg {
		sites, err := s.GetBlockedSites()
		if err != nil {
			return BlockedListErrorMsg{Err: err}
		}
		return BlockedListLoadedMsg{Sites: sites}
	}
}

func FetchBlockedAppsCmd(s service.AppService) tea.Cmd {
	return func() tea.Msg {
		apps, err := s.GetBlockedApps()
		if err != nil {
			return BlockedAppsListErrorMsg{Err: err}
		}
		return BlockedAppsListLoadedMsg{Apps: apps}
	}
}