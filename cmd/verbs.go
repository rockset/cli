package cmd

import "github.com/spf13/cobra"

func addVerbs(root *cobra.Command) {
	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create resources",
		Long:    "create Rockset resources",
	}

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete resources",
		Long:  "delete Rockset resource",
	}

	getCmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "get resources",
		Long:    "get Rockset resources",
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list resources",
		Long:    "list Rockset resources",
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "resume resources",
		Long:  "resume Rockset resources",
	}

	suspendCmd := &cobra.Command{
		Use:   "suspend",
		Short: "suspend resources",
		Long:  "suspend Rockset resources",
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update resources",
		Long:  "update Rockset resources",
	}

	// workspace
	createCmd.AddCommand(newCreateWorkspaceCmd())
	deleteCmd.AddCommand(newDeleteWorkspaceCmd())
	getCmd.AddCommand(newGetWorkspaceCmd())
	listCmd.AddCommand(newListWorkspacesCmd())

	// sample
	sampleCmd := &cobra.Command{
		Use:   "sample",
		Short: "create sample collections",
		Long:  "create sample collections",
	}

	createCmd.AddCommand(sampleCmd)
	sampleCmd.AddCommand(newCreateSampleCollectionCmd())

	// s3
	s3Cmd := &cobra.Command{
		Use:   "s3",
		Short: "create s3 resources",
		Long:  "s3 integration and collection commands",
	}

	createCmd.AddCommand(s3Cmd)
	s3Cmd.AddCommand(newCreateS3CollectionCmd())
	s3Cmd.AddCommand(newCreateS3IntegrationsCmd())

	// collection
	createCmd.AddCommand(newCreateCollectionCmd())
	deleteCmd.AddCommand(newDeleteCollectionCmd())
	getCmd.AddCommand(newGetCollectionCmd())
	listCmd.AddCommand(newListCollectionsCmd())

	// integration
	deleteCmd.AddCommand(newDeleteIntegrationsCmd())
	getCmd.AddCommand(newGetIntegrationCmd())
	listCmd.AddCommand(newListIntegrationsCmd())

	// org
	getCmd.AddCommand(newGetOrganizationCmd())

	listCmd.AddCommand(newListQueryCmd())

	// user
	getCmd.AddCommand(newGetUserCmd())
	listCmd.AddCommand(newListUsersCmd())

	// virtual instance
	createCmd.AddCommand(newCreateVirtualInstanceCmd())
	deleteCmd.AddCommand(newDeleteVirtualInstanceCmd())
	getCmd.AddCommand(newGetVirtualInstancesCmd())
	listCmd.AddCommand(newListVirtualInstancesCmd())
	resumeCmd.AddCommand(newResumeVirtualInstanceCmd())
	suspendCmd.AddCommand(newSuspendVirtualInstanceCmd())
	updateCmd.AddCommand(newUpdateVirtualInstanceCmd())

	// query lambda
	listCmd.AddCommand(newListLambdaCmd())

	// documents
	deleteCmd.AddCommand(newDeleteDocumentsCmd())

	listCmd.AddCommand(newListConfigCmd())
	updateCmd.AddCommand(newUpdateConfigCmd())

	root.AddCommand(createCmd)
	root.AddCommand(deleteCmd)
	root.AddCommand(getCmd)
	root.AddCommand(listCmd)
	root.AddCommand(resumeCmd)
	root.AddCommand(suspendCmd)
	root.AddCommand(updateCmd)
	root.AddCommand(newVersionCmd())

	root.AddCommand(newQueryCmd())
	root.AddCommand(newIngestCmd())
}
