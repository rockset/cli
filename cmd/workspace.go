package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"

	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/option"
)

func newCreateWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "create workspace",
		Long:    "create Rockset workspace",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			var opts []option.WorkspaceOption

			// safe to ignore the error it is added below
			desc, _ := cmd.Flags().GetString("description")
			if desc != "" {
				opts = append(opts, option.WithWorkspaceDescription(desc))
			}

			_, err = rs.CreateWorkspace(ctx, args[0], opts...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "workspace '%s' created\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "workspace description")
	return cmd
}

func newDeleteWorkspaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "delete workspace",
		Long:    "delete Rockset workspace",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			err = rs.DeleteWorkspace(ctx, args[0])
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "workspace '%s' deleted\n", args[0])
			return nil
		},
	}
}

func newGetWorkspaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "get workspace",
		Long:    "get Rockset workspace",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			ws, err := rs.GetWorkspace(ctx, args[0])
			if err != nil {
				return err
			}

			f, err := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			if err != nil {
				return err
			}

			f.Workspace(ws)
			return nil
		},
	}
}

func newListWorkspacesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"workspace", "ws"},
		Short:   "list workspaces",
		Long:    "list Rockset workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			list, err := rs.ListWorkspaces(ctx)
			if err != nil {
				return err
			}

			f, err := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			if err != nil {
				return err
			}

			f.Workspaces(list)
			return nil
		},
	}
}
