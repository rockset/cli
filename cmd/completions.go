package cmd

import (
	"github.com/spf13/cobra"
)

func collectionCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	rs, err := rockClient(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	collections, err := rs.ListCollections(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	list := make([]string, len(collections))
	for i, ws := range collections {
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
