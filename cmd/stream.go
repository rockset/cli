package cmd

import (
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
)

func newStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stream",
		Short: "stream data to a collection",
		Long:  "stream data to a collection from either a list of files or from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			_ = ctx
			_ = rs

			if len(args) > 0 {

			} else {
				cmd.InOrStdin()
			}

			return nil
		},
	}
}
