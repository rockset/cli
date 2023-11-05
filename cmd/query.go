package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"log/slog"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"

	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/tui"
)

func newListQueryCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "query [virtual instance ID|NAME]",
		Aliases:     []string{"queries", "q"},
		Short:       "list queries",
		Long:        "list all active queries, or on a specific virtual instance",
		Args:        cobra.RangeArgs(0, 1),
		Annotations: group("query"), // TODO should this be in the VI group too?
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			var list []openapi.QueryInfo
			if len(args) == 0 {
				list, err = rs.ListActiveQueries(ctx)
			} else {
				id, err := viNameOrIDtoID(ctx, rs, args[0])
				if err != nil {
					return err
				}

				list, err = rs.ListVirtualInstanceQueries(ctx, id)
			}
			if err != nil {
				return err
			}

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")

	return &cmd
}

func newQueryCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "query SQL",
		Short:       "execute SQL query",
		Long:        "query Rockset collections",
		Annotations: group("query"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			vi, _ := cmd.Flags().GetString("vi")
			file, _ := cmd.Flags().GetString(flag.File)
			validate, _ := cmd.Flags().GetBool(flag.Validate)

			// TODO handle a parameterized query
			var sql string

			// start an interactive session
			if file == "" && len(args) == 0 {
				if validate {
					return fmt.Errorf("can't validate interactive commands")
				}
				return interactiveQuery(ctx, io.NopCloser(cmd.InOrStdin()), cmd.OutOrStdout(), rs)
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

	cmd.Flags().Bool(flag.Validate, false, "validate SQL")
	cmd.Flags().String(flag.File, "", "read SQL from file")
	cmd.Flags().String("vi", "", "execute query on virtual instance")
	_ = cobra.MarkFlagFilename(cmd.Flags(), flag.File, ".sql")

	return &cmd
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
		t := tui.NewTable(out)

		if len(result.GetColumnFields()) == 0 {
			// pick out the headers from the first doc
			var headers []string
			for k := range result.Results[0] {
				headers = append(headers, k)
			}
			t.Headers(headers...)

			for _, row := range result.Results {
				var r []string
				for _, h := range headers {
					r = append(r, fmt.Sprintf("%v", row[h]))
				}
				t.Row(r...)
			}
		} else {
			var headers []string
			for _, h := range result.GetColumnFields() {
				headers = append(headers, h.Name)
			}
			t.Headers(headers...)

			for _, row := range result.Results {
				var r []string
				for _, column := range result.GetColumnFields() {
					r = append(r, fmt.Sprintf("%v", row[column.GetName()]))
				}
				t.Row(r...)
			}
		}

		_, _ = fmt.Fprintln(out, t.Render())
		_, _ = fmt.Fprintf(out, "Elapsed time: %d ms\n\n", result.Stats.GetElapsedTimeMs())
	default:
		return fmt.Errorf("unexpected query status: %s", result.GetStatus())
	}

	return nil
}

func interactiveQuery(ctx context.Context, in io.ReadCloser, out io.Writer, rs *rockset.RockClient) error {
	histFile, err := config.HistoryFile()
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(out, "%s interactive console. End your SQL with ;\n", tui.Rockset)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      tui.Prompt,
		Stdin:       in,
		Stdout:      out,
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
			if len(cmds) > 0 {
				// in case someone is sending the SQL over a pipe and there isn't a ";" at the end
				executeQuery(ctx, out, rs, strings.Join(cmds, " "))
			}
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt(tui.ContinuationPrompt)
			continue
		}

		sql := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt(tui.Prompt)

		if strings.HasPrefix(sql, "help") {
			_, _ = fmt.Fprintf(out, "help command TBW\n")
			continue
		}

		if err = rl.SaveHistory(sql); err != nil {
			slog.Error("failed to save history", err)
		}

		executeQuery(ctx, out, rs, sql)
	}

	return nil
}

func executeQuery(ctx context.Context, out io.Writer, rs *rockset.RockClient, sql string) {
	result, err := rs.Query(ctx, sql)
	if err != nil {
		slog.Error("query failed", err)
		return
	}

	if err = showQueryResult(out, result); err != nil {
		slog.Error("failed to show result", err)
	}
}
