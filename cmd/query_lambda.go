package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

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
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}
			var opts []option.ListQueryLambdaOption
			if ws != "" && ws != AllWorkspaces {
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
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, AllWorkspaces, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

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
		ValidArgsFunction: lambdaCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			tag, _ := cmd.Flags().GetString(TagFlag)
			tags, _ := cmd.Flags().GetBool(TagsFlag)
			version, _ := cmd.Flags().GetString(VersionFlag)
			versions, _ := cmd.Flags().GetBool(VersionsFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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

	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().Bool(VersionsFlag, false, "show all versions of this query lambda")
	cmd.Flags().Bool(TagsFlag, false, "show all tags for this query lambda")

	cmd.Flags().String(TagFlag, "", "only show this query lambda tag")
	_ = cmd.RegisterFlagCompletionFunc("tag", lambdaTagsCompletion)

	cmd.Flags().String(VersionFlag, "", "only show this query lambda version")
	_ = cmd.RegisterFlagCompletionFunc(VersionFlag, lambdaVersionsCompletion)

	cmd.MarkFlagsMutuallyExclusive(VersionFlag, "versions", "tag", "tags")

	return &cmd
}

func newExecuteQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "lambda NAME",
		Aliases:           []string{"ql"},
		Short:             "execute lambda",
		Long:              "execute Rockset query lambda",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: lambdaCompletion,
		Annotations:       group("lambda"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var opts []option.QueryLambdaOption
			if version, _ := cmd.Flags().GetString(VersionFlag); version != "" {
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
					fields := strings.SplitN(p, ":", 1)
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

	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace name")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().String(VersionFlag, "", "query lambda version")
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
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			sql, _ := cmd.Flags().GetString(SQLFlag)
			description, _ := cmd.Flags().GetString(DescriptionFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().String(DescriptionFlag, "", "description of the query lambda")

	cmd.Flags().String(SQLFlag, "", "file containing SQL")
	_ = cobra.MarkFlagRequired(cmd.Flags(), SQLFlag)
	_ = cobra.MarkFlagFilename(cmd.Flags(), SQLFlag, ".sql")

	return &cmd
}

func newDeleteQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "lambda",
		Aliases:           []string{"ql"},
		Args:              cobra.ExactArgs(1),
		Short:             "delete query lambda",
		Annotations:       group("lambda"),
		ValidArgsFunction: lambdaCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().String(DescriptionFlag, "", "description of the query lambda")

	return &cmd
}

func newUpdateQueryLambdaCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "lambda",
		Aliases:           []string{"ql"},
		Args:              cobra.ExactArgs(1),
		Short:             "update query lambda",
		Annotations:       group("lambda"),
		ValidArgsFunction: lambdaCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			sql, _ := cmd.Flags().GetString(SQLFlag)
			description, _ := cmd.Flags().GetString(DescriptionFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().String(DescriptionFlag, "", "description of the query lambda")

	cmd.Flags().String(SQLFlag, "", "file containing SQL")
	_ = cobra.MarkFlagRequired(cmd.Flags(), SQLFlag)
	_ = cobra.MarkFlagFilename(cmd.Flags(), SQLFlag, ".sql")

	return &cmd
}
