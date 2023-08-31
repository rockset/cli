package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
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

// TODO change to HTTPS

const VersionURL = "http://rockset.sh/install/version.json"

type Versions struct {
	Stable string `json:"stable"`
	Beta   string `json:"beta"`
}

func VersionCheck(ctx context.Context, ch chan string) {
	// always send something on the channel
	defer func() { ch <- "" }()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, VersionURL, nil)
	if err != nil {
		logger.Error("failed to create http request", "err", err)
		return
	}

	c := http.Client{}
	response, err := c.Do(request)
	if err != nil {
		logger.Error("failed to perform http request", "err", err)
		return
	}

	var v Versions
	dec := json.NewDecoder(response.Body)
	if err = dec.Decode(&v); err != nil {
		logger.Error("failed to unmarshal json", "err", err)
		return
	}

	if v.Stable != Version {
		ch <- fmt.Sprintf("version %s is available", v.Stable)
	}
}
