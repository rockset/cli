package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	Red    = lipgloss.Color("9")
	Purple = lipgloss.Color("13")
	Cyan   = lipgloss.Color("14")
	Yellow = lipgloss.Color("11")

	focusedStyle = lipgloss.NewStyle().Foreground(Purple)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ErrorStyle   = lipgloss.NewStyle().Foreground(Red)
	WarningStyle = lipgloss.NewStyle().Foreground(Yellow)
	BracketStyle = lipgloss.NewStyle().Foreground(Cyan)
	RocksetStyle = lipgloss.NewStyle().Foreground(Purple)
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()

	focusedButton = fmt.Sprintf("%s %s %s",
		BracketStyle.Render("["),
		RocksetStyle.Render("Submit"),
		BracketStyle.Render("]"),
	)
	blurredButton = fmt.Sprintf("%s %s %s",
		BracketStyle.Render("["),
		blurredStyle.Render("Submit"),
		BracketStyle.Render("]"),
	)
)

var (
	R       = Bracketed("R")
	Rockset = Bracketed("Rockset")

	Prompt             = R + "> "
	ContinuationPrompt = BracketStyle.Render(">") + RocksetStyle.Render(">") + "> "
)

func Bracketed(msg string) string {
	return BracketStyle.Render("[") + RocksetStyle.Render(msg) + BracketStyle.Render("]")
}
