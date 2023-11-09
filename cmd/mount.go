package cmd

import (
	"fmt"
	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/lookup"
	"github.com/rockset/cli/sort"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/spf13/cobra"
	"strings"
)

func NewListMountsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "mounts [NAME | ID]",
		Aliases:           []string{"m", "mount"},
		Args:              cobra.ExactArgs(1),
		Short:             "list collection mounts for a virtual instance",
		Annotations:       group("mount"),
		ValidArgsFunction: completion.VirtualInstance(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			id, err := lookup.VirtualInstanceNameOrIDtoID(ctx, rs, args[0])
			if err != nil {
				return err
			}

			mounts, err := rs.ListCollectionMounts(ctx, id)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.CollectionMount]{
				LessFuncs: []func(p1 *openapi.CollectionMount, p2 *openapi.CollectionMount) bool{
					sort.ByCollectionPath[*openapi.CollectionMount],
				},
			}
			ms.Sort(mounts)

			return formatList(cmd, format.ToInterfaceArray(mounts))
		},
	}

	return &cmd
}

func NewGetMountCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "mount PATH",
		Aliases:           []string{"m"},
		Args:              cobra.ExactArgs(1),
		Short:             "get collection mount information",
		Annotations:       group("mount"),
		ValidArgsFunction: completion.CollectionMount(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			vi, _ := cmd.Flags().GetString(flag.VI)
			id, err := lookup.VirtualInstanceNameOrIDtoID(ctx, rs, vi)
			if err != nil {
				return err
			}

			mount, err := rs.GetCollectionMount(ctx, id, args[0])
			if err != nil {
				return err
			}

			return formatOne(cmd, mount)
		},
	}

	cmd.Flags().String(flag.VI, "", "virtual instance id or name")
	cmd.MarkFlagRequired(flag.VI)
	_ = cmd.RegisterFlagCompletionFunc(flag.VI, completion.VirtualInstance(Version))

	return &cmd
}

func NewMountCollectionsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "mount PATH",
		Aliases:     []string{"m"},
		Args:        cobra.MinimumNArgs(1),
		Short:       "mount one or more collections on a virtual instance",
		Annotations: group("mount"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			vi, _ := cmd.Flags().GetString(flag.VI)
			id, err := lookup.VirtualInstanceNameOrIDtoID(ctx, rs, vi)
			if err != nil {
				return err
			}

			mounts, err := rs.MountCollections(ctx, id, args)
			if err != nil {
				return err
			}

			var mounted = make([]string, len(mounts))
			for i, mount := range mounts {
				mounted[i] = mount.GetCollectionPath()
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "mounted %s on %s\n", strings.Join(mounted, ", "), vi)

			return nil
		},
	}

	cmd.Flags().String(flag.VI, "", "virtual instance id or name")
	cmd.MarkFlagRequired(flag.VI)
	_ = cmd.RegisterFlagCompletionFunc(flag.VI, completion.VirtualInstance(Version))

	return &cmd
}

func NewUnmountCollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "unmount PATH",
		Aliases:           []string{"m"},
		Args:              cobra.MinimumNArgs(1),
		Short:             "unmount a collection from a virtual instance",
		Annotations:       group("mount"),
		ValidArgsFunction: completion.CollectionMount(Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd, Version)
			if err != nil {
				return err
			}

			vi, _ := cmd.Flags().GetString(flag.VI)
			id, err := lookup.VirtualInstanceNameOrIDtoID(ctx, rs, vi)
			if err != nil {
				return err
			}

			mount, err := rs.UnmountCollection(ctx, id, args[0])
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "unmounted %s from %s\n", mount.GetCollectionPath(), vi)

			return nil
		},
	}

	cmd.Flags().String(flag.VI, "", "virtual instance id or name")
	cmd.MarkFlagRequired(flag.VI)
	_ = cmd.RegisterFlagCompletionFunc(flag.VI, completion.VirtualInstance(Version))

	return &cmd
}
