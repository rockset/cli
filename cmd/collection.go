package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rockset/rockset-go-client"
)

func newDeleteCollectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "collection",
		Aliases: []string{"coll"},
		Short:   "delete collection",
		Long:    "delete Rockset collection",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockset.NewClient(rocksetAPI(cmd))
			if err != nil {
				return err
			}

			ws, name, err := toWsAndName(args)
			if err != nil {
				return err
			}

			err = rs.DeleteCollection(ctx, ws, name)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' deleted\n", ws, name)
			return nil
		},
	}
}

func newGetCollectionCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collection",
		Aliases: []string{"coll"},
		Short:   "get collection",
		Long:    "get Rockset collection",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wide, _ := cmd.Flags().GetBool("wide")

			rs, err := rockset.NewClient(rocksetAPI(cmd))
			if err != nil {
				return err
			}

			ws, name, err := toWsAndName(args)
			if err != nil {
				return err
			}

			c, err := rs.GetCollection(ctx, ws, name)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			f.Format(wide, c)
			return nil
		},
	}
	c.Flags().Bool("wide", false, "display more information")

	return &c
}

func newListCollectionsCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collections",
		Aliases: []string{"collection", "coll"},
		Short:   "list collections",
		Long:    "list Rockset collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wide, _ := cmd.Flags().GetBool("wide")

			rs, err := rockset.NewClient(rocksetAPI(cmd))
			if err != nil {
				return err
			}

			list, err := rs.ListCollections(ctx)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			f.FormatList(wide, format.ToInterfaceArray(list))
			return nil
		},
	}

	c.Flags().Bool("wide", false, "display more information")

	return &c
}

func toWsAndName(args []string) (string, string, error) {
	var ws string
	var name string
	if len(args) == 1 {
		fields := strings.Split(args[0], ".")
		if len(fields) != 2 {
			return "", "", fmt.Errorf("could not split %s into workspace and collection", args[0])
		}
		ws = fields[0]
		name = fields[1]
	} else {
		ws = args[0]
		name = args[1]
	}
	return ws, name, nil
}
