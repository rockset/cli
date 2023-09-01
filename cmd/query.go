package cmd

import (
	"context"
	"fmt"
	"github.com/rockset/cli/format"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
)

func newListQueryCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "query ID",
		Short: "list queries",
		Long:  "list queries on a virtual instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wide, _ := cmd.Flags().GetBool(WideFlag)

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			list, err := rs.GetVirtualInstanceQueries(ctx, args[0])
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			return f.FormatList(wide, format.ToInterfaceArray(list))
		},
	}

	cmd.Flags().Bool(WideFlag, false, "display more information")

	return &cmd
}

func newQueryCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "query SQL",
		Short: "execute SQL query",
		Long:  "query Rockset collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			vi, _ := cmd.Flags().GetString("vi")
			file, _ := cmd.Flags().GetString(FileFlag)
			validate, _ := cmd.Flags().GetBool(ValidateFlag)

			// TODO handle a parameterized query
			var sql string

			// start an interactive session
			if file == "" && len(args) == 0 {
				if validate {
					return fmt.Errorf("can't validate interactive commands")
				}
				return interactiveQuery(ctx, cmd.OutOrStdout(), rs)
			}

			if file != "" && len(args) > 0 {
				return fmt.Errorf("you can only specify one of --file or a SQL query")
			}

			if file != "" {
				data, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				sql = string(data)
			} else {
				sql = args[0]
			}

			if validate {
				_, err = rs.ValidateQuery(ctx, sql)
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "SQL is valid\n")
			}

			var result openapi.QueryResponse
			if vi == "" {
				result, err = rs.Query(ctx, sql)
			} else {
				result, err = rs.ExecuteQueryOnVirtualInstance(ctx, vi, sql)
			}

			if err != nil {
				return err
			}

			err = showQueryResult(cmd.OutOrStdout(), result)
			if err != nil {
				return err
			}

			return nil
		},
	}

	c.Flags().Bool(ValidateFlag, false, "validate SQL")
	c.Flags().String(FileFlag, "", "read SQL from file")
	c.Flags().String("vi", "", "execute query on virtual instance")
	_ = cobra.MarkFlagFilename(c.Flags(), FileFlag, ".sql")

	return &c
}

func showQueryResult(out io.Writer, result openapi.QueryResponse) error {
	switch result.GetStatus() {
	case "ERROR":
		var errs []string
		for _, e := range result.GetQueryErrors() {
			errs = append(errs, e.GetMessage())
		}
		_, _ = fmt.Fprintf(out, "query %s failed:\n%s\n", result.GetQueryId(), strings.Join(errs, "\n"))
	case "QUEUED", "RUNNING":
		_, _ = fmt.Fprintf(out, "your query %s is %s\n", result.GetQueryId(), result.GetStatus())
	case "COMPLETED":
		t := tablewriter.NewWriter(out)

		if len(result.GetColumnFields()) == 0 {
			// in a "SELECT *" query the ColumnFields isn't populated
			// what order should the columns be presented in?

			var headers []string
			for k := range result.Results[0] {
				headers = append(headers, k)
			}
			t.SetHeader(headers)

			for _, row := range result.Results {
				var r []string
				for _, h := range headers {
					r = append(r, fmt.Sprintf("%v", row[h]))
				}
				t.Append(r)
			}
		} else {
			var headers []string
			for _, h := range result.GetColumnFields() {
				headers = append(headers, h.Name)
			}
			t.SetHeader(headers)

			for _, row := range result.Results {
				var r []string
				for _, column := range result.GetColumnFields() {
					r = append(r, fmt.Sprintf("%v", row[column.GetName()]))
				}
				t.Append(r)
			}
		}

		t.Render()
		_, _ = fmt.Fprintf(out, "Elapsed time: %d ms\n\n", result.Stats.GetElapsedTimeMs())
	default:
		return fmt.Errorf("unexpected query status: %s", result.GetStatus())
	}

	return nil
}

var (
	prompt             = R + "> "
	continuationPrompt = color.CyanString(">") + color.MagentaString(">") + "> "
)

func interactiveQuery(ctx context.Context, out io.Writer, rs *rockset.RockClient) error {
	histFile, err := historyFile()
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(out, "%s interactive console. End your SQL with ;\n", Rockset)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      prompt,
		HistoryFile: histFile,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := rl.Close(); err != nil {
			slog.Error("failed to close readline", err)
		}
	}()

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt(continuationPrompt)
			continue
		}

		sql := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt(prompt)

		if strings.HasPrefix(sql, "help") {
			_, _ = fmt.Fprintf(out, "help command TBW\n")
			continue
		}

		if err = rl.SaveHistory(sql); err != nil {
			slog.Error("failed to save history", err)
		}

		result, err := rs.Query(ctx, sql)
		if err != nil {
			slog.Error("query failed", err)
			continue
		}

		if err = showQueryResult(out, result); err != nil {
			slog.Error("failed to show result", err)
		}
	}

	return nil
}