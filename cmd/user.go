package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/openapi"
)

func newListUsersCmd() *cobra.Command {
	c := cobra.Command{
		Use:         "users",
		Args:        cobra.NoArgs,
		Short:       "list users",
		Long:        "list Rockset users",
		Annotations: group("user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			list, err := rs.ListUsers(ctx)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), FormatFromCommand(cmd), true)
			return f.FormatList(true, format.ToInterfaceArray(list))
		},
	}

	c.Flags().Bool(WideFlag, false, "display more information")

	return &c
}

func newGetUserCmd() *cobra.Command {
	c := cobra.Command{
		Use:         "user [EMAIL]",
		Short:       "get user information",
		Long:        "get Rockset user, if no email address is specified the current user is returned",
		Args:        cobra.RangeArgs(0, 1),
		Annotations: group("user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var u openapi.User
			if len(args) == 0 {
				u, err = rs.GetCurrentUser(ctx)
			} else {
				u, err = rs.GetUser(ctx, args[0])
			}
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), FormatFromCommand(cmd), true)

			return f.Format(true, u)
		},
	}

	c.Flags().Bool(WideFlag, false, "display more information")

	return &c
}
