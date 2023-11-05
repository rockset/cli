package cmd

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
)

func newCreateVirtualInstanceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Short:       "create a virtual instance",
		Long:        "create a Rockset virtual instance",
		Args:        cobra.ExactArgs(1),
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			size, _ := cmd.Flags().GetString(flag.Size)

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			var options = []option.VirtualInstanceOption{
				option.WithVirtualInstanceSize(option.VirtualInstanceSize(size)),
			}

			// TODO uncomment after go client supports setting it
			//if desc, _ := cmd.Flags().GetString(flag.Description); desc != "" {
			//options = append(options, option.WithVirtualInstanceDescription(desc))
			//}

			result, err := rs.CreateVirtualInstance(ctx, args[0], options...)
			if err != nil {
				return err
			}

			if err = waitUntilVIActive(rs, cmd, result.GetId()); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "virtual instance '%s' created\n", result.GetName())
			return nil
		},
	}

	// TODO uncomment after go client supports setting it
	//cmd.Flags().StringP(flag.Description, "d", "", "virtual instance description")
	cmd.Flags().Bool(flag.Wait, false, "wait until virtual instance is active")
	cmd.Flags().String(flag.Size, "", "virtual instance size")
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.Size)
	// TODO completion of sizes

	return &cmd
}

func newUpdateVirtualInstanceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "virtualinstance ID|NAME",
		Aliases:           []string{"vi"},
		Short:             "update a virtual instance",
		Long:              "update a Rockset virtual instance",
		Args:              cobra.ExactArgs(1),
		Annotations:       group("virtual instance"),
		ValidArgsFunction: completion.VirtualInstance,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			size, _ := cmd.Flags().GetString(flag.Size)

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			id, err := viNameOrIDtoID(ctx, rs, args[0])
			if err != nil {
				return err
			}

			logger.Info("setting option to", "size", size)
			var options = []option.VirtualInstanceOption{
				option.WithVirtualInstanceSize(option.VirtualInstanceSize(size)),
			}

			result, err := rs.UpdateVirtualInstance(ctx, id, options...)
			if err != nil {
				return err
			}

			if err = waitUntilVIActive(rs, cmd, result.GetId()); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "virtual instance '%s' updated\n", result.GetName())
			return nil
		},
	}

	// TODO uncomment after go client supports setting it
	//cmd.Flags().StringP(flag.Description, "d", "", "virtual instance description")
	cmd.Flags().String(flag.Size, "", "virtual instance size")
	cmd.Flags().Bool(flag.Wait, false, "wait until virtual instance is active")
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.Size)
	// TODO completion of sizes

	return &cmd
}

func newListVirtualInstancesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "virtualinstances",
		Aliases:     []string{"vi", "vis"},
		Args:        cobra.NoArgs,
		Short:       "list virtual instances",
		Long:        "list Rockset virtual instances",
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			list, err := rs.ListVirtualInstances(ctx)
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.VirtualInstance]{
				LessFuncs: []func(p1 *openapi.VirtualInstance, p2 *openapi.VirtualInstance) bool{
					sort.ByName[*openapi.VirtualInstance],
				},
			}
			ms.Sort(list)

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")

	return &cmd
}

func newGetVirtualInstancesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "virtualinstance ID|NAME",
		Aliases:           []string{"vi"},
		Args:              cobra.ExactArgs(1),
		Short:             "get virtual instance",
		Long:              "get Rockset virtual instances",
		Annotations:       group("virtual instance"),
		ValidArgsFunction: completion.VirtualInstance,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			id, err := viNameOrIDtoID(ctx, rs, args[0])
			if err != nil {
				return err
			}

			vi, err := rs.GetVirtualInstance(ctx, id)
			if err != nil {
				return err
			}

			return formatOne(cmd, vi)
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")

	return &cmd
}

func newDeleteVirtualInstanceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "virtualinstance ID|NAME",
		Aliases:           []string{"vi"},
		Short:             "delete virtual instance",
		Long:              "delete Rockset virtual instance",
		Args:              cobra.ExactArgs(1),
		Annotations:       group("virtual instance"),
		ValidArgsFunction: completion.VirtualInstance,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			id, err := viNameOrIDtoID(ctx, rs, args[0])
			if err != nil {
				return err
			}

			result, err := rs.DeleteVirtualInstance(ctx, id)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "virtual instance '%s' deleted\n", result.GetName())

			return nil
		},
	}

	return &cmd
}

func newSuspendVirtualInstanceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "virtualinstance ID|NAME",
		Aliases:           []string{"vi"},
		Short:             "suspend virtual instance",
		Long:              "suspend Rockset virtual instance",
		Args:              cobra.ExactArgs(1),
		Annotations:       group("virtual instance"),
		ValidArgsFunction: completion.VirtualInstance,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			id, err := viNameOrIDtoID(ctx, rs, args[0])
			if err != nil {
				return err
			}

			result, err := rs.SuspendVirtualInstance(ctx, id)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "virtual instance '%s' suspended\n", result.GetName())

			return nil
		},
	}

	return &cmd
}

func newResumeVirtualInstanceCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "virtualinstance ID|NAME",
		Aliases:           []string{"vi"},
		Short:             "resume virtual instance",
		Long:              "resume Rockset virtual instance",
		Args:              cobra.ExactArgs(1),
		Annotations:       group("virtual instance"),
		ValidArgsFunction: completion.VirtualInstance,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			id, err := viNameOrIDtoID(ctx, rs, args[0])
			if err != nil {
				return err
			}

			result, err := rs.ResumeVirtualInstance(ctx, id)
			if err != nil {
				return err
			}

			if err = waitUntilVIActive(rs, cmd, id); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "virtual instance '%s' resumed\n", result.GetName())

			return nil
		},
	}

	cmd.Flags().Bool(flag.Wait, false, "wait until virtual instance is ready")

	return &cmd
}

// TODO should this move to the Rockset go client instead?
func viNameOrIDtoID(ctx context.Context, rs *rockset.RockClient, nameOrID string) (string, error) {
	if !isUUID(nameOrID) {
		id, err := viNameToID(ctx, rs, nameOrID)
		if err != nil {
			return "", fmt.Errorf("failed to get virtual instance id for %s: %v", nameOrID, err)
		}

		return id, nil
	}

	return nameOrID, nil
}

var uuidRe = regexp.MustCompile(`[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4-[[:xdigit:]]{12`)

func isUUID(id string) bool {
	return uuidRe.MatchString(id)
}

func viNameToID(ctx context.Context, rs *rockset.RockClient, name string) (string, error) {
	vis, err := rs.ListVirtualInstances(ctx)
	if err != nil {
		return "", err
	}

	for _, vi := range vis {
		if vi.GetName() == name {
			return vi.GetId(), nil
		}
	}

	return "", VINotFoundErr
}

var VINotFoundErr = errors.New("virtual instance not found")

func waitUntilVIActive(rs *rockset.RockClient, cmd *cobra.Command, vID string) error {
	wait, err := cmd.Flags().GetBool(flag.Wait)
	if err != nil {
		return err
	}
	if wait {
		// TODO notify the user that we're waiting
		if err := rs.Wait.UntilVirtualInstanceActive(cmd.Context(), vID); err != nil {
			return fmt.Errorf("failed to wait for %s to be active: %v", vID, err)
		}
	}

	return nil
}
