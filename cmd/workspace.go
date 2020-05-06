package cmd

import (
	"fmt"
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
)

func createWorkspaceCmd(cmd *cobra.Command, args []string) error {
	rs, err := rockset.NewClient(rockset.FromEnv())
	if err != nil {
		return err
	}

	// safe to ignore the error
	desc, _ := cmd.Flags().GetString("description")

	_, _, err = rs.CreateWorkspace(args[0], desc)
	if err != nil {
		if err, ok := rockset.AsRocksetError(err); ok {
			return fmt.Errorf("%s", err.Message)
		}
		return err
	}

	fmt.Printf("workspace '%s' created\n", args[0])
	return nil
}

func deleteWorkspaceCmd(cmd *cobra.Command, args []string) error {
	rs, err := rockset.NewClient(rockset.FromEnv())
	if err != nil {
		return err
	}

	_, _, err = rs.DeleteWorkspace(args[0])
	if err != nil {
		if err, ok := rockset.AsRocksetError(err); ok {
			return fmt.Errorf("%s", err.Message)
		}
		return err
	}

	fmt.Printf("workspace '%s' deleted\n", args[0])
	return nil
}

func getWorkspaceCmd(cmd *cobra.Command, args []string) error {
	rs, err := rockset.NewClient(rockset.FromEnv())
	if err != nil {
		return err
	}

	ws, _, err := rs.GetWorkspace(args[0])
	if err != nil {
		if err, ok := rockset.AsRocksetError(err); ok {
			return fmt.Errorf("%s", err.Message)
		}
		return err
	}

	fmt.Printf("workspace info: %+v\n", ws)
	return nil
}

func listWorkspaceCmd(cmd *cobra.Command, args []string) error {
	rs, err := rockset.NewClient(rockset.FromEnv())
	if err != nil {
		return err
	}

	list, _, err := rs.ListWorkspaces()
	if err != nil {
		if err, ok := rockset.AsRocksetError(err); ok {
			return fmt.Errorf("%s", err.Message)
		}
		return err
	}

	for _, ws := range list {
		fmt.Printf("%+v\n", ws)
	}
	return nil
}

func addWorkspaceCommands(createCmd, deleteCmd, getCmd, listCmd *cobra.Command) {
	createWs := &cobra.Command{
		Use:   "workspace",
		Short: "create workspace",
		Long:  "create Rockset workspace",
		Args:  cobra.ExactArgs(1),
		RunE:  createWorkspaceCmd,
	}

	deleteWs := &cobra.Command{
		Use:   "workspace",
		Short: "delete workspace",
		Long:  "delete Rockset workspace",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteWorkspaceCmd,
	}

	listWs := &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"workspace"},
		Short:   "list workspaces",
		Long:    "list Rockset workspace",
		RunE:    listWorkspaceCmd,
	}

	getWs := &cobra.Command{
		Use:   "workspace",
		Short: "get workspace",
		Long:  "get Rockset workspace",
		Args:  cobra.ExactArgs(1),
		RunE:  getWorkspaceCmd,
	}

	createWs.Flags().StringP("description", "d", "", "workspace description")
	createCmd.AddCommand(createWs)
	deleteCmd.AddCommand(deleteWs)
	getCmd.AddCommand(getWs)
	listCmd.AddCommand(listWs)
}
