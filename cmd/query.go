package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"log/slog"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/rockset/rockset-go-client/paginate"

	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/lookup"
	"github.com/rockset/cli/tui"
)

func newListQueriesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "queries [ID|NAME]",
		Aliases:     []string{"query", "q"},
		Short:       "list queries",
		Long:        "list all actively queued and executing queries, or on a specific virtual instance",
		Args:        cobra.RangeArgs(0, 1),
		Annotations: group("query"), // TODO should this be in the VI group too?
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var list []openapi.QueryInfo
			if len(args) == 0 {
				list, err = rs.ListActiveQueries(ctx)
			} else {
				id, err := lookup.VirtualInstanceNameOrIDtoID(ctx, rs, args[0])
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

func newGetQueryInfoCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "info ID",
		Aliases:     []string{"i"},
		Short:       "get query info",
		Long:        "get information about a query",
		Args:        cobra.ExactArgs(1),
		Annotations: group("query"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			q, err := rs.GetQueryInfo(ctx, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, q)
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")

	return &cmd
}

func newGetQueryResultCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "results ID",
		Aliases: []string{"r"},
		Short:   "get query results",
		Long: fmt.Sprintf("Get query results for a previously executed query. "+
			"If --%s isn't specified, all documents are retrieved.", flag.Docs),
		Args:        cobra.ExactArgs(1),
		Annotations: group("query"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			// have to get info to get thew query stats
			info, err := rs.GetQueryInfo(ctx, args[0])
			if err != nil {
				return err
			}

			var options []option.QueryResultOption
			if cursor, _ := cmd.Flags().GetString(flag.Cursor); cursor != "" {
				options = append(options, option.WithQueryResultCursor(cursor))
			}
			docs, _ := cmd.Flags().GetInt32(flag.Docs)
			if docs != 0 {
				options = append(options, option.WithQueryResultDocs(docs))
			}
			if offset, _ := cmd.Flags().GetInt32(flag.Offset); offset != 0 {
				options = append(options, option.WithQueryResultOffset(offset))
			}

			stats := info.GetStats()
			var list []map[string]interface{}
			var cursor string
			// if we don't have any options, we fetch all documents
			if len(options) == 0 {
				p := paginate.New(rs)
				docCh := make(chan map[string]any, 100)

				go func() {
					err = p.GetQueryResults(ctx, docCh, args[0])
				}()

				for doc := range docCh {
					list = append(list, doc)
				}
			} else {
				result, err := rs.GetQueryResults(ctx, args[0], options...)
				if err != nil {
					return err
				}
				page := result.GetPagination()
				cursor = page.GetNextCursor()
				list = result.Results
			}

			return showQueryPaginationResponse(cmd.OutOrStdout(), cursor, stats.GetElapsedTimeMs(), list)
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")
	cmd.Flags().String(flag.Cursor, "", "cursor to current page, defaults to first page")
	cmd.Flags().Int32(flag.Docs, 0, "number of documents to fetch, 0 means fetch max")
	cmd.Flags().Int32(flag.Offset, 0, "offset from the cursor of the first document to be returned")

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
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			async, _ := cmd.Flags().GetBool(flag.Async)
			vi, _ := cmd.Flags().GetString(flag.VI)
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
				return nil
			}

			var options []option.QueryOption
			if async {
				// TODO inform the user that --validate and --async are mutually exclusive?
				options = append(options, option.WithAsync())
			}

			var result openapi.QueryResponse
			if vi == "" {
				result, err = rs.Query(ctx, sql, options...)
			} else {
				result, err = rs.ExecuteQueryOnVirtualInstance(ctx, vi, sql, options...)
			}
			if err != nil {
				return err
			}

			if async {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "query ID is: %s\n", result.GetQueryId())
				return nil
			}

			err = showQueryResponse(cmd.OutOrStdout(), result)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().Bool(flag.Async, false, "execute the query asynchronously")
	cmd.Flags().Bool(flag.Validate, false, "validate SQL")
	cmd.Flags().String(flag.File, "", "read SQL from file")
	cmd.Flags().String(flag.VI, "", "execute query on virtual instance")
	_ = cobra.MarkFlagFilename(cmd.Flags(), flag.File, ".sql")

	return &cmd
}

func showQueryPaginationResponse(out io.Writer, cursor string, elapsedMs int64, results []map[string]interface{}) error {
	if len(results) == 0 {
		return errors.New("query returned no rows")
	}

	var headers []string
	for h := range results[0] {
		headers = append(headers, h)
	}

	showQueryResponseTable(out, cursor, elapsedMs, headers, results)

	return nil
}

func showQueryResponse(out io.Writer, result openapi.QueryResponse) error {
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
		var headers []string
		if len(result.GetColumnFields()) == 0 {
			// in a "SELECT *" query the ColumnFields isn't populated what order should the columns be presented in?
			for h := range result.Results[0] {
				headers = append(headers, h)
			}
		} else {
			for _, h := range result.GetColumnFields() {
				headers = append(headers, h.Name)
			}
		}
		stats := result.GetStats()
		showQueryResponseTable(out, "", stats.GetElapsedTimeMs(), headers, result.Results)
	default:
		return fmt.Errorf("unexpected query status: %s", result.GetStatus())
	}

	return nil
}

func showQueryResponseTable(out io.Writer, cursor string, elapsedMs int64, headers []string, results []map[string]interface{}) {
	t := tui.NewTable(out)
	t.Headers(headers...)

	for _, row := range results {
		var r []string
		for _, h := range headers {
			r = append(r, fmt.Sprintf("%v", row[h]))
		}
		t.Row(r...)
	}

	_, _ = fmt.Fprintln(out, t.Render())
	if cursor != "" {
		_, _ = fmt.Fprintf(out, "Next cursor: %s\n", cursor)
	}
	_, _ = fmt.Fprintf(out, "Elapsed time: %d ms\n\n", elapsedMs)
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
			slog.Error("failed to close readline", "err", err)
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
			slog.Error("failed to save history", "err", err)
		}

		executeQuery(ctx, out, rs, sql)
	}

	return nil
}

func executeQuery(ctx context.Context, out io.Writer, rs *rockset.RockClient, sql string) {
	result, err := rs.Query(ctx, sql)
	if err != nil {
		// TODO should this use tui.ShowError()?
		_, _ = fmt.Fprintf(out, "%s\n", tui.ErrorStyle.Render("query failed:", err.Error()))
		return
	}

	if err = showQueryResponse(out, result); err != nil {
		slog.Error("failed to show result", "err", err)
	}
}
