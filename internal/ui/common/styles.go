package common

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor   = lipgloss.Color("#4361EE") // Royal Blue
	AccentColor    = lipgloss.Color("#A8DADC") // Ice Blue
	SuccessColor   = lipgloss.Color("#dc596f") // Dusty Lavender
	ErrorColor = lipgloss.Color("#E63946") // Soft red


	TitleStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			MarginLeft(2).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(PrimaryColor).
			PaddingLeft(1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	SelectedItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
    Foreground(ErrorColor).
    Bold(true)

)