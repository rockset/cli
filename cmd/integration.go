package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/rockset/rockset-go-client/option"
)

func newGetIntegrationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "integration",
		Short: "get integration",
		Long:  "get Rockset integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
			slog.Debug("integration: %+v", i)
			if i.S3 != nil {
				slog.Debug("s3", "struct", *i.S3)
				if i.S3.AwsAccessKey != nil {
					slog.Debug("aws", "key", *i.S3.AwsAccessKey)
				}
			}

			return nil
		},
	}
}

func newListIntegrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "integrations",
		Aliases: []string{"integration"},
		Args:    cobra.NoArgs,
		Short:   "list integration",
		Long:    "list Rockset integrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
				slog.Debug("integration", "struct", i)
			}
			return nil
		},
	}
}

func newCreateS3IntegrationsCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "integration NAME",
		Short: "create S3 integration",
		Long:  "create S3 integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			role, _ := cmd.Flags().GetString(RoleARNFlag)
			result, err := rs.CreateS3Integration(ctx, args[0], option.AWSRole(role))
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "integration'%s' is created\n", result.Name)

			return nil
		},
	}

	c.Flags().String(RoleARNFlag, "", "AWS IAM role ARN")
	_ = cobra.MarkFlagRequired(c.Flags(), RoleARNFlag)

	return &c
}

func newDeleteIntegrationsCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "integration NAME",
		Short: "delete integration",
		Long:  "delete an integration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			if err = rs.DeleteIntegration(ctx, args[0]); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "integration '%s' deleted\n", args[0])

			return nil
		},
	}

	return &c
}
