package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/rockset/cli/config"
	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"
)

func newListConfigCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:         "configs",
		Aliases:     []string{"config", "cfg"},
		Annotations: group("config"),
		Args:        cobra.NoArgs,
		Short:       "list configurations",
		Long: fmt.Sprintf(`list configurations and show the currently selected context

YAML file located in %s of the format 
---
current: dev
configs:
  dev:
    apikey: ...
    apiserver: api.usw2a1.rockset.com
  prod:
    apikey: ...
    apiserver: api.use1a1.rockset.com`, config.FileName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("config file %s not readable available: %v", config.FileName, err)
				}
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "available contexts:\n")
			var names []string
			for name := range cfg.Keys {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				var arrow = "  "
				if cfg.Current == name {
					arrow = "->"
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s (%s)\n", arrow, name, cfg.Keys[name].Server)
			}

			return nil
		},
	}

	return &cmd
}

func newCreateConfigCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "config NAME",
		Aliases:     []string{"cfg"},
		Short:       "create configuration",
		Long:        "create new configuration command",
		Annotations: group("config"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, _ := cmd.Flags().GetString("apikey")
			server, _ := cmd.Flags().GetString("server")
			force, _ := cmd.Flags().GetBool(ForceFlag)

			if key == "" || server == "" {
				// TODO open up a form
				return fmt.Errorf("both --apikey and --server are required")
			}

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			// use --force to add anyway
			if _, found := cfg.Keys[args[0]]; found {
				if force {
					logger.Info("config already exist, adding anyway")
				} else {
					return fmt.Errorf("configuration %s already exists", args[0])
				}
			}

			if !force {
				rs, err := rockset.NewClient(rockset.WithAPIKey(key), rockset.WithAPIServer(server))
				if err != nil {
					return fmt.Errorf("failed to create Rockset client using the new credentials: %v", err)
				}
				org, err := rs.GetOrganization(cmd.Context())
				if err != nil {
					return fmt.Errorf("failed to get org info using the new credentials: %v", err)
				}
				logger.Info("connected to Rockset", "org", org.DisplayName)
			} else {
				logger.Warn("skipping auth check due to --force option")
			}

			cfg.Keys[args[0]] = config.APIKey{
				Key:    key,
				Server: server,
			}

			if err = config.Store(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "config %s created\n", args[0])

			return nil
		},
	}

	cmd.Flags().String("server", "", "api server name")
	cmd.Flags().String("apikey", "", "apikey")
	cmd.Flags().Bool(ForceFlag, false, "force add the config even if the name exists or the credentials can't be used connect to the API server")

	return &cmd
}

func newUseConfigCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "config NAME",
		Aliases:     []string{"cfg"},
		Short:       "use configuration",
		Long:        "configuration command",
		Annotations: group("config"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if err = cfg.Use(args[0]); err != nil {
				return err
			}

			if err = config.Store(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "using config %s\n", args[0])

			return nil
		},
	}

	return &cmd
}

func rockClient(cmd *cobra.Command) (*rockset.RockClient, error) {
	// load from config, ok if none is found
	cfg, err := config.Load()
	if err != nil {
		if !errors.Is(err, config.NotFoundErr) {
			return nil, err
		}
	}

	override, _ := cmd.Flags().GetString(ContextFLag)
	if override != "" {
		slog.Debug("using override", "name", override)
	}

	var options = []rockset.RockOption{
		rockset.WithCustomHeader("rockset-go-cli", Version),
	}

	opts, err := cfg.AsOptions(override)
	if err != nil {
		return nil, err
	}
	options = append(options, opts...)

	/*
		// load from environment
		if key, found := os.LookupEnv(rockset.APIKeyEnvironmentVariableName); found {
			slog.Debug("set apikey")
			apikey = key
		}
		if server, found := os.LookupEnv(rockset.APIServerEnvironmentVariableName); found {
			slog.Debug("set apiserver")
			apiserver = server
		}

		// let the --cluster flag override the apiserver
		if cluster, _ := cmd.Flags().GetString(ClusterFLag); cluster != "" {
			slog.Debug("override apiserver")
			apiserver = fmt.Sprintf("https://api.%s.rockset.com", cluster)
		}
	*/
	return rockset.NewClient(options...)
}
