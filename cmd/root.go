package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	var cfgFile string
	root := &cobra.Command{
		Use:   "rockset",
		Short: "A cli for Rockset",
		Long: `The rockset cli is used to...

To use the CLI you need an API Key, which you initially have to create using the console:
https://console.rockset.com/apikeys

It should either be stored as an environment variable ROCKSET_APIKEY.

For more configuration options, see the 'rockset config' command.

  rockset create sample collection --wait --dataset movies movies
  rockset query "SELECT COUNT(*) FROM movies"

`,
	}

	cobra.OnInitialize(initConfig(cfgFile))

	var current string
	cfg, _ := loadConfig()
	if cfg != nil {
		current = fmt.Sprintf(" (\"%s\")", cfg.Current)
	}

	// any persistent flag defined here will be visible in all commands
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rockset.yaml)")
	root.PersistentFlags().Bool("debug", false, "enable debug output")
	root.PersistentFlags().String(FormatFlag, DefaultFormat, "output format")
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

			// Search config in home directory with name ".r7" (without extension).
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
