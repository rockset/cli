package cmd

import (
	"context"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
)

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

			file, _ := cmd.Flags().GetString("file")
			validate, _ := cmd.Flags().GetBool(ValidateFlag)

			// TODO handle a parameterized query
			var sql string

			// start an interactive session
			if file == "" && len(args) == 0 {
				if validate {
					return fmt.Errorf("can't validate interactive commands")
				}
				return interactive(ctx, cmd.OutOrStdout(), rs)
			}

			if file != "" && len(args) > 0 {
				return fmt.Errorf("you can only specify one of --file or a SQL query")
			}

			if file != "" {
				data, err := ioutil.ReadFile(file)
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

			result, err := rs.Query(ctx, sql)
			if err != nil {
				return err
			}

			err = showResult(cmd.OutOrStdout(), result)
			if err != nil {
				return err
			}

			return nil
		},
	}

	c.Flags().Bool(ValidateFlag, false, "validate SQL")
	c.Flags().String("file", "", "read SQL from file")
	_ = cobra.MarkFlagFilename(c.Flags(), FileFlag, ".sql")

	return &c
}

func showResult(out io.Writer, result openapi.QueryResponse) error {
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
	default:
		return fmt.Errorf("unexpected query status: %s", result.GetStatus())
	}

	return nil
}

const prompt = "rockset> "

func interactive(ctx context.Context, out io.Writer, rs *rockset.RockClient) error {
	histFile, err := historyFile()
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      prompt,
		HistoryFile: histFile,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := rl.Close(); err != nil {
			log.Printf("failed to close readline: %v", err)
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
			rl.SetPrompt(">>> ")
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
			log.Printf("failed to save history: %v", err)
		}

		result, err := rs.Query(ctx, sql)
		if err != nil {
			log.Printf("query failed: %v", err)
			continue
		}

		if err = showResult(out, result); err != nil {
			log.Printf("failed to show result: %v", err)
		}
	}

	return nil
}
