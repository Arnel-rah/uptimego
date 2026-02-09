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

		// Chargement unique des endpoints
		rawEndpoints := viper.Get("endpoints")
		var endpoints []interface{}
		if rawEndpoints != nil {
			var ok bool
			endpoints, ok = rawEndpoints.([]interface{})
			if !ok {
				_, err := fmt.Fprintf(os.Stderr, "Erreur: 'endpoints' n'est pas une liste (type trouvé: %T)\n", rawEndpoints)
				if err != nil {
					return
				}
				os.Exit(1)
			}
		} else {
			fmt.Println("Aucun endpoint configuré")
		}

		// Fonction réutilisable pour checker un endpoint
		checkEndpoint := func(endpoint map[string]interface{}) {
			name, _ := endpoint["name"].(string)
			url, _ := endpoint["url"].(string)
			timeoutStr, _ := endpoint["timeout"].(string)

			timeout, err := time.ParseDuration(timeoutStr)
			if err != nil {
				fmt.Printf("Timeout invalide pour %s: %v\n", name, err)
				return
			}

			result := checker.CheckEndpoint(url, timeout)
			fmt.Println(checker.FormatResult(name, url, result))
		}

		// Checks initiaux
		fmt.Println("--- Initial checks ---")
		for _, ep := range endpoints {
			endpoint, ok := ep.(map[string]interface{})
			if !ok {
				fmt.Println("Endpoint invalide, saut...")
				continue
			}
			checkEndpoint(endpoint)
		}

		fmt.Println("Initial checks done. Monitoring loop starting...")
		fmt.Println("Monitoring started. Press Ctrl+C to stop.")

		// Boucle périodique
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("--- New check cycle ---")
				for _, ep := range endpoints {
					endpoint, ok := ep.(map[string]interface{})
					if !ok {
						fmt.Println("Endpoint invalide, saut...")
						continue
					}
					checkEndpoint(endpoint)
				}

			case <-cmd.Context().Done():
				fmt.Println("Daemon stopped gracefully")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("config", "c", "", "path to config file (default: config.yaml)")
}
