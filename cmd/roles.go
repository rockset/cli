package cmd

import (
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

func newListRolesCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:         "roles",
		Aliases:     []string{"r", "role"},
		Args:        cobra.NoArgs,
		Short:       "list roles",
		Annotations: group("role"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			roles, err := rs.ListRoles(ctx)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.Role]{
				LessFuncs: []func(p1 *openapi.Role, p2 *openapi.Role) bool{
					sort.ByRoleName[*openapi.Role],
				},
			}
			ms.Sort(roles)

			return formatList(cmd, format.ToInterfaceArray(roles))
		},
	}

	return &cmd
}

func newGetRoleCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:               "role NAME",
		Aliases:           []string{"r"},
		Args:              cobra.ExactArgs(1),
		Short:             "get role information",
		Annotations:       group("role"),
		ValidArgsFunction: completion.Role,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			alias, err := rs.GetRole(ctx, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, alias)
		},
	}

	return &cmd
}
