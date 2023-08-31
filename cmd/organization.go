package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func newGetOrganizationCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "organization",
		Aliases: []string{"org"},
		Short:   "get organization",
		Long:    "get Rockset organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wide, _ := cmd.Flags().GetBool("wide")

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			c, err := rs.GetOrganization(ctx)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			if err = f.Format(wide, c); err != nil {
				slog.Error("failed to format data", err)
			}
			return nil
		},
	}
	c.Flags().Bool("wide", false, "display more information")

	return &c
}
