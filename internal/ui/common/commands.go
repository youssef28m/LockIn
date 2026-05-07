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