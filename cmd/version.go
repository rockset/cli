package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime/debug"
)

func newVersionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Rockset CLI version: %s\n", Version)

			if d, _ := cmd.Flags().GetBool(DebugFlag); !d {
				return nil
			}

			info, ok := debug.ReadBuildInfo()
			if !ok {
				logger.Error("no build info found")
				return nil
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\tgo version: %s\n", info.GoVersion)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\tpath: %s\n", info.Path)

			for _, dep := range info.Deps {
				if dep.Path == "github.com/rockset/rockset-go-client" {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\tRockset go client version: %s\n", dep.Version)
					break
				}
			}

			for _, s := range info.Settings {
				if s.Key == "vcs.revision" {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\tcommit hash: %s\n", s.Value)
				}
				if s.Key == "vcs.modified" {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\tdirty: %s\n", s.Value)
				}
				if s.Key == "vcs.time" {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\tbuild time: %s\n", s.Value)
				}
			}

			return nil
		},
	}

	return &cmd
}
