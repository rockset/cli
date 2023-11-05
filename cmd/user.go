package cmd

import (
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/openapi"
)

func newListUsersCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "users",
		Args:        cobra.NoArgs,
		Short:       "list users",
		Long:        "list Rockset users",
		Annotations: group("user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			list, err := rs.ListUsers(ctx)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.User]{
				LessFuncs: []func(p1 *openapi.User, p2 *openapi.User) bool{
					sort.ByEmail[*openapi.User],
				},
			}
			ms.Sort(list)

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")

	return &cmd
}

func newGetUserCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "user [EMAIL]",
		Short:       "get user information",
		Long:        "get Rockset user, if no email address is specified the current user is returned",
		Args:        cobra.RangeArgs(0, 1),
		Annotations: group("user"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			var user openapi.User
			if len(args) == 0 {
				user, err = rs.GetCurrentUser(ctx)
			} else {
				user, err = rs.GetUser(ctx, args[0])
			}
			if err != nil {
				return err
			}

			return formatOne(cmd, user)
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")

	return &cmd
}
