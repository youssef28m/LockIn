package common

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryColor   = lipgloss.Color("#4361EE") // Royal Blue
	AccentColor    = lipgloss.Color("#A8DADC") // Ice Blue
	SuccessColor   = lipgloss.Color("#dc596f") // Dusty Lavender
	ErrorColor     = lipgloss.Color("#E63946") // Soft red

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

	// Big clock digit style
	BigDigitStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)
)

// DigitMap - ASCII art representations of digits (5 lines tall)
var DigitMap = map[rune][]string{
	'0': {
		" ███ ",
		"█   █",
		"█   █",
		"█   █",
		" ███ ",
	},
	'1': {
		"  █  ",
		" ██  ",
		"  █  ",
		"  █  ",
		" ███ ",
	},
	'2': {
		" ███ ",
		"    █",
		" ███ ",
		"█    ",
		"█████",
	},
	'3': {
		"████ ",
		"    █",
		" ███ ",
		"    █",
		"████ ",
	},
	'4': {
		"█   █",
		"█   █",
		"█████",
		"    █",
		"    █",
	},
	'5': {
		"█████",
		"█    ",
		"████ ",
		"    █",
		"████ ",
	},
	'6': {
		" ███ ",
		"█    ",
		"████ ",
		"█   █",
		" ███ ",
	},
	'7': {
		"█████",
		"    █",
		"   █ ",
		"  █  ",
		"  █  ",
	},
	'8': {
		" ███ ",
		"█   █",
		" ███ ",
		"█   █",
		" ███ ",
	},
	'9': {
		" ███ ",
		"█   █",
		" ████",
		"    █",
		" ███ ",
	},
	':': {
		"     ",
		"  █  ",
		"     ",
		"  █  ",
		"     ",
	},
}

// RenderBigTime creates a large ASCII art representation of time
func RenderBigTime(timeStr string) string {
	lines := make([]string, 5)
	
	for _, char := range timeStr {
		digit, exists := DigitMap[char]
		if !exists {
			continue
		}
		
		for i := 0; i < 5; i++ {
			lines[i] += digit[i] + " "
		}
	}
	
	// Join lines and apply style
	result := ""
	for _, line := range lines {
		result += BigDigitStyle.Render(line) + "\n"
	}
	
	return result
}