package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rockset/cli/format"
	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
)

func newListLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Short:   "list lambda",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString("workspace")

			ctx := cmd.Context()
			rs, err := rockset.NewClient(rocksetAPI(cmd))
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

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			f.FormatList(true, format.ToInterfaceArray(lambdas))

			return nil
		},
	}
	cmd.Flags().String("workspace", "", "only show query lambdas for the selected workspace")

	return cmd
}

func newExecuteLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lambda",
		Aliases: []string{"ql"},
		Short:   "execute lambda",
		Long:    "execute Rockset query lambda",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient(rocksetAPI(cmd))
			if err != nil {
				return err
			}

			fields := strings.Split(args[0], ".")
			if len(fields) != 2 {
				return fmt.Errorf("could not parse '%s' into workspace and query lambda name", args[0])
			}
			ws := fields[0]
			name := fields[1]

			version, _ := cmd.Flags().GetString("version")
			opts := []option.QueryLambdaOption{option.WithVersion(version)}

			params, _ := cmd.Flags().GetString("params")
			if params != "" {
				f, err := os.Open(params)
				if err != nil {
					return fmt.Errorf("failed to read paramater file %s: %w", params, err)
				}
				_ = f
			}

			resp, err := rs.ExecuteQueryLambda(ctx, ws, name, opts...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "ElapsedTimeMs: %d\n", *resp.Stats.ElapsedTimeMs)
			writeTable(cmd.OutOrStdout(), *resp.ColumnFields, *resp.Results)
			return nil
		},
	}

	cmd.Flags().StringP("version", "v", "latest", "query lambda version")
	cmd.Flags().StringP("params", "p", "", "query parameters file")
	return cmd
}

func writeTable(w io.Writer, columns []openapi.QueryFieldType, data []map[string]interface{}) {
	table := tablewriter.NewWriter(w)

	var headers = make([]string, len(columns))
	for i, column := range columns {
		headers[i] = column.Name
	}
	table.SetHeader(headers)

	for _, row := range data {

		var values = make([]string, len(headers))
		for i, h := range headers {
			switch row[h].(type) {
			case string:
				values[i] = row[h].(string)
			case int64:
				values[i] = strconv.FormatInt(row[h].(int64), 10)
			default:
				values[i] = fmt.Sprintf("%v", row[h])
			}
		}
		table.Append(values)
	}

	table.Render()
}
