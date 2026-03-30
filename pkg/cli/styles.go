package cli

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorGreen  = lipgloss.Color("10")
	colorRed    = lipgloss.Color("9")
	colorYellow = lipgloss.Color("11")
	colorBlue   = lipgloss.Color("12")
	colorDim    = lipgloss.Color("8")
	colorWhite  = lipgloss.Color("15")

	greenStyle  = lipgloss.NewStyle().Foreground(colorGreen)
	redStyle    = lipgloss.NewStyle().Foreground(colorRed)
	yellowStyle = lipgloss.NewStyle().Foreground(colorYellow)
	blueStyle   = lipgloss.NewStyle().Foreground(colorBlue)
	dimStyle    = lipgloss.NewStyle().Foreground(colorDim)
	boldStyle   = lipgloss.NewStyle().Bold(true)
	whiteStyle  = lipgloss.NewStyle().Foreground(colorWhite)

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	promptStyle = lipgloss.NewStyle().Bold(true).Foreground(colorGreen)
	keyStyle    = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	labelStyle  = lipgloss.NewStyle().Foreground(colorDim)
	valueStyle  = lipgloss.NewStyle().Foreground(colorWhite)
)

func renderEnabled(b bool) string {
	if b {
		return lipgloss.NewStyle().Foreground(colorGreen).Render("● enabled")
	}
	return lipgloss.NewStyle().Foreground(colorRed).Render("● disabled")
}

func separator(w int) string {
	return dimStyle.Render(strings.Repeat("─", w))
}

func stringsRepeat(s string, n int) string {
	return strings.Repeat(s, n)
}
