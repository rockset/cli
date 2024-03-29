package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
)

func newDeleteDocumentsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "documents [id] [id] ...",
		Aliases: []string{"doc", "docs"},
		Short:   "delete documents",
		Long:    "delete documents from a collection",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			coll, _ := cmd.Flags().GetString("collection")

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			// TODO make it possible to read document IDs from stdin
			res, err := rs.DeleteDocuments(ctx, ws, coll, args)
			if err != nil {
				return err
			}

			var count, failed int
			for _, d := range res {
				if d.GetStatus() != "DELETED" {
					failed++
					fmt.Fprintf(cmd.OutOrStdout(), "failed to delete document %s\n", d.GetId())
					continue
				}
				count++
			}

			fmt.Fprintf(cmd.OutOrStdout(), "deleted %d documents\n", count)
			if failed > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "failed to delete %d documents\n", failed)
			}

			return nil
		},
	}

	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace name")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().String(flag.Collection, "", "collection name")
	_ = cmd.MarkFlagRequired(flag.Collection)
	_ = cmd.RegisterFlagCompletionFunc(flag.Collection, completion.Collection(Version))

	return &cmd
}

func newIngestCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "ingest",
		Short: "ingest documents to a collection",
		Long:  "ingest documents to a collection from either a list of files or from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			collection, _ := cmd.Flags().GetString(flag.Collection)
			batchSize, _ := cmd.Flags().GetUint64("batch-size")

			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			cfg := StreamConfig{
				Workspace:  ws,
				Collection: collection,
				BatchSize:  batchSize,
			}

			s := NewStreamer(rs, cfg)

			if len(args) == 0 {
				slog.Debug("streaming data from stdin to", "workspace", cfg.Workspace, "collection", cfg.Collection)
				count, err := s.Stream(ctx, cmd.InOrStdin())
				slog.Debug("wrote records", "count", count)
				return err
			}

			for _, a := range args {
				slog.Debug("reading", "file", a)
				count, err := s.Stream(ctx, cmd.InOrStdin())
				slog.Debug("wrote records", "count", count)
				if err != nil {
					slog.Error("failed to write", "err", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace name")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace(Version))

	cmd.Flags().String(flag.Collection, "", "collection name")
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.Collection)
	_ = cmd.RegisterFlagCompletionFunc(flag.Collection, completion.Collection(Version))

	cmd.Flags().Uint64("batch-size", 100,
		"number of documents to batch together each write")

	return &cmd
}
