package common

import "github.com/charmbracelet/lipgloss"

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true).
			MarginLeft(2)

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#005F87"))

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
)