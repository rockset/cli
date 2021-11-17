package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
)

func newListUsersCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "users",
		Short: "list users",
		Long:  "list Rockset users",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			list, err := rs.ListUsers(ctx)
			if err != nil {
				return err
			}

			f, err := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			if err != nil {
				return err
			}

			f.Users(list)
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

			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			u, err := rs.GetCurrentUser(ctx)
			if err != nil {
				return err
			}

			f, err := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			if err != nil {
				return err
			}

			f.User(u)
			return nil
		},
	}

	c.Flags().Bool("wide", false, "display more information")

	return &c
}
