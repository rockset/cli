package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/rockset/rockset-go-client"
	"github.com/spf13/cobra"

	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
)

// config is the file containing the auth contexts

func newListContextsCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:         "contexts",
		Aliases:     []string{"context", "ctx"},
		Annotations: group("context"),
		Args:        cobra.NoArgs,
		Short:       "list authentication contexts",
		Long:        fmt.Sprintf(`list authentication contexts and show the currently selected one. YAML file located in %s of the format`, config.FileName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("config file %s not readable available: %v", config.FileName, err)
				}
				return err
			}

			// TODO redo as tui
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Available Authentication Contexts:\n")
			if len(cfg.Keys) > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "apikeys:\n")
				listContext(cmd.OutOrStdout(), cfg.Current, cfg.Keys)
			}
			if len(cfg.Tokens) > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "bearer tokens:\n")
				listContext(cmd.OutOrStdout(), cfg.Current, cfg.Tokens)
			}
			return nil
		},
	}

	return &cmd
}

type apiserver interface {
	APIServer() string
}

func listContext[T apiserver](out io.Writer, current string, m map[string]T) {
	var names []string
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		var arrow = "  "
		if current == name {
			arrow = "->"
		}
		_, _ = fmt.Fprintf(out, "%s %s (%s)\n", arrow, name, m[name].APIServer())
	}
}

func newDeleteContextCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "context NAME",
		Aliases:     []string{"ctx"},
		Short:       "delete an authentication context",
		Annotations: group("context"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if err = cfg.DeleteContext(args[0]); err != nil {
				return err
			}

			if err = config.Store(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "context %s deleted\n", args[0])

			return nil
		},
	}

	return &cmd
}

func newCreateContextCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "context NAME",
		Aliases:     []string{"ctx"},
		Short:       "create a new authentication context",
		Long:        "create a new authentication context and save it in the config file",
		Annotations: group("context"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, _ := cmd.Flags().GetString("apikey")
			server, _ := cmd.Flags().GetString("server")
			force, _ := cmd.Flags().GetBool(flag.Force)

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
					logger.Info("context already exist, adding anyway")
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
	cmd.Flags().Bool(flag.Force, false, "force add the context even if the name exists or the credentials can't be used connect to the API server")

	return &cmd
}

func newUseContextCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "context NAME",
		Aliases:     []string{"ctx"},
		Short:       "use auth context",
		Long:        "use authentication context",
		Annotations: group("context"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			// TODO if len(args) == it should open a tui and let the user select
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
