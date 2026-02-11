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
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Daemon starting...\n")
		fmt.Printf("Listening on port: %d\n", viper.GetInt("port"))

		rawEndpoints := viper.Get("endpoints")
		var endpoints []interface{}
		if rawEndpoints != nil {
			var ok bool
			endpoints, ok = rawEndpoints.([]interface{})
			if !ok {
				fmt.Fprintf(os.Stderr, "Erreur: 'endpoints' n'est pas une liste (type trouvé: %T)\n", rawEndpoints)
				os.Exit(1)
			}
		} else {
			fmt.Println("Aucun endpoint configuré")
		}

		checkAndLogEndpoint := func(endpoint map[string]interface{}) {
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

		fmt.Println("--- Initial checks ---")
		for _, ep := range endpoints {
			endpoint, ok := ep.(map[string]interface{})
			if !ok {
				fmt.Println("Endpoint invalide, saut...")
				continue
			}
			checkAndLogEndpoint(endpoint)
		}

		fmt.Println("Initial checks done. Monitoring loop starting...")
		fmt.Println("Monitoring started. Press Ctrl+C to stop.")

		globalTickInterval := 15 * time.Second
		ticker := time.NewTicker(globalTickInterval)
		defer ticker.Stop()

		cycleCount := 0

		for {
			select {
			case <-ticker.C:
				cycleCount++
				fmt.Printf("--- Cycle %d (%s) ---\n", cycleCount, time.Now().Format("15:04:05"))

				for _, ep := range endpoints {
					endpoint, ok := ep.(map[string]interface{})
					if !ok {
						fmt.Println("Endpoint invalide, saut...")
						continue
					}

					name, _ := endpoint["name"].(string)
					intervalStr, _ := endpoint["interval"].(string)
					interval, err := time.ParseDuration(intervalStr)
					if err != nil {
						fmt.Printf("Interval invalide pour %s: %v\n", name, err)
						continue
					}

					expectedTicks := int(interval / globalTickInterval)
					if expectedTicks > 0 && cycleCount%expectedTicks == 0 {
						checkAndLogEndpoint(endpoint)
					}
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
