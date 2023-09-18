package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/rockset/rockset-go-client"
)

func newListConfigCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "configurations",
		Aliases:     []string{"configuration", "configs", "config", "cfg"},
		Annotations: group("config"),
		Args:        cobra.NoArgs,
		Short:       "list configurations",
		Long: `list configurations and show the currently selected

YAML file located in ~/.config/rockset/cli.yaml of the format 
---
current: dev
configs:
  dev:
    apikey: ...
    apiserver: api.usw2a1.rockset.com
  prod:
    apikey: ...
    apiserver: api.use1a1.rockset.com`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "available configs:\n")
			var names []string
			for name := range cfg.Configs {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				var arrow = "  "
				if cfg.Current == name {
					arrow = "->"
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s (%s)\n", arrow, name, cfg.Configs[name].APIServer)
			}

			return nil
		},
	}

	return &cmd
}

func newUpdateConfigCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "configuration",
		Aliases:     []string{"config", "cfg"},
		Short:       "configuration",
		Long:        "configuration command",
		Annotations: group("config"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			if _, found := cfg.Configs[args[0]]; !found {
				return fmt.Errorf("configuration %s not found", args[0])
			}

			cfg.Current = args[0]
			if err = storeConfig(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "using %s\n", args[0])

			return nil
		},
	}

	return &cmd
}

const (
	ConfigFile  = "cli.yaml"
	HistoryFile = "cli.hist"
)

type Configs struct {
	Current string            `yaml:"current"`
	Configs map[string]Config `yaml:"configs"`
}

type Config struct {
	APIKey    string `yaml:"apikey"`
	APIServer string `yaml:"apiserver"`
}

func configFile() (string, error) {
	return rocksetFile(ConfigFile)
}

func historyFile() (string, error) {
	return rocksetFile(HistoryFile)
}

func rocksetFile(name string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".config", "rockset", name), nil
}

func storeConfig(cfg *Configs) error {
	file, err := configFile()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(f)
	if err = enc.Encode(cfg); err != nil {
		return err
	}

	if err = enc.Close(); err != nil {
		logger.Error("failed to close config", "err", err)
	}

	return nil
}

func loadConfig() (*Configs, error) {
	cfg, err := configFile()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(cfg)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, nil
	}
	var contexts Configs
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&contexts)
	if err != nil {
		return nil, err
	}
	return &contexts, nil
}

func rockClient(cmd *cobra.Command) (*rockset.RockClient, error) {
	var apikey, apiserver string

	// load from config, ok if none is found
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	if cfg != nil {
		context, _ := cmd.Flags().GetString(ContextFLag)
		if context == "" {
			slog.Debug("using context from file")
			context = cfg.Current
		}
		if ctx, found := cfg.Configs[context]; found {
			slog.Debug("using", context, context)
			apikey = ctx.APIKey
			apiserver = ctx.APIServer
		} else {
			slog.Debug("not found", "context", context)
		}
	}

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

	return rockset.NewClient(
		rockset.WithAPIKey(apikey),
		rockset.WithAPIServer(apiserver),
		// TODO add rockset.WithCustomHeader("rockset-go-cli", Version) once the new client it released
	)
}
