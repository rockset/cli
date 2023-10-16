package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func newDeleteDocumentsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "documents",
		Aliases: []string{"doc", "docs"},
		Short:   "delete documents",
		Long:    "delete documents from a collection",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			coll, _ := cmd.Flags().GetString("collection")

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

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

	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace name")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().String("collection", "", "collection name")
	_ = cmd.MarkFlagRequired("collection")

	return &cmd
}

func newIngestCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "ingest",
		Short: "ingest documents to a collection",
		Long:  "ingest documents to a collection from either a list of files or from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			coll, _ := cmd.Flags().GetString("collection")
			batchSize, _ := cmd.Flags().GetUint64("batch-size")

			if coll == "" {
				return fmt.Errorf("must specify --collection")
			}

			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			cfg := StreamConfig{
				Workspace:  ws,
				Collection: coll,
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
					slog.Error("failed to write", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace name")
	_ = cmd.RegisterFlagCompletionFunc(WorkspaceFlag, workspaceCompletion)

	cmd.Flags().String("collection", "", "collection name")
	cmd.Flags().Uint64("batch-size", 100,
		"number of documents to batch together each write")

	return &cmd
}
