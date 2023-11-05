package cmd

import (
	"fmt"

	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

// TODO when implementing create apikey:
//  the role should have auto-completion for the list of roles the user has access to

func NewListAPIKeysCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "apikeys [USER]",
		Aliases:           []string{"ak", "api", "apikey"},
		Args:              cobra.RangeArgs(0, 1),
		Short:             "list apikeys for the current user, or the specified USER",
		Annotations:       group("apikey"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			logger.Info("args", "args", args)
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

func NewGetAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "apikey NAME",
		Aliases:           []string{"ak"},
		Args:              cobra.ExactArgs(1),
		Short:             "get apikey information",
		Annotations:       group("apikey"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var options []option.APIKeyOption
			if email, _ := cmd.Flags().GetString(flag.Email); email != "" {
				options = append(options, option.ForUser(email))
			}

			apikey, err := rs.GetAPIKey(ctx, args[0], options...)
			if err != nil {
				return err
			}

			return formatOne(cmd, apikey)
		},
	}

	cmd.Flags().String(flag.Email, "", "the email address of the user who's key to get, defaults to self")
	_ = cmd.RegisterFlagCompletionFunc(flag.Email, completion.Alias(Version))

	return &cmd
}

func NewDeleteAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "apikey NAME",
		Aliases:           []string{"ak"},
		Args:              cobra.ExactArgs(1),
		Short:             "delete an apikey",
		Annotations:       group("apikey"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var options []option.APIKeyOption
			if email, _ := cmd.Flags().GetString(flag.Email); email != "" {
				options = append(options, option.ForUser(email))
			}

			err = rs.DeleteAPIKey(ctx, args[0], options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "apikey %s deleted\n", args[0])

			return nil
		},
	}

	cmd.Flags().String(flag.Email, "", "the email address of the user who's key to delete, defaults to self")
	_ = cmd.RegisterFlagCompletionFunc(flag.Email, completion.Alias(Version))

	return &cmd
}

func NewCreateAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "apikey NAME",
		Aliases:     []string{"ak"},
		Args:        cobra.ExactArgs(1),
		Short:       "create apikey",
		Annotations: group("apikey"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var options []option.APIKeyRoleOption
			if role, _ := cmd.Flags().GetString(flag.Role); role != "" {
				options = append(options, option.WithRole(role))
			}

			apikey, err := rs.CreateAPIKey(ctx, args[0], options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "apikey %s created\n", apikey.GetName())

			return nil
		},
	}

	cmd.Flags().String(flag.Role, "", "role for the apikey")
	_ = cmd.RegisterFlagCompletionFunc(flag.Role, completion.Role(Version))

	return &cmd
}

func newUpdateAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "apikey NAME",
		Aliases:           []string{"ak"},
		Args:              cobra.ExactArgs(1),
		Short:             "update the state of an apikey",
		Annotations:       group("apikey"),
		ValidArgsFunction: completion.Alias(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			var options []option.APIKeyOption
			if email, _ := cmd.Flags().GetString(flag.Email); email != "" {
				options = append(options, option.ForUser(email))
			}
			state, _ := cmd.Flags().GetString(flag.State)
			options = append(options, option.State(option.KeyState(state)))

			apikey, err := rs.UpdateAPIKey(ctx, args[0], options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "apikey %s updated to %s\n", apikey.GetName(), apikey.GetState())

			return nil
		},
	}

	cmd.Flags().String(flag.Email, "", "the email address of the user who's key to update, defaults to self")
	_ = cmd.RegisterFlagCompletionFunc(flag.Email, completion.Alias(Version))

	cmd.Flags().String(flag.State, "",
		fmt.Sprintf("the state of the apikey, either %s or %s", option.KeyActive, option.KeySuspended))
	_ = cmd.RegisterFlagCompletionFunc(flag.State,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{string(option.KeyActive), string(option.KeySuspended)}, cobra.ShellCompDirectiveNoFileComp
		})
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.State)

	return &cmd
}
