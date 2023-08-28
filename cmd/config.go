package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
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
			log.Printf("using context from file")
			context = cfg.Current
		}
		if ctx, found := cfg.Configs[context]; found {
			log.Printf("using context %s", context)
			apikey = ctx.APIKey
			apiserver = ctx.APIServer
		} else {
			log.Printf("context %s not found", context)
		}
	}

	// load from environment
	if key, found := os.LookupEnv(rockset.APIKeyEnvironmentVariableName); found {
		log.Printf("set apikey")
		apikey = key
	}
	if server, found := os.LookupEnv(rockset.APIServerEnvironmentVariableName); found {
		log.Printf("set apiserver")
		apiserver = server
	}

	// let the --cluster flag override the apiserver
	if cluster, _ := cmd.Flags().GetString(ClusterFLag); cluster != "" {
		log.Printf("override apiserver")
		apiserver = fmt.Sprintf("https://api.%s.rockset.com", cluster)
	}

	return rockset.NewClient(rockset.WithAPIKey(apikey), rockset.WithAPIServer(apiserver))
}
