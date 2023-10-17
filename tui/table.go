package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"io"
)

const (
	cyan   = lipgloss.Color("14")
	purple = lipgloss.Color("93")
	white  = lipgloss.Color("252")
)

func NewTable(out io.Writer) *table.Table {
	re := lipgloss.NewRenderer(out)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(cyan).Bold(true).Align(lipgloss.Center)

	return table.New().Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			case row%2 == 0:
				return baseStyle.Copy().Foreground(lipgloss.Color("245"))
			default:
				return baseStyle.Copy().Foreground(lipgloss.Color("252"))
			}
		})
}
