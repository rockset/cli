package cmd

import (
	"fmt"
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
	"log"
)

func newStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "stream data to a collection",
		Long:  "stream data to a collection from either a list of files or from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			ws, _ := cmd.Flags().GetString("workspace")
			coll, _ := cmd.Flags().GetString("collection")
			batchSize, _ := cmd.Flags().GetUint64("batch-size")

			if coll == "" {
				return fmt.Errorf("must specify --collection")
			}

			ctx := cmd.Context()
			rs, err := rockset.NewClient()
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
