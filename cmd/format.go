package cmd

import (
	"github.com/rockset/cli/flag"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/format"
)

func formatOne(cmd *cobra.Command, a any) error {
	wide, _ := cmd.Flags().GetBool(flag.Header)
	selector, _ := cmd.Flags().GetString(flag.Selector)

	f, err := formatterFor(cmd)
	if err != nil {
		return err
	}

	var s format.Selector
	if selector != "" {
		s, err = format.ParseSelectionString(selector)
		if err != nil {
			return err
		}
	}
	return f.Format(wide, s, a)
}

func formatList(cmd *cobra.Command, a []any) error {
	wide, _ := cmd.Flags().GetBool(flag.Wide)
	selector, _ := cmd.Flags().GetString(flag.Selector)

	f, err := formatterFor(cmd)
	if err != nil {
		return err
	}

	var s format.Selector
	if selector != "" {
		s, err = format.ParseSelectionString(selector)
		if err != nil {
			return err
		}
	}

	return f.FormatList(wide, s, a)
}

func formatterFor(cmd *cobra.Command) (format.Formatter, error) {
	f, _ := cmd.Flags().GetString(flag.Format)
	header, _ := cmd.Flags().GetBool(flag.Header)

	return format.FormatterFor(cmd.OutOrStdout(), format.Format(f), header)
}
