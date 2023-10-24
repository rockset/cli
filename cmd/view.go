package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/option"
)

func newListViewsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "views",
		Aliases:     []string{"view", "v"},
		Args:        cobra.NoArgs,
		Short:       "list views",
		Annotations: group("view"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var opts []option.ListViewOption
			if ws != "" && ws != AllWorkspaces {
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
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, AllWorkspaces, "only show views for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	return &cmd
}

func newGetViewCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "view NAME",
		Aliases:           []string{"v"},
		Args:              cobra.ExactArgs(1),
		Short:             "get view information",
		Annotations:       group("view"),
		ValidArgsFunction: viewCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show views for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	return &cmd
}
