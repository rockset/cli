package cmd

import (
	"fmt"

	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

// https://docs.rockset.com/documentation/reference/aliases

func NewListAliasesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "aliases",
		Aliases: []string{"a", "alias"},
		Args:    cobra.NoArgs,
		Short:   "list aliases",
		Long: ` list aliases and all workspaces, or in a specific workspace
	# Documentation URLs
	https://docs.rockset.com/documentation/reference/listaliases
	https://docs.rockset.com/documentation/reference/workspacealiases`,
		Annotations: group("alias"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var opts []option.ListAliasesOption
			if ws != "" && ws != flag.AllWorkspaces {
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
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.AllWorkspaces, "only show aliases for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	return &cmd
}

func NewGetAliasCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "alias NAME",
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
		Short:   "get alias information",
		Long: `get an alias

	# Documentation URL
	https://docs.rockset.com/documentation/reference/getalias`,
		Annotations:       group("alias"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
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
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "get an alias for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	return &cmd
}

func NewCreateAliasCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "alias NAME COLLECTION [[COLLECTION] ...]",
		Aliases: []string{"a"},
		Args:    cobra.MinimumNArgs(2),
		Short:   "create alias",
		Long: `create a new alias

	# Documentation URL
	https://docs.rockset.com/documentation/reference/createalias`,
		Annotations: group("alias"),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			collections := args[1:]
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			alias, err := rs.CreateAlias(ctx, ws, name, collections)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "alias %s created\n", alias.GetName())
			return nil
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "create alias in the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	return &cmd
}

func NewUpdateAliasCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "alias NAME COLLECTION [[COLLECTION] ...]",
		Aliases: []string{"a"},
		Args:    cobra.MinimumNArgs(2),
		Short:   "update alias",
		Long: `update an alias

	# Documentation URL
	https://docs.rockset.com/documentation/reference/updatealias`,

		Annotations:       group("alias"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			collections := args[1:]
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			err = rs.UpdateAlias(ctx, ws, name, collections)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "alias %s updated\n", name)
			return nil
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "create alias in the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	return &cmd
}

func NewDeleteAliasCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "alias NAME ",
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
		Short:   "delete alias",
		Long: `delete an alias

	# Documentation URL
	https://docs.rockset.com/documentation/reference/deletealias`,
		Annotations:       group("alias"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			err = rs.DeleteAlias(ctx, ws, name)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "alias %s deleted\n", name)
			return nil
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "delete alias for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	return &cmd
}
