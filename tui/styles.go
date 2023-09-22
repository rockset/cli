package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	BracketStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	RocksetStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
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
