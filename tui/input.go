package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type Input struct {
	title      string
	focusIndex int
	inputs     []textinput.Model
	Fields     []string
	Err        error
	cursorMode cursor.Mode
}

type InputConfig struct {
	Prompt      string
	Placeholder string
	Validate    textinput.ValidateFunc
}

func NewInput(title string, cfg []InputConfig) *Input {
	m := Input{
		title:  title,
		inputs: make([]textinput.Model, len(cfg)),
	}

	for i, c := range cfg {
		t := textinput.New()
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		t.Cursor.Style = cursorStyle
		t.CharLimit = 64
		t.Placeholder = c.Placeholder
		t.Prompt = c.Prompt
		t.Validate = c.Validate

		m.inputs[i] = t
	}

	return &m
}

func (i *Input) Init() tea.Cmd {
	return textinput.Blink
}

func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return i, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused? If so, exit.
			if s == "enter" && i.focusIndex == len(i.inputs) {
				for _, input := range i.inputs {
					if input.Err != nil {
						// TODO this only reports the last error, should it use a slice of error instead?
						i.Err = input.Err
					}
					i.Fields = append(i.Fields, input.Value())
				}

				return i, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				i.focusIndex--
			} else {
				i.focusIndex++
			}

			if i.focusIndex > len(i.inputs) {
				i.focusIndex = 0
			} else if i.focusIndex < 0 {
				i.focusIndex = len(i.inputs)
			}

			cmds := make([]tea.Cmd, len(i.inputs))
			for j := 0; j <= len(i.inputs)-1; j++ {
				if j == i.focusIndex {
					// Set focused state
					cmds[j] = i.inputs[j].Focus()
					i.inputs[j].PromptStyle = focusedStyle
					i.inputs[j].TextStyle = focusedStyle
					continue
				}

				// Remove focused state
				i.inputs[j].Blur()
				i.inputs[j].PromptStyle = noStyle
				i.inputs[j].TextStyle = noStyle
			}

			return i, tea.Batch(cmds...)
		}
	}
	// Handle character input and blinking
	cmd := i.updateInputs(msg)

	return i, tea.Batch(cmd)
}

func (i *Input) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(i.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for j := range i.inputs {
		i.inputs[j], cmds[j] = i.inputs[j].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (i *Input) View() string {
	var b strings.Builder

	b.WriteString(BracketStyle.Render(i.title))
	b.WriteRune('\n')

	var err error
	for j := range i.inputs {
		if i.inputs[j].Err != nil {
			err = i.inputs[j].Err
		}
		b.WriteString(i.inputs[j].View())
		if j < len(i.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if err == nil {
		_, _ = fmt.Fprintf(&b, "\n\n")
	} else {
		_, _ = fmt.Fprintf(&b, "\n%s\n", errorStyle.Render(err.Error()))
	}

	button := blurredButton
	if i.focusIndex == len(i.inputs) {
		button = focusedButton
	}
	_, _ = fmt.Fprintf(&b, "%s\n\n", button)

	b.WriteString(helpStyle.Render("ESC or ctrl+C to quit"))
	b.WriteRune('\n')

	return b.String()
}
