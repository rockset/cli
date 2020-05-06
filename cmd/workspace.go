package cmd

import (
	"fmt"
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
)

func newCreateWorkspaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "create workspace",
		Long:  "create Rockset workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rs, err := rockset.NewClient(rockset.FromEnv())
			if err != nil {
				return err
			}

			// safe to ignore the error it is added below
			desc, _ := cmd.Flags().GetString("description")

			_, _, err = rs.CreateWorkspace(args[0], desc)
			if err != nil {
				if err, ok := rockset.AsRocksetError(err); ok {
					return fmt.Errorf("%s", err.Message)
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(),"workspace '%s' created\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringP("description", "d", "", "workspace description")
	return cmd
}

func newDeleteWorkspaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "workspace",
		Short: "delete workspace",
		Long:  "delete Rockset workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}
}

func newGetWorkspaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "workspace",
		Short: "get workspace",
		Long:  "get Rockset workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			fmt.Fprintf(cmd.OutOrStdout(),"workspace info: %+v\n", ws)
			return nil
		},
	}
}

func newListWorkspaceCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"workspace"},
		Short:   "list workspaces",
		Long:    "list Rockset workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}
}

