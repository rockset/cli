package cmd

import (
	"github.com/rockset/cli/config"
	"github.com/spf13/cobra"
)

func newGetOrganizationCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "organization",
		Aliases:     []string{"org"},
		Args:        cobra.NoArgs,
		Short:       "get organization",
		Long:        "get Rockset organization",
		Annotations: group("org"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			org, err := rs.GetOrganization(ctx)
			if err != nil {
				return err
			}

			return formatOne(cmd, org)
		},
	}

	cmd.Flags().Bool("wide", false, "display more information")

	return &cmd
}
