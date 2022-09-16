package cmd

import "github.com/spf13/cobra"

func addVerbs(root *cobra.Command) {
	//addCmd := &cobra.Command{
	//	Use:   "add",
	//	Short: "add sub-command",
	//	Long:  "add Rockset resource",
	//}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create sub-command",
		Long:  "create Rockset resource",
	}

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete sub-command",
		Long:  "delete Rockset resource",
	}

	//executeCmd := &cobra.Command{
	//	Use:   "execute",
	//	Aliases: []string{"exec"},
	//	Short: "execute sub-command",
	//	Long:  "execute Rockset resource",
	//}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list sub-command",
		Long:  "list Rockset resources",
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "get sub-command",
		Long:  "get Rockset resource",
	}

	streamCmd := &cobra.Command{
		Use:   "stream",
		Short: "stream sub-command",
		Long:  "stream data to Rockset",
	}

	// workspace
	createCmd.AddCommand(newCreateWorkspaceCmd())
	deleteCmd.AddCommand(newDeleteWorkspaceCmd())
	getCmd.AddCommand(newGetWorkspaceCmd())
	listCmd.AddCommand(newListWorkspacesCmd())

	// collection
	deleteCmd.AddCommand(newDeleteCollectionCmd())
	getCmd.AddCommand(newGetCollectionCmd())
	listCmd.AddCommand(newListCollectionsCmd())

	// integration
	getCmd.AddCommand(newGetIntegrationCmd())
	listCmd.AddCommand(newListIntegrationsCmd())

	// org
	getCmd.AddCommand(newGetOrganizationCmd())

	// user
	getCmd.AddCommand(newGetUserCmd())
	listCmd.AddCommand(newListUsersCmd())

	// query lambda
	//executeCmd.AddCommand(newExecuteLambdaCmd())
	listCmd.AddCommand(newListLambdaCmd())

	// documents
	deleteCmd.AddCommand(newDeleteDocumentsCmd())
	streamCmd.AddCommand(newStreamDocumentsCmd())

	//root.AddCommand(addCmd)
	root.AddCommand(createCmd)
	root.AddCommand(deleteCmd)
	//root.AddCommand(executeCmd)
	root.AddCommand(getCmd)
	root.AddCommand(listCmd)
	root.AddCommand(streamCmd)

}
