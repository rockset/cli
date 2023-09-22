package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type ProgressModel struct {
	progress progress.Model
	d        time.Duration
	untilFn  func() error
	err      error
}

const (
	padding  = 2
	maxWidth = 80
)

type tickMsg time.Time
type DoneMsg struct{}
type ErrMsg struct{ Err error }

func (m ProgressModel) Init() tea.Cmd {
	return tickCmd()
}

func NewTimeProgress(d time.Duration) *ProgressModel {
	return &ProgressModel{
		d:        d,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m *ProgressModel) Error() error {
	return m.err
}

func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case ErrMsg:
		m.err = msg.Err
		return m, tea.Quit

	case DoneMsg:
		m.progress.SetPercent(1.0)
		return m, tea.Sequence(finalPause(), tea.Quit)

	case tickMsg:
		if m.err != nil {
			return m, tea.Quit
		}

		var cmd tea.Cmd
		if m.progress.Percent() > 1.0 {
			m.d = m.d * 2
			cmd = m.progress.SetPercent(0.5)
		} else {
			cmd = m.progress.IncrPercent(1 / m.d.Seconds())
		}

		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m ProgressModel) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle.Render("Press q to stop waiting")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func finalPause() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}
