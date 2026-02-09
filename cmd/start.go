/*
Copyright © 2026 raharinandrasana <ton@email.com>  // ← Mets ton nom et email ici
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the uptime monitoring daemon",
	Long: `Starts the background daemon that periodically checks configured HTTP endpoints,
collects uptime/latency metrics, and triggers alerts if needed.

The daemon runs indefinitely until stopped (Ctrl+C or SIGTERM).
Configuration is loaded from a YAML file (default: config.yamlin current dir).`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			viper.AddConfigPath("$HOME/.uptimego")
		}
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}

		port := viper.GetInt("port")
		endpoints := viper.Get("endpoints")
		fmt.Printf("Daemon starting...\n")
		fmt.Printf("Listening on port: %d\n", port)
		fmt.Printf("Endpoints loaded: %v\n", endpoints)
		fmt.Println("Monitoring started. Press Ctrl+C to stop.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("config", "c", "", "path to config file (default: config.yaml)")
}
