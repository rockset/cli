package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/rockset/rockset-go-client/openapi"

	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/option"
)

func newCreateWorkspaceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "workspace",
		Aliases:     []string{"ws"},
		Short:       "create workspace",
		Long:        "create Rockset workspace",
		Args:        cobra.ExactArgs(1),
		Annotations: group("workspace"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var opts []option.WorkspaceOption

			// safe to ignore the error it is added below
			desc, _ := cmd.Flags().GetString(DescriptionFlag)
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

	cmd.Flags().StringP(DescriptionFlag, "d", "", "workspace description")
	return &cmd
}

func newDeleteWorkspaceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "workspace",
		Aliases:           []string{"ws"},
		Short:             "delete workspace",
		Long:              "delete Rockset workspace",
		Args:              cobra.ExactArgs(1),
		Annotations:       group("workspace"),
		ValidArgsFunction: workspaceCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ws := args[0]
			recurse, _ := cmd.Flags().GetBool("recurse")

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			if recurse {
				// TODO these should be moved into the go client

				collections, err := rs.ListCollections(ctx, option.WithWorkspace(ws))
				if err != nil {
					return err
				}
				for _, collection := range collections {
					if err = rs.DeleteCollection(ctx, ws, collection.GetName()); err != nil {
						return err
					}

					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "deleted collection %s\n", collection.GetName())
				}

				views, err := rs.ListViews(ctx, option.WithViewWorkspace(ws))
				if err != nil {
					return err
				}
				for _, view := range views {
					if err = rs.DeleteView(ctx, ws, view.GetName()); err != nil {
						return err
					}

					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "deleted view %s\n", view.GetName())
				}

				qls, err := rs.ListQueryLambdas(ctx, option.WithQueryLambdaWorkspace(ws))
				if err != nil {
					return err
				}
				for _, ql := range qls {
					if err = rs.DeleteQueryLambda(ctx, ws, ql.GetName()); err != nil {
						return err
					}

					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "deleted query lambda %s\n", ql.GetName())
				}

				aliases, err := rs.ListAliases(ctx, option.WithAliasWorkspace(ws))
				if err != nil {
					return err
				}
				for _, alias := range aliases {
					if err = rs.DeleteAlias(ctx, ws, alias.GetName()); err != nil {
						return err
					}

					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "deleted alias %s\n", alias.GetName())
				}

				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "waiting for all resources to be gone...\n")

				for _, collection := range collections {
					if err = rs.Wait.UntilCollectionGone(ctx, ws, collection.GetName()); err != nil {
						return err
					}
				}

				for _, view := range views {
					if err = rs.Wait.UntilViewGone(ctx, ws, view.GetName()); err != nil {
						return err
					}
				}

				for _, alias := range aliases {
					if err = rs.Wait.UntilAliasGone(ctx, ws, alias.GetName()); err != nil {
						return err
					}
				}

				// TODO wait until QLs are gone (missing support in go client)
			}

			err = rs.DeleteWorkspace(ctx, ws)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "workspace '%s' deleted\n", ws)
			return nil
		},
	}

	cmd.Flags().Bool("recurse", false, "recursively delete everything in the workspace, i.e. collections, query lambdas and views")

	return &cmd
}

func newGetWorkspaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "workspace",
		Aliases:           []string{"ws"},
		Short:             "get workspace",
		Long:              "get Rockset workspace",
		Args:              cobra.ExactArgs(1),
		Annotations:       group("workspace"),
		ValidArgsFunction: workspaceCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			ws, err := rs.GetWorkspace(ctx, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, ws)
		},
	}
}

func newListWorkspacesCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "workspaces",
		Aliases:     []string{"workspace", "ws"},
		Args:        cobra.NoArgs,
		Short:       "list workspaces",
		Long:        "list Rockset workspaces",
		Annotations: group("workspace"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			list, err := rs.ListWorkspaces(ctx)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.Workspace]{
				LessFuncs: []func(p1 *openapi.Workspace, p2 *openapi.Workspace) bool{
					sort.ByName[*openapi.Workspace],
				},
			}
			ms.Sort(list)

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}
}
