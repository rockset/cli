package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/option"
)

func newGetIntegrationCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "integration",
		Short:       "get integration",
		Long:        "get Rockset integration",
		Annotations: group("integration"),
		Args:        cobra.ExactArgs(1),
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

			return formatOne(cmd, i)
		},
	}
}

func newListIntegrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "integrations",
		Aliases:     []string{"integration"},
		Args:        cobra.NoArgs,
		Short:       "list integration",
		Long:        "list Rockset integrations",
		Annotations: group("integration"),
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

			ms := sort.Multi[openapi.Integration]{
				LessFuncs: []func(p1 *openapi.Integration, p2 *openapi.Integration) bool{
					sort.ByName[*openapi.Integration],
				},
			}
			ms.Sort(list)

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}
}

func newCreateS3IntegrationsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "integration NAME",
		Short:       "create S3 integration",
		Long:        "create S3 integration",
		Args:        cobra.ExactArgs(1),
		Annotations: group("integration"),
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

	cmd.Flags().String(RoleARNFlag, "", "AWS IAM role ARN")
	_ = cobra.MarkFlagRequired(cmd.Flags(), RoleARNFlag)

	return &cmd
}

func newDeleteIntegrationsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "integration NAME",
		Short:       "delete integration",
		Long:        "delete an integration",
		Args:        cobra.ExactArgs(1),
		Annotations: group("integration"),
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

	return &cmd
}
