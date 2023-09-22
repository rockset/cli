package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Selector struct {
	values   []string
	cursor   int
	Selected int
}

func NewSelector(values []string, selected int) *Selector {
	return &Selector{values, selected, selected}
}

func (s *Selector) Init() tea.Cmd {
	return textinput.Blink
}

func (s *Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			s.Selected = -1
			return s, tea.Quit
		case "down":
			if s.cursor < len(s.values)-1 {
				s.cursor++
			}
		case "up":
			if s.cursor > 0 {
				s.cursor--
			}
		case "enter":
			s.Selected = s.cursor
			return s, tea.Quit
		}
	}

	return s, nil
}

func (s *Selector) View() string {
	var b strings.Builder

	b.WriteString("Select:\n")
	for i, v := range s.values {
		if i == s.cursor {
			b.WriteString(focusedStyle.Render("-> "))
		} else {
			b.WriteString(blurredStyle.Render("-  "))
		}
		if i == s.Selected {
			b.WriteString(focusedStyle.Render(v))
		} else {
			b.WriteString(v)
		}

		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("\nEnter to select, Q to quit\n"))

	return b.String()
}
