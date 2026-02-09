package cmd

import (
	"fmt"
	"os"
	"time"

	checker "github.com/Arnel-Rah/uptimego/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the uptime monitoring daemon",
	Long: `Starts the background daemon that periodically checks configured HTTP endpoints,
collects uptime/latency metrics, and triggers alerts if needed.

The daemon runs indefinitely until stopped (Ctrl+C or SIGTERM).
Configuration is loaded from a YAML file (default: config.yaml in current dir).`,
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
			_, err := fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}

		fmt.Printf("Daemon starting...\n")
		fmt.Printf("Listening on port: %d\n", viper.GetInt("port"))

		rawEndpoints := viper.Get("endpoints")
		if rawEndpoints == nil {
			fmt.Println("Aucun endpoint configuré")
		} else {
			endpoints, ok := rawEndpoints.([]interface{})
			if !ok {
				fmt.Fprintf(os.Stderr, "Erreur: 'endpoints' n'est pas une liste (type trouvé: %T)\n", rawEndpoints)
				os.Exit(1)
			}

			for _, ep := range endpoints {
				endpoint, ok := ep.(map[string]interface{})
				if !ok {
					fmt.Println("Endpoint invalide, saut...")
					continue
				}

				name, _ := endpoint["name"].(string)
				url, _ := endpoint["url"].(string)
				timeoutStr, _ := endpoint["timeout"].(string)
				timeout, err := time.ParseDuration(timeoutStr)
				if err != nil {
					fmt.Printf("Timeout invalide pour %s: %v\n", name, err)
					continue
				}

				result := checker.CheckEndpoint(url, timeout)
				fmt.Println(checker.FormatResult(name, url, result))
			}
		}

		fmt.Println("Initial checks done. Monitoring loop coming soon...")
		fmt.Println("Monitoring started. Press Ctrl+C to stop.")
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		ticker = time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("--- New check cycle ---")
				endpoints, _ := rawEndpoints.([]interface{})
				for _, ep := range endpoints {
					endpoint, ok := ep.(map[string]interface{})
					if !ok {
						fmt.Println("Endpoint invalide, saut...")
						continue
					}

					name, _ := endpoint["name"].(string)
					url, _ := endpoint["url"].(string)
					timeoutStr, _ := endpoint["timeout"].(string)
					timeout, err := time.ParseDuration(timeoutStr)
					if err != nil {
						fmt.Printf("Timeout invalide pour %s: %v\n", name, err)
						continue
					}

					result := checker.CheckEndpoint(url, timeout)
					fmt.Println(checker.FormatResult(name, url, result))
				}

			case <-cmd.Context().Done():
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("config", "c", "", "path to config file (default: config.yaml)")
}
