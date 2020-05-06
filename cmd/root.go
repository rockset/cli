package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	var cfgFile string
	root := &cobra.Command{
		Use:   "rock",
		Short: "A cli for Rockset",
	}

	cobra.OnInitialize(initConfig(cfgFile))

	// any persistent flag defined here will be visible in all commands
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rockset.yaml)")
	root.PersistentFlags().Bool("debug", false, "enable debug output")
	root.PersistentFlags().String("profile", "", "configuration profile")

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
