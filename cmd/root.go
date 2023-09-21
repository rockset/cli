package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rockset/cli/format"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var (
	R       = color.CyanString("[") + color.MagentaString("R") + color.CyanString("]")
	Rockset = color.CyanString("[") + color.MagentaString("Rockset") + color.CyanString("]")
	Version = "development"
	logger  = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
)

func NewRootCmd(version string) *cobra.Command {
	var cfgFile string
	Version = version

	root := &cobra.Command{
		Use:   "rockset",
		Short: "A cli for Rockset",
		Long: fmt.Sprintf(`The Rockset cli is used as a companion to the console. 

To use the CLI you need an API Key, which you initially have to create using the console:
https://console.rockset.com/apikeys

It should either be stored as an environment variable ROCKSET_APIKEY or in a
platform dependent configuration file, %s on the current computer.

For more configuration options, see the 'rockset create config' command.`, APIKeysFile),
		Example: `	## Create a sample collection and run a query against it
	rockset create sample collection --wait --dataset movies movies
	rockset query "SELECT COUNT(*) FROM movies"`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if debug, _ := cmd.Flags().GetBool(DebugFlag); debug {
				logger = slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
				slog.SetDefault(logger)
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig(cfgFile))

	var current string
	cfg, err := loadAPIKeys()
	if err != nil {
		current = fmt.Sprintf(" (\"%s\")", cfg.Current)
	}

	// any persistent flag defined here will be visible in all commands
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rockset.yaml)")
	root.PersistentFlags().Bool(DebugFlag, false, "enable debug output")

	root.PersistentFlags().String(FormatFlag, DefaultFormat, fmt.Sprintf("output format (%s)",
		strings.Join(format.SupportedFormats.ToStringArray(), ", ")))
	root.PersistentFlags().Bool(HeaderFlag, true, "show header")
	root.PersistentFlags().Bool(WideFlag, false, "show extended fields")

	root.PersistentFlags().String(ContextFLag, "", fmt.Sprintf("override currently selected configuration context%s", current))
	root.PersistentFlags().String(ClusterFLag, "", "override Rockset cluster")

	// this binds the environment variable DEBUG to the flag debug
	_ = viper.BindPFlag("debug", root.PersistentFlags().Lookup("debug"))

	addVerbs(root)
	return root
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cfgFile string) func() {
	return func() {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}

			// Search config in home directory with name ".rocket" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".rockset")
		}

		viper.SetEnvPrefix("rockset")
		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
