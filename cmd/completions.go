package cmd

import (
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"
)

func collectionCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	rs, err := rockClient(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var options []option.ListCollectionOption
	if ws, _ := cmd.Flags().GetString(WorkspaceFlag); ws != "" {
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

func integrationCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	rs, err := rockClient(cmd)
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

func workspaceCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	rs, err := rockClient(cmd)
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

func lambdaCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ws, err := cmd.Flags().GetString(WorkspaceFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	rs, err := rockClient(cmd)
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

func aliasCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ws, err := cmd.Flags().GetString(WorkspaceFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	rs, err := rockClient(cmd)
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

func lambdaVersionsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ws, err := cmd.Flags().GetString(WorkspaceFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	rs, err := rockClient(cmd)
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

func lambdaTagsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ws, err := cmd.Flags().GetString(WorkspaceFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	rs, err := rockClient(cmd)
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

func roleCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	rs, err := rockClient(cmd)
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

func virtualInstanceCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	rs, err := rockClient(cmd)
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
