package cmd

import (
	"github.com/nanoteck137/nosepass"
	"github.com/nanoteck137/nosepass/config"
	"github.com/nanoteck137/nosepass/core/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     nosepass.AppName,
	Version: nosepass.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(nosepass.VersionTemplate(nosepass.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
