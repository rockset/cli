package completion

import (
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/lookup"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"
)

func Collection(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var options []option.ListCollectionOption
		if ws, _ := cmd.Flags().GetString(flag.Workspace); ws != "" {
			options = append(options, option.WithWorkspace(ws))
		}

		collections, err := rs.ListCollections(cmd.Context(), options...)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(collections))
		for i, ws := range collections {
			list[i] = ws.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func Integration(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		integrations, err := rs.ListIntegrations(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(integrations))
		for i, ws := range integrations {
			list[i] = ws.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func Workspace(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		workspaces, err := rs.ListWorkspaces(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(workspaces))
		for i, ws := range workspaces {
			list[i] = ws.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func Lambda(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ws, err := cmd.Flags().GetString(flag.Workspace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var options []option.ListQueryLambdaOption
		if ws != "" {
			options = append(options, option.WithQueryLambdaWorkspace(ws))
		}

		versions, err := rs.ListQueryLambdas(cmd.Context(), options...)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func Alias(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ws, err := cmd.Flags().GetString(flag.Workspace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var options []option.ListAliasesOption
		if ws != "" {
			options = append(options, option.WithAliasWorkspace(ws))
		}

		versions, err := rs.ListAliases(cmd.Context(), options...)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func LambdaVersion(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ws, err := cmd.Flags().GetString(flag.Workspace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		versions, err := rs.ListQueryLambdaVersions(cmd.Context(), ws, args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func LambdaTag(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ws, err := cmd.Flags().GetString(flag.Workspace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		versions, err := rs.ListQueryLambdaTags(cmd.Context(), ws, args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetTagName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func Role(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		versions, err := rs.ListRoles(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetRoleName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func VirtualInstance(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		versions, err := rs.ListVirtualInstances(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func View(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ws, err := cmd.Flags().GetString(flag.Workspace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var options []option.ListViewOption
		if ws != "" {
			options = append(options, option.WithViewWorkspace(ws))
		}

		versions, err := rs.ListViews(cmd.Context(), options...)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(versions))
		for i, v := range versions {
			list[i] = v.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func Email(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		users, err := rs.ListUsers(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(users))
		for i, user := range users {
			list[i] = user.GetEmail()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func APIKey(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var options []option.APIKeyOption
		if email, _ := cmd.Flags().GetString(flag.Email); email != "" {
			options = append(options, option.ForUser(email))
		}

		keys, err := rs.ListAPIKeys(cmd.Context(), options...)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(keys))
		for i, key := range keys {
			list[i] = key.GetName()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}

func CollectionMount(version string) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ctx := cmd.Context()
		rs, err := config.Client(cmd, version)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		vi, _ := cmd.Flags().GetString(flag.VI)
		id, err := lookup.VirtualInstanceNameOrIDtoID(ctx, rs, vi)

		mounts, err := rs.ListCollectionMounts(ctx, id)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := make([]string, len(mounts))
		for i, mount := range mounts {
			list[i] = mount.GetCollectionPath()
		}

		return list, cobra.ShellCompDirectiveNoFileComp
	}
}
