package cmd

import (
	"fmt"
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
	"log"
)

func newDeleteDocumentsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "documents",
		Aliases: []string{"doc", "docs"},
		Short:   "delete documents",
		Long:    "delete documents from a collection",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString("workspace")
			coll, _ := cmd.Flags().GetString("collection")

			ctx := cmd.Context()
			rs, err := rockset.NewClient(rocksetAPI(cmd))
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
					fmt.Printf("failed to delete document %s\n", d.GetId())
					continue
				}
				count++
			}

			fmt.Printf("deleted %d documents\n", count)
			if failed > 0 {
				fmt.Printf("failed to delete %d documents\n", failed)
			}

			return nil
		},
	}

	cmd.Flags().String("workspace", "commons", "workspace name")
	cmd.Flags().String("collection", "", "collection name")
	_ = cmd.MarkFlagRequired("collection")

	return &cmd
}

func newStreamDocumentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "documents",
		Aliases: []string{"doc", "docs", "document"},
		Short:   "stream documents to a collection",
		Long:    "stream documents to a collection from either a list of files or from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString("workspace")
			coll, _ := cmd.Flags().GetString("collection")
			batchSize, _ := cmd.Flags().GetUint64("batch-size")

			if coll == "" {
				return fmt.Errorf("must specify --collection")
			}

			ctx := cmd.Context()
			rs, err := rockset.NewClient(rocksetAPI(cmd))
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
				log.Printf("streaming data from stdin to %s.%s", cfg.Workspace, cfg.Collection)
				count, err := s.Stream(ctx, cmd.InOrStdin())
				log.Printf("wrote %d records", count)
				return err
			}

			for _, a := range args {
				log.Printf("reading from file %s", a)
				count, err := s.Stream(ctx, cmd.InOrStdin())
				log.Printf("wrote %d records", count)
				if err != nil {
					log.Printf("failed to write")
				}
			}

			return nil
		},
	}

	cmd.Flags().String("workspace", "common", "workspace name")
	cmd.Flags().String("collection", "", "collection name")
	cmd.Flags().Uint64("batch-size", 100,
		"number of documents to batch together each write")

	return cmd
}
