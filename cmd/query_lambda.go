package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/rockset/rockset-go-client/option"
)

func newListQueryLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Args:    cobra.NoArgs,
		Short:   "list lambda",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}
			var opts []option.ListQueryLambdaOption
			if ws != "" {
				opts = append(opts, option.WithQueryLambdaWorkspace(ws))
			}

			lambdas, err := rs.ListQueryLambdas(ctx, opts...)
			if err != nil {
				return nil
			}

			return formatList(cmd, format.ToInterfaceArray(lambdas))
		},
	}
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")

	return cmd
}

func newGetQueryLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Args:    cobra.ExactArgs(1),
		Short:   "get query lambda",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			tag, _ := cmd.Flags().GetString("tag")

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			ql, err := rs.GetQueryLambdaVersionByTag(ctx, ws, args[0], tag)
			if err != nil {
				return err
			}

			return formatOne(cmd, ql)
		},
	}
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")
	cmd.Flags().String("tag", "latest", "query lambda tag")

	return cmd
}

func newExecuteQueryLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Short:   "execute lambda",
		Long:    "execute Rockset query lambda",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var opts []option.QueryLambdaOption
			if version, _ := cmd.Flags().GetString("version"); version != "" {
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

			showQueryResult(cmd.OutOrStdout(), resp)
			return nil
		},
	}

	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace name")
	cmd.Flags().String("version", "", "query lambda version")
	cmd.Flags().String("tag", "", "query lambda tag")
	cmd.Flags().StringP("params-file", "P", "", "query parameters file")
	cmd.Flags().StringArrayP("param", "p", nil, "query parameters")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "params", ".json")

	return cmd
}

func newCreateQueryLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Args:    cobra.ExactArgs(1),
		Short:   "create query lambda",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			sql, _ := cmd.Flags().GetString("sql")

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}
			var options []option.CreateQueryLambdaOption

			ql, err := rs.CreateQueryLambda(ctx, ws, args[0], sql, options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "created query lambda %s in %s", ql.GetName(), ql.GetWorkspace())

			return nil
		},
	}
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "only show query lambdas for the selected workspace")
	cmd.Flags().String("sql", "", "file containing SQL")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "sql")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "sql", ".sql")

	return cmd
}