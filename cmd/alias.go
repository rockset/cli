package cmd

import (
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

func newListAliasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "ailases",
		Aliases:     []string{"a", "alias"},
		Args:        cobra.NoArgs,
		Short:       "list aliases",
		Annotations: group("alias"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var opts []option.ListAliasesOption
			if ws != "" && ws != AllWorkspaces {
				opts = append(opts, option.WithAliasWorkspace(ws))
			}

			aliases, err := rs.ListAliases(ctx, opts...)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.Alias]{
				LessFuncs: []func(p1 *openapi.Alias, p2 *openapi.Alias) bool{
					sort.ByWorkspace[*openapi.Alias],
					sort.ByName[*openapi.Alias],
				},
			}
			ms.Sort(aliases)

			return formatList(cmd, format.ToInterfaceArray(aliases))
		},
	}
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, AllWorkspaces, "only show aliases for the selected workspace")

	return cmd
}

func newGetAliasCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "alias NAME",
		Aliases:     []string{"a"},
		Args:        cobra.ExactArgs(1),
		Short:       "get alias information",
		Annotations: group("alias"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			alias, err := rs.GetAlias(ctx, ws, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, alias)
		},
	}
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "get an alias for the selected workspace")

	return cmd
}
