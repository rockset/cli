package cmd

import (
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
	"log"
)

func newGetIntegrationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "integration",
		Short: "get integration",
		Long:  "get Rockset integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			i, err := rs.GetIntegration(ctx, args[0])
			if err != nil {
				return err
			}

			//f, err := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			//if err != nil {
			//	return err
			//}
			//
			//f.Workspace(ws)
			log.Printf("integration: %+v", i)
			if i.S3 != nil {
				log.Printf("s3: %+v", *i.S3)
				if i.S3.AwsAccessKey != nil {
					log.Printf("aws key: %+v", *i.S3.AwsAccessKey)
				}
			}
			if i.Gcs != nil {
				log.Printf("gcs: %+v", i.Gcs.GcpServiceAccount.ServiceAccountKeyFileJson)
			}
			return nil
		},
	}
}

func newListIntegrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "integrations",
		Aliases: []string{"integration"},
		Short:   "list integration",
		Long:    "list Rockset integrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient()
			if err != nil {
				return err
			}

			list, err := rs.ListIntegrations(ctx)
			if err != nil {
				return err
			}

			//f, err := format.FormatterFor(cmd.OutOrStdout(), "table", true)
			//if err != nil {
			//	return err
			//}
			//
			//f.Workspace(ws)
			for _, i := range list {
				log.Printf("integration: %+v", i)
			}
			return nil
		},
	}
}
