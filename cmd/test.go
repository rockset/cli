package cmd

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/tui"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func newTestCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:    "test",
		Short:  "test",
		Hidden: true,
	}

	cmd.AddCommand(newTestProgressCmd())
	cmd.AddCommand(newTestInputCmd())
	cmd.AddCommand(newTestSelectorCmd())

	return &cmd
}

func newTestInputCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "input",
		Short:       "test input",
		Long:        "used for testing input fields",
		Args:        cobra.NoArgs,
		Annotations: group("test"),
		RunE: func(cmd *cobra.Command, args []string) error {
			model := tui.NewInput("Enter context options", []tui.InputConfig{
				{Placeholder: "name", Prompt: "Name: "},
				{Placeholder: "server", Prompt: "Cluster: ", Validate: func(s string) error {
					for _, c := range config.Clusters {
						if strings.HasPrefix(c, s) {
							return nil
						}
					}
					return fmt.Errorf("cluster must be one of: %s", strings.Join(config.Clusters, ", "))
				}},
				{Placeholder: "organization", Prompt: "Organization: "},
			})
			p := tea.NewProgram(model)

			if _, err := p.Run(); err != nil {
				return err
			}

			if model.Err != nil {
				return model.Err
			}

			for _, f := range model.Fields {
				fmt.Printf("%s\n", f)
			}

			return nil
		},
	}

	return &cmd
}

func newTestProgressCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "progress",
		Short:       "test progress",
		Long:        "used for testing progress bar",
		Args:        cobra.NoArgs,
		Annotations: group("test"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fail, _ := cmd.Flags().GetBool("fail")
			done, _ := cmd.Flags().GetDuration("done")
			estimate, _ := cmd.Flags().GetDuration("estimate")

			model := tui.NewTimeProgress(estimate)
			p := tea.NewProgram(model)
			ctx := cmd.Context()

			go func() {
				select {
				case <-ctx.Done():
					if err := ctx.Err(); err != nil {
						p.Send(tui.ErrMsg{Err: err})
					}
				case <-time.After(done):
					if fail {
						p.Send(tui.ErrMsg{Err: errors.New("boom")})
					} else {
						p.Send(tui.DoneMsg{})
					}
				}
			}()

			if _, err := p.Run(); err != nil {
				return err
			}

			if err := model.Error(); err == nil {
				fmt.Printf("done\n")
			} else {
				fmt.Printf("err: %v\n", err)
			}

			return nil
		},
	}

	cmd.Flags().Duration("estimate", 10*time.Second, "estimate how long it takes")
	cmd.Flags().Duration("done", 5*time.Second, "done after")
	cmd.Flags().Bool("fail", false, "should it fail")

	return &cmd
}

func newTestSelectorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "selector",
		Short:       "test selector",
		Long:        "used for testing the selector",
		Args:        cobra.NoArgs,
		Annotations: group("test"),
		RunE: func(cmd *cobra.Command, args []string) error {
			values := []string{"a", "b", "c"}
			model := tui.NewSelector(values, 2)
			p := tea.NewProgram(model)

			if _, err := p.Run(); err != nil {
				return err
			}

			if model.Selected > 0 {
				fmt.Printf("selection: %s\n", values[model.Selected])
			} else {
				fmt.Printf("no selection\n")
			}

			return nil
		},
	}

	return &cmd
}
