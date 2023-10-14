package cmd

import (
	"github.com/spf13/cobra"
)

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

	executeCmd := &cobra.Command{
		Use:     "execute",
		Aliases: []string{"e"},
		Short:   "execute query",
		Long:    "execute Rockset query",
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

	useCmd := &cobra.Command{
		Use:   "use",
		Short: "use configuration",
		Long:  "use a Rockset apikey and server configuration",
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

	// view
	getCmd.AddCommand(newGetViewCmd())
	listCmd.AddCommand(newListViewsCmd())

	// virtual instance
	createCmd.AddCommand(newCreateVirtualInstanceCmd())
	deleteCmd.AddCommand(newDeleteVirtualInstanceCmd())
	getCmd.AddCommand(newGetVirtualInstancesCmd())
	listCmd.AddCommand(newListVirtualInstancesCmd())
	resumeCmd.AddCommand(newResumeVirtualInstanceCmd())
	suspendCmd.AddCommand(newSuspendVirtualInstanceCmd())
	updateCmd.AddCommand(newUpdateVirtualInstanceCmd())

	getCmd.AddCommand(newGetAliasCmd())
	listCmd.AddCommand(newListAliasesCmd())

	getCmd.AddCommand(newGetRoleCommand())
	listCmd.AddCommand(newListRolesCommand())

	getCmd.AddCommand(newGetAPIKeyCmd())
	listCmd.AddCommand(newListAPIKeysCmd())

	// query lambda
	createCmd.AddCommand(newCreateQueryLambdaCmd())
	executeCmd.AddCommand(newExecuteQueryLambdaCmd())
	getCmd.AddCommand(newGetQueryLambdaCmd())
	listCmd.AddCommand(newListQueryLambdaCmd())

	// documents
	deleteCmd.AddCommand(newDeleteDocumentsCmd())

	// config
	createCmd.AddCommand(newCreateConfigCmd())
	listCmd.AddCommand(newListConfigCmd())
	useCmd.AddCommand(newUseConfigCmd())

	root.AddCommand(newAuthCmd())
	root.AddCommand(createCmd)
	root.AddCommand(deleteCmd)
	root.AddCommand(executeCmd)
	root.AddCommand(getCmd)
	root.AddCommand(listCmd)
	root.AddCommand(resumeCmd)
	root.AddCommand(suspendCmd)
	root.AddCommand(updateCmd)
	root.AddCommand(useCmd)
	root.AddCommand(newVersionCmd())

	root.AddCommand(newQueryCmd())
	root.AddCommand(newIngestCmd())

	root.AddCommand(newTestCmd())

	// TODO set help func for the root command to show commands grouped by the resource they operate on
}
