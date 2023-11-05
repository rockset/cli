package cmd

import (
	"github.com/spf13/cobra"
)

func addVerbs(root *cobra.Command) {
	authCmd := cobra.Command{
		Use:   "auth",
		Short: "authenticate",
		Long:  "authenticate using an bearer token or an apikey",
	}

	createCmd := cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "create resources",
		Long:    "create Rockset resources",
	}

	deleteCmd := cobra.Command{
		Use:   "delete",
		Short: "delete resources",
		Long:  "delete Rockset resource",
	}

	executeCmd := cobra.Command{
		Use:     "execute",
		Aliases: []string{"e"},
		Short:   "execute query",
		Long:    "execute Rockset query",
	}

	getCmd := cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "get resources",
		Long:    "get Rockset resources",
	}

	listCmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list resources",
		Long:    "list Rockset resources",
	}

	resumeCmd := cobra.Command{
		Use:   "resume",
		Short: "resume resources",
		Long:  "resume Rockset resources",
	}

	suspendCmd := cobra.Command{
		Use:   "suspend",
		Short: "suspend resources",
		Long:  "suspend Rockset resources",
	}

	tailCmd := cobra.Command{
		Use:   "tail",
		Short: "tail collections",
		Long:  "tail Rockset collections",
	}

	updateCmd := cobra.Command{
		Use:   "update",
		Short: "update resources",
		Long:  "update Rockset resources",
	}

	useCmd := cobra.Command{
		Use:   "use",
		Short: "use configuration",
		Long:  "use a Rockset apikey and server configuration",
	}

	// sample
	sampleCmd := cobra.Command{
		Use:   "sample",
		Short: "create sample collections",
		Long:  "create sample collections",
	}

	createCmd.AddCommand(&sampleCmd)

	// s3
	s3Cmd := cobra.Command{
		Use:   "s3",
		Short: "create s3 resources",
		Long:  "s3 integration and collection commands",
	}

	// authentication
	authCmd.AddCommand(newAuthLoginCmd())
	authCmd.AddCommand(newAuthKeyCmd())
	authCmd.AddCommand(newAuthRefreshCmd())

	createCmd.AddCommand(&s3Cmd)
	s3Cmd.AddCommand(newCreateS3CollectionCmd())
	s3Cmd.AddCommand(newCreateS3IntegrationsCmd())

	// collections
	createCmd.AddCommand(newCreateCollectionCmd())
	deleteCmd.AddCommand(newDeleteCollectionCmd())
	getCmd.AddCommand(newGetCollectionCmd())
	listCmd.AddCommand(newListCollectionsCmd())
	listCmd.AddCommand(newListQueryCmd())
	sampleCmd.AddCommand(newCreateSampleCollectionCmd())
	tailCmd.AddCommand(newCreateTailCollectionCmd())

	// integrations
	deleteCmd.AddCommand(newDeleteIntegrationsCmd())
	getCmd.AddCommand(newGetIntegrationCmd())
	listCmd.AddCommand(newListIntegrationsCmd())

	// org
	getCmd.AddCommand(newGetOrganizationCmd())

	// users
	getCmd.AddCommand(newGetUserCmd())
	listCmd.AddCommand(newListUsersCmd())

	// views
	getCmd.AddCommand(newGetViewCmd())
	listCmd.AddCommand(newListViewsCmd())

	// virtual instances
	createCmd.AddCommand(newCreateVirtualInstanceCmd())
	deleteCmd.AddCommand(newDeleteVirtualInstanceCmd())
	getCmd.AddCommand(newGetVirtualInstancesCmd())
	listCmd.AddCommand(newListVirtualInstancesCmd())
	resumeCmd.AddCommand(newResumeVirtualInstanceCmd())
	suspendCmd.AddCommand(newSuspendVirtualInstanceCmd())
	updateCmd.AddCommand(newUpdateVirtualInstanceCmd())

	// aliases
	getCmd.AddCommand(NewGetAliasCmd())
	listCmd.AddCommand(NewListAliasesCmd())
	createCmd.AddCommand(NewCreateAliasCmd())
	deleteCmd.AddCommand(NewDeleteAliasCmd())
	updateCmd.AddCommand(NewUpdateAliasCmd())

	// roles
	getCmd.AddCommand(newGetRoleCommand())
	listCmd.AddCommand(newListRolesCommand())

	// API keys
	createCmd.AddCommand(NewCreateAPIKeyCmd())
	deleteCmd.AddCommand(NewDeleteAPIKeyCmd())
	getCmd.AddCommand(NewGetAPIKeyCmd())
	listCmd.AddCommand(NewListAPIKeysCmd())
	updateCmd.AddCommand(newUpdateAPIKeyCmd())

	// query lambda
	createCmd.AddCommand(newCreateQueryLambdaCmd())
	deleteCmd.AddCommand(newDeleteQueryLambdaCmd())
	updateCmd.AddCommand(newUpdateQueryLambdaCmd())
	executeCmd.AddCommand(NewExecuteQueryLambdaCmd())
	getCmd.AddCommand(newGetQueryLambdaCmd())
	listCmd.AddCommand(newListQueryLambdasCmd())

	// documents
	deleteCmd.AddCommand(newDeleteDocumentsCmd())

	// workspace
	createCmd.AddCommand(newCreateWorkspaceCmd())
	deleteCmd.AddCommand(newDeleteWorkspaceCmd())
	getCmd.AddCommand(NewGetWorkspaceCmd())
	listCmd.AddCommand(newListWorkspacesCmd())

	// context
	createCmd.AddCommand(newCreateContextCmd())
	deleteCmd.AddCommand(newDeleteContextCmd())
	listCmd.AddCommand(newListContextsCmd())
	useCmd.AddCommand(newUseContextCmd())

	root.AddCommand(&authCmd)
	root.AddCommand(&createCmd)
	root.AddCommand(&deleteCmd)
	root.AddCommand(&executeCmd)
	root.AddCommand(&getCmd)
	root.AddCommand(&listCmd)
	root.AddCommand(&resumeCmd)
	root.AddCommand(&suspendCmd)
	root.AddCommand(&tailCmd)
	root.AddCommand(&updateCmd)
	root.AddCommand(&useCmd)
	root.AddCommand(newVersionCmd())

	root.AddCommand(newQueryCmd())
	root.AddCommand(newIngestCmd())

	root.AddCommand(newTestCmd())

	// TODO set help func for the root command to show commands grouped by the resource they operate on
}
