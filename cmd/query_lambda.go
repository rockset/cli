package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

func newListQueryLambdasCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "lambdas",
		Aliases:     []string{"ql", "qls", "querylambdas"},
		Args:        cobra.NoArgs,
		Annotations: group("lambda"),
		Short:       "list query lambdas",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}
			var opts []option.ListQueryLambdaOption
			if ws != "" && ws != flag.AllWorkspaces {
				opts = append(opts, option.WithQueryLambdaWorkspace(ws))
			}

			lambdas, err := rs.ListQueryLambdas(ctx, opts...)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.QueryLambda]{
				LessFuncs: []func(p1 *openapi.QueryLambda, p2 *openapi.QueryLambda) bool{
					sort.ByWorkspace[*openapi.QueryLambda],
					sort.ByName[*openapi.QueryLambda],
				},
			}
			ms.Sort(lambdas)

			return formatList(cmd, format.ToInterfaceArray(lambdas))
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.AllWorkspaces, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	return &cmd
}

func newGetQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Short:   "get query lambda",
		Long: `get query lambda information, has options to get a specific tag or version,
or to get all tags or versions`,
		Args:              cobra.ExactArgs(1),
		Annotations:       group("lambda"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			tag, _ := cmd.Flags().GetString(flag.Tag)
			tags, _ := cmd.Flags().GetBool(flag.Tags)
			version, _ := cmd.Flags().GetString(flag.Version)
			versions, _ := cmd.Flags().GetBool(flag.Versions)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			if tag != "" {
				ql, err := rs.GetQueryLambdaVersionByTag(ctx, ws, args[0], tag)
				if err != nil {
					return err
				}

				return formatOne(cmd, ql)
			}

			if version != "" {
				ql, err := rs.GetQueryLambdaVersion(ctx, ws, args[0], version)
				if err != nil {
					return err
				}

				return formatOne(cmd, ql)
			}

			if tags {
				list, err := rs.ListQueryLambdaTags(ctx, ws, args[0])
				if err != nil {
					return err
				}

				return formatList(cmd, format.ToInterfaceArray(list))
			}

			if versions {
				list, err := rs.ListQueryLambdaVersions(ctx, ws, args[0])
				if err != nil {
					return err
				}

				return formatList(cmd, format.ToInterfaceArray(list))
			}

			ql, err := rs.GetQueryLambdaVersionByTag(ctx, ws, args[0], "latest")
			if err != nil {
				return err
			}

			return formatOne(cmd, ql)
		},
	}

	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().Bool(flag.Versions, false, "show all versions of this query lambda")
	cmd.Flags().Bool(flag.Tags, false, "show all tags for this query lambda")

	cmd.Flags().String(flag.Tag, "", "only show this query lambda tag")
	_ = cmd.RegisterFlagCompletionFunc("tag", completion.LambdaTag(Version))

	cmd.Flags().String(flag.Version, "", "only show this query lambda version")
	_ = cmd.RegisterFlagCompletionFunc(flag.Version, completion.LambdaVersion(Version))

	cmd.MarkFlagsMutuallyExclusive(flag.Version, "versions", "tag", "tags")

	return &cmd
}

func NewExecuteQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "lambda NAME",
		Aliases:           []string{"ql"},
		Short:             "execute lambda",
		Long:              "execute Rockset query lambda",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.Alias(Version),
		Annotations:       group("lambda"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var opts []option.QueryLambdaOption
			if version, _ := cmd.Flags().GetString(flag.Version); version != "" {
				opts = []option.QueryLambdaOption{option.WithVersion(version)}
			}
			if tag, _ := cmd.Flags().GetString("tag"); tag != "" {
				opts = []option.QueryLambdaOption{option.WithTag(tag)}
			}

			params, _ := cmd.Flags().GetStringArray("param")
			paramFile, _ := cmd.Flags().GetString("params-file")
			if len(params) > 0 && paramFile != "" {

			}

			if paramFile != "" {
				f, err := os.Open(paramFile)
				if err != nil {
					return fmt.Errorf("failed to read paramater file %s: %w", paramFile, err)
				}
				_ = f
				panic("not implemented - need to define file format")
			} else {
				for _, p := range params {
					fields := strings.SplitN(p, ":", 2)
					opts = append(opts, option.WithQueryLambdaParameter(fields[0], "", fields[1]))
				}
			}

			resp, err := rs.ExecuteQueryLambda(ctx, ws, args[0], opts...)
			if err != nil {
				return err
			}

			return showQueryResult(cmd.OutOrStdout(), resp)
		},
	}

	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace name")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().String(flag.Version, "", "query lambda version")
	cmd.Flags().String("tag", "", "query lambda tag")
	cmd.Flags().StringP("params-file", "P", "", "query parameters file")
	cmd.Flags().StringArrayP("param", "p", nil, "query parameters")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "params", ".json")

	return &cmd
}

func newCreateQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "lambda",
		Aliases:     []string{"ql"},
		Args:        cobra.ExactArgs(1),
		Short:       "create query lambda",
		Annotations: group("lambda"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			sql, _ := cmd.Flags().GetString(flag.SQL)
			description, _ := cmd.Flags().GetString(flag.Description)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var options []option.CreateQueryLambdaOption
			if description != "" {
				options = append(options, option.WithQueryLambdaDescription(description))
			}

			ql, err := rs.CreateQueryLambda(ctx, ws, args[0], sql, options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "created query lambda %s in %s", ql.GetName(), ql.GetWorkspace())

			return nil
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().String(flag.Description, "", "description of the query lambda")

	cmd.Flags().String(flag.SQL, "", "file containing SQL")
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.SQL)
	_ = cobra.MarkFlagFilename(cmd.Flags(), flag.SQL, ".sql")

	return &cmd
}

func newDeleteQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "lambda",
		Aliases:           []string{"ql"},
		Args:              cobra.ExactArgs(1),
		Short:             "delete query lambda",
		Annotations:       group("lambda"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			err = rs.DeleteQueryLambda(ctx, ws, args[0])
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "deleted query lambda %s in %s", args[0], ws)

			return nil
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().String(flag.Description, "", "description of the query lambda")

	return &cmd
}

func newUpdateQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "lambda",
		Aliases:           []string{"ql"},
		Args:              cobra.ExactArgs(1),
		Short:             "update query lambda",
		Annotations:       group("lambda"),
		ValidArgsFunction: completion.Lambda(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			sql, _ := cmd.Flags().GetString(flag.SQL)
			description, _ := cmd.Flags().GetString(flag.Description)

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var options []option.CreateQueryLambdaOption
			if description != "" {
				options = append(options, option.WithQueryLambdaDescription(description))
			}

			ql, err := rs.UpdateQueryLambda(ctx, ws, args[0], sql, options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "updated query lambda %s in %s", ql.GetName(), ql.GetWorkspace())

			return nil
		},
	}
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().String(flag.Description, "", "description of the query lambda")

	cmd.Flags().String(flag.SQL, "", "file containing SQL")
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.SQL)
	_ = cobra.MarkFlagFilename(cmd.Flags(), flag.SQL, ".sql")

	return &cmd
}
