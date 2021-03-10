package cmd

import (
	"github.com/pawmart/wp-atrd-task/api"
	"github.com/pawmart/wp-atrd-task/service"
	"github.com/spf13/cobra"
)

var config service.Config
var configName string

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Simple HTTP service",
	Long: `HTTP server implementing simple user-secrets API.
  	Its purpose is for demonstrating simple GO software development
  	for Wirtualna Polska Media SA recruitment process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.Unmarshal(configName)
		if err != nil {
			return err
		}
		return api.NewApi(
			service.NewRedisSecret(config.Redis),
		).Run()
	},
}

func Execute() (err error) {
	rootCmd.PersistentFlags().StringVar(
		&configName,
		"config",
		"config",
		"name of the config file to read on (without extension)",
	)
	rootCmd.PersistentFlags().StringP(
		service.OptionRedisAddress,
		"r",
		"wp-atrd-task-database:6379",
		"address for the redis database",
	)

	err = config.Init()
	if err != nil {
		return err
	}

	err = config.BindPFlag(service.OptionRedisAddress, rootCmd.PersistentFlags().Lookup(service.OptionRedisAddress))
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}
