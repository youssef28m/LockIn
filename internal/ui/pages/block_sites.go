package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type BlockSitesModel struct{ cursor int }

func NewBlockSitesModel() BlockSitesModel { return BlockSitesModel{} }

func (m BlockSitesModel) Init() tea.Cmd { return nil }

func (m BlockSitesModel) Update(msg tea.Msg) (BlockSitesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.cursor++
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m BlockSitesModel) View() string {
	return fmt.Sprintf("🔒 Block Sites\n\nCursor: %d\n\nPress q → Quit", m.cursor)
}
