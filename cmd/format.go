package cmd

import (
	"github.com/spf13/cobra"

	"github.com/rockset/cli/format"
)

func formatOne(cmd *cobra.Command, a any) error {
	wide, _ := cmd.Flags().GetBool(HeaderFlag)

	f, err := formatterFor(cmd)
	if err != nil {
		return err
	}

	return f.Format(wide, a)
}

func formatList(cmd *cobra.Command, a []any) error {
	wide, _ := cmd.Flags().GetBool(HeaderFlag)

	f, err := formatterFor(cmd)
	if err != nil {
		return err
	}

	return f.FormatList(wide, a)
}

func formatterFor(cmd *cobra.Command) (format.Formatter, error) {
	f, _ := cmd.Flags().GetString(FormatFlag)
	header, _ := cmd.Flags().GetBool(HeaderFlag)

	return format.FormatterFor(cmd.OutOrStdout(), format.Format(f), header)
}
