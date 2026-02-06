package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/youssef28m/LockIn/internal/service"
	"github.com/youssef28m/LockIn/internal/ui/common"
)

type BlockListModel struct {
	sites  []string
	service service.AppService
}


func NewSiteListModel(s service.AppService) BlockListModel {
    return BlockListModel{
        service: s,
    }
}

func (b BlockListModel) Init() tea.Cmd { 
    return nil
}

func (b BlockListModel) Update(msg tea.Msg) (BlockListModel, tea.Cmd) {
    switch msg := msg.(type) {
    case common.BlockedListLoadedMsg:
        b.sites = msg.Sites
    }

    return b, nil
}

func (b BlockListModel) View() string {
    return "Blocked Sites"
}

func fetchBlockedSitesCmd(s service.AppService) tea.Cmd {
    return func() tea.Msg {
        sites, err := s.GetBlockedSites()
        if err != nil {
            return err
        }
        return common.BlockedListLoadedMsg{Sites: sites}
    }
}
