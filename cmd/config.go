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
    apiserver: api.use1a1.rockset.com`, APIKeysFile),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadAPIKeys()
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("config file %s not readable available: %v", APIKeysFile, err)
				}
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "available configs:\n")
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

			cfg, err := loadAPIKeys()
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

			cfg.Keys[args[0]] = APIKey{
				Key:    key,
				Server: server,
			}

			if err = storeAPIKeys(cfg); err != nil {
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
			cfg, err := loadAPIKeys()
			if err != nil {
				return err
			}

			if _, found := cfg.Keys[args[0]]; !found {
				return fmt.Errorf("configuration %s not found", args[0])
			}

			cfg.Current = args[0]
			if err = storeAPIKeys(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "using config %s\n", args[0])

			return nil
		},
	}

	return &cmd
}

const (
	APIKeysFileName = "apikeys.yaml"
	HistoryFile     = "cli.hist"
)

var APIKeysFile string

func init() {
	file, err := apikeysPath()
	if err != nil {
		panic(fmt.Sprintf("unable to locate apikey file %s: %v", APIKeysFileName, err))
	}
	APIKeysFile = file
}

type APIKeys struct {
	Current string            `yaml:"current"`
	Keys    map[string]APIKey `yaml:"configs"`
}

type APIKey struct {
	Key    string `yaml:"apikey"`
	Server string `yaml:"apiserver"`
}

func apikeysPath() (string, error) {
	return rocksetFile(APIKeysFileName)
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

func storeAPIKeys(cfg APIKeys) error {
	file, err := apikeysPath()
	if err != nil {
		return err
	}

	// TODO create directory

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

func loadAPIKeys() (APIKeys, error) {
	var keys APIKeys

	cfg, err := apikeysPath()
	if err != nil {
		return keys, err
	}

	f, err := os.Open(cfg)
	if err != nil {
		return keys, fmt.Errorf("failed to read apikey config file: %w", err)
	}

	dec := yaml.NewDecoder(f)
	err = dec.Decode(&keys)
	if err != nil {
		return keys, err
	}
	return keys, nil
}

func rockClient(cmd *cobra.Command) (*rockset.RockClient, error) {
	var apikey, apiserver string

	// load from config, ok if none is found
	cfg, err := loadAPIKeys()
	if err != nil {
		return nil, err
	}

	context, _ := cmd.Flags().GetString(ContextFLag)
	if context == "" {
		slog.Debug("using context from file")
		context = cfg.Current
	}
	if ctx, found := cfg.Keys[context]; found {
		slog.Debug("using", context, context)
		apikey = ctx.Key
		apiserver = ctx.Server
	} else {
		slog.Debug("not found", "context", context)
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
