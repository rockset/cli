package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client/option"
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
		ValidArgsFunction: emailCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
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
		ValidArgsFunction: apikeyCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var options []option.APIKeyOption
			if email, _ := cmd.Flags().GetString(EmailFlag); email != "" {
				options = append(options, option.ForUser(email))
			}

			apikey, err := rs.GetAPIKey(ctx, args[0], options...)
			if err != nil {
				return err
			}

			return formatOne(cmd, apikey)
		},
	}

	cmd.Flags().String(EmailFlag, "", "the email address of the user who's key to get, defaults to self")
	_ = cmd.RegisterFlagCompletionFunc(EmailFlag, emailCompletion)

	return &cmd
}

func NewDeleteAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "apikey NAME",
		Aliases:           []string{"ak"},
		Args:              cobra.ExactArgs(1),
		Short:             "delete an apikey",
		Annotations:       group("apikey"),
		ValidArgsFunction: apikeyCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var options []option.APIKeyOption
			if email, _ := cmd.Flags().GetString(EmailFlag); email != "" {
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

	cmd.Flags().String(EmailFlag, "", "the email address of the user who's key to delete, defaults to self")
	_ = cmd.RegisterFlagCompletionFunc(EmailFlag, emailCompletion)

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
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var options []option.APIKeyRoleOption
			if role, _ := cmd.Flags().GetString(RoleFlag); role != "" {
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

	cmd.Flags().String(RoleFlag, "", "role for the apikey")
	_ = cmd.RegisterFlagCompletionFunc(RoleFlag, roleCompletion)

	return &cmd
}

func newUpdateAPIKeyCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "apikey NAME",
		Aliases:           []string{"ak"},
		Args:              cobra.ExactArgs(1),
		Short:             "update the state of an apikey",
		Annotations:       group("apikey"),
		ValidArgsFunction: apikeyCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var options []option.APIKeyOption
			if email, _ := cmd.Flags().GetString(EmailFlag); email != "" {
				options = append(options, option.ForUser(email))
			}
			state, _ := cmd.Flags().GetString(StateFlag)
			options = append(options, option.State(option.KeyState(state)))

			apikey, err := rs.UpdateAPIKey(ctx, args[0], options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "apikey %s updated to %s\n", apikey.GetName(), apikey.GetState())

			return nil
		},
	}

	cmd.Flags().String(EmailFlag, "", "the email address of the user who's key to update, defaults to self")
	_ = cmd.RegisterFlagCompletionFunc(EmailFlag, emailCompletion)

	cmd.Flags().String(StateFlag, "",
		fmt.Sprintf("the state of the apikey, either %s or %s", option.KeyActive, option.KeySuspended))
	_ = cmd.RegisterFlagCompletionFunc(StateFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{string(option.KeyActive), string(option.KeySuspended)}, cobra.ShellCompDirectiveNoFileComp
		})
	_ = cobra.MarkFlagRequired(cmd.Flags(), StateFlag)

	return &cmd
}
