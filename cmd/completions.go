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

	return list, cobra.ShellCompDirectiveDefault
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

	return list, cobra.ShellCompDirectiveDefault
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

	return list, cobra.ShellCompDirectiveDefault
}
