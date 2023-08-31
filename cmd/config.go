package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"

	"github.com/rockset/rockset-go-client"
)

func newConfigCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "configuration",
		Aliases: []string{"config", "cfg"},
		Short:   "configuration",
		Long:    "configuration command",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "selected config is '%s'\n", cfg.Current)

			return nil
		},
	}

	return &c
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

	return rockset.NewClient(rockset.WithAPIKey(apikey), rockset.WithAPIServer(apiserver))
}
