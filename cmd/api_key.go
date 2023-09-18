package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/option"
)

func newListAPIKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "apikey [USER]",
		Aliases:     []string{"ak", "api"},
		Args:        cobra.RangeArgs(0, 1),
		Short:       "list apikeys for the current user, or the specified USER",
		Annotations: group("apikey"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var opts []option.APIKeyOption
			if len(args) > 0 {
				opts = append(opts, option.ForUser(args[0]))
			}

			keys, err := rs.ListAPIKeys(ctx, opts...)
			if err != nil {
				return err
			}

			return formatList(cmd, format.ToInterfaceArray(keys))
		},
	}

	return cmd
}

func newGetAPIKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "apikey NAME",
		Aliases:     []string{"ak"},
		Args:        cobra.ExactArgs(1),
		Short:       "get apikey information",
		Annotations: group("apikey"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			apikey, err := rs.GetAPIKey(ctx, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, apikey)
		},
	}

	return cmd
}
