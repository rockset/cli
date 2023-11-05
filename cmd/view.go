package cmd

import (
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

func newListViewsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "views",
		Aliases:     []string{"view", "v"},
		Args:        cobra.NoArgs,
		Short:       "list views",
		Annotations: group("view"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			var opts []option.ListViewOption
			if ws != "" && ws != flag.AllWorkspaces {
				opts = append(opts, option.WithViewWorkspace(ws))
			}

			views, err := rs.ListViews(ctx, opts...)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.View]{
				LessFuncs: []func(p1 *openapi.View, p2 *openapi.View) bool{
					sort.ByName[*openapi.View],
				},
			}
			ms.Sort(views)

			return formatList(cmd, format.ToInterfaceArray(views))
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.AllWorkspaces, "only show views for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	return &cmd
}

func newGetViewCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "view NAME",
		Aliases:           []string{"v"},
		Args:              cobra.ExactArgs(1),
		Short:             "get view information",
		Annotations:       group("view"),
		ValidArgsFunction: completion.View,
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			view, err := rs.GetView(ctx, ws, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, view)
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.Description, "only show views for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	return &cmd
}
