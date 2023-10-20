package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/option"
)

// TODO when implementing create apikey:
//  the role should have auto-completion for the list of roles the user has access to

func newListAPIKeysCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "apikeys [USER]",
		Aliases:     []string{"ak", "api", "apikey"},
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

			ms := sort.Multi[openapi.ApiKey]{
				LessFuncs: []func(p1 *openapi.ApiKey, p2 *openapi.ApiKey) bool{
					sort.ByName[*openapi.ApiKey],
				},
			}
			ms.Sort(keys)

			return formatList(cmd, format.ToInterfaceArray(keys))
		},
	}

	return &cmd
}

func newGetAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
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

	return &cmd
}
