package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
	"regexp"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/option"
)

func newCreateVirtualInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Short:       "create a virtual instance",
		Long:        "create a Rockset virtual instance",
		Args:        cobra.ExactArgs(1),
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			size, _ := cmd.Flags().GetString(SizeFlag)

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var options = []option.VirtualInstanceOption{
				option.WithVirtualInstanceSize(option.VirtualInstanceSize(size)),
			}

			// TODO uncomment after go client supports setting it
			//if desc, _ := cmd.Flags().GetString(DescriptionFlag); desc != "" {
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
	//cmd.Flags().StringP(DescriptionFlag, "d", "", "virtual instance description")
	cmd.Flags().Bool(WaitFlag, false, "wait until virtual instance is active")
	cmd.Flags().String(SizeFlag, "", "virtual instance size")
	_ = cobra.MarkFlagRequired(cmd.Flags(), SizeFlag)

	return cmd
}

func newUpdateVirtualInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Short:       "update a virtual instance",
		Long:        "update a Rockset virtual instance",
		Args:        cobra.ExactArgs(1),
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			size, _ := cmd.Flags().GetString(SizeFlag)

			rs, err := rockClient(cmd)
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
	//cmd.Flags().StringP(DescriptionFlag, "d", "", "virtual instance description")
	cmd.Flags().String(SizeFlag, "", "virtual instance size")
	cmd.Flags().Bool(WaitFlag, false, "wait until virtual instance is active")
	_ = cobra.MarkFlagRequired(cmd.Flags(), SizeFlag)

	return cmd
}

func newListVirtualInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstances",
		Aliases:     []string{"vi", "vis"},
		Args:        cobra.NoArgs,
		Short:       "list virtual instances",
		Long:        "list Rockset virtual instances",
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			list, err := rs.ListVirtualInstances(ctx)
			if err != nil {
				return err
			}

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}

	cmd.Flags().Bool(WideFlag, false, "display more information")

	return cmd
}

func newGetVirtualInstancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Args:        cobra.ExactArgs(1),
		Short:       "get virtual instance",
		Long:        "get Rockset virtual instances",
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
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

	cmd.Flags().Bool(WideFlag, false, "display more information")

	return cmd
}

func newDeleteVirtualInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Short:       "delete virtual instance",
		Long:        "delete Rockset virtual instance",
		Args:        cobra.ExactArgs(1),
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
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

	return cmd
}

func newSuspendVirtualInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Short:       "suspend virtual instance",
		Long:        "suspend Rockset virtual instance",
		Args:        cobra.ExactArgs(1),
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
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

	return cmd
}

func newResumeVirtualInstanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "virtualinstance ID|NAME",
		Aliases:     []string{"vi"},
		Short:       "resume virtual instance",
		Long:        "resume Rockset virtual instance",
		Args:        cobra.ExactArgs(1),
		Annotations: group("virtual instance"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			rs, err := rockClient(cmd)
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

	cmd.Flags().Bool(WaitFlag, false, "wait until virtual instance is ready")

	return cmd
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
	wait, err := cmd.Flags().GetBool(WaitFlag)
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
