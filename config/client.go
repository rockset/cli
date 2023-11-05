package config

import (
	"errors"
	"log/slog"

	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/flag"
)

func Client(cmd *cobra.Command, version string) (*rockset.RockClient, error) {
	// load from config, ok if none is found
	cfg, err := Load()
	if err != nil {
		if !errors.Is(err, NotFoundErr) {
			return nil, err
		}
	}

	override, _ := cmd.Flags().GetString(flag.Context)
	if override != "" {
		slog.Debug("using override", "name", override)
	}

	var options = []rockset.RockOption{
		rockset.WithUserAgent("rockset-go-cli/" + version),
	}

	opts, err := cfg.AsOptions(override)
	if err != nil {
		return nil, err
	}
	options = append(options, opts...)

	return rockset.NewClient(options...)
}
