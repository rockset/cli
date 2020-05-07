package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/olekukonko/tablewriter"
	"github.com/rockset/rockset-go-client"
	models "github.com/rockset/rockset-go-client/lib/go"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
)

func newExecuteLambdaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lambda",
		Aliases: []string{"ql"},
		Short: "execute lambda",
		Long:  "execute Rockset query lambda",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := rockset.NewClient(rockset.FromEnv())
			if err != nil {
				return err
			}

			fields := strings.Split(args[0], ".")
			if len(fields) != 2 {
				return fmt.Errorf("could not parse '%s' into workspace and query lambda name", args[0])
			}

			version, _ := cmd.Flags().GetInt32("version")
			if version == 0 {
				if version, err = getLatestVersion(rs, fields[0], fields[1]); err != nil {
					return err
				}
			}

			execOpts := &models.ExecuteOpts{}

			params, _ := cmd.Flags().GetString("params")
			if params != "" {
				f, err := os.Open(params)
				if err != nil {
					return fmt.Errorf("failed to read paramater file %s: %w", params, err)
				}

				var qp []models.QueryParameter
				d := json.NewDecoder(f)
				err = d.Decode(&qp)
				if err != nil {
					return fmt.Errorf("failed to unmarshal %s: %w", params, err)
				}

				execOpts.Body = optional.NewInterface(models.ExecuteQueryLambdaRequest{
					Parameters: qp,
				})
			}

			res, _, err := rs.QueryLambdas.Execute(fields[0], fields[1], version, execOpts)
			if err != nil {
				if err, ok := rockset.AsRocksetError(err); ok {
					return fmt.Errorf("%s", err.Message)
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "stats: %+v\n", *res.Stats)
			writeTable(cmd.OutOrStdout(), res.ColumnFields, res.Results)
			return nil
		},
	}

	cmd.Flags().Int32P("version", "v", 0, "query lambda version")
	cmd.Flags().StringP("params", "p", "", "query parameters file")
	return cmd
}

func getLatestVersion(rs *rockset.RockClient, workspace, queryLambda string) (int32, error) {
	rsp, _, err := rs.QueryLambdas.List_2(workspace, queryLambda)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest version of %s.%s: %w", workspace, queryLambda, err)
	}
	var maxVersion int32
	for _, ql := range rsp.Data {
		if ql.Version > maxVersion {
			maxVersion = ql.Version
		}
	}
	return maxVersion, nil
}

func writeTable(w io.Writer, columns []models.QueryFieldType, data []interface{}) {
	table := tablewriter.NewWriter(w)

	var headers = make([]string, len(columns))
	for i, column := range columns {
		headers[i] = column.Name
	}
	table.SetHeader(headers)

	for _, row := range data {
		m := row.(map[string]interface{})

		var values = make([]string, len(headers))
		for i, h := range headers {
			switch m[h].(type) {
			case string:
				values[i] = m[h].(string)
			case int64:
				values[i] = strconv.FormatInt(m[h].(int64), 10)
			default:
				values[i] = fmt.Sprintf("%v", m[h])
			}
		}
		table.Append(values)
	}

	table.Render()
}
