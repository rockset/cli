package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client"
)

func newListUsersCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "users",
		Short: "list users",
		Long:  "list Rockset users",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockset.NewClient(rockOption(cmd))
			if err != nil {
				return err
			}

			list, err := rs.ListUsers(ctx)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			f.FormatList(true, format.ToInterfaceArray(list))
			return nil
		},
	}

	c.Flags().Bool("wide", false, "display more information")

	return &c
}

func newGetUserCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "user",
		Short: "get user",
		Long:  "get current Rockset user",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockset.NewClient(rockOption(cmd))
			if err != nil {
				return err
			}

			u, err := rs.GetCurrentUser(ctx)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			f.Format(true, u)
			return nil
		},
	}

	c.Flags().Bool("wide", false, "display more information")

	return &c
}
