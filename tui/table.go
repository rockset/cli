package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"io"
)

func NewTable(out io.Writer) *table.Table {
	re := lipgloss.NewRenderer(out)
	baseStyle := re.NewStyle().Padding(0, 1).Align(lipgloss.Left)
	headerStyle := baseStyle.Copy().Foreground(Cyan).Bold(true).Align(lipgloss.Center)
	oddStyle := baseStyle.Copy().Foreground(Grey)
	evenStyle := baseStyle.Copy().Foreground(White)

	return table.New().Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(Purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			case row%2 == 0:
				return oddStyle
			default:
				return evenStyle
			}
		})
}
