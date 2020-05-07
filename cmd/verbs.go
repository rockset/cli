package cmd

import "github.com/spf13/cobra"

func addVerbs(root *cobra.Command) {
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

	executeCmd := &cobra.Command{
		Use:   "execute",
		Aliases: []string{"exec"},
		Short: "execute sub-command",
		Long:  "execute Rockset resource",
	}

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

	// add workspace commands
	createCmd.AddCommand(newCreateWorkspaceCmd())
	deleteCmd.AddCommand(newDeleteWorkspaceCmd())
	getCmd.AddCommand(newGetWorkspaceCmd())
	listCmd.AddCommand(newListWorkspaceCmd())

	executeCmd.AddCommand(newExecuteLambdaCmd())

	root.AddCommand(createCmd)
	root.AddCommand(deleteCmd)
	root.AddCommand(executeCmd)
	root.AddCommand(getCmd)
	root.AddCommand(listCmd)
}
