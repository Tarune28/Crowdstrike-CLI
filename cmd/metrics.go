package cmd

import (
	"encoding/json"
	"fmt"

	"crowdstrike-cli/api"
	"crowdstrike-cli/config"
	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Get CrowdStrike metrics and summaries",
	Long:  `Retrieve various metrics and summaries from CrowdStrike Falcon.`,
}

var detectionsCmd = &cobra.Command{
	Use:   "detections",
	Short: "Get detection summary metrics",
	Long:  `Retrieve detection summary metrics from CrowdStrike.`,
	Run:   getDetections,
}

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Get host summary metrics",
	Long:  `Retrieve host summary metrics from CrowdStrike.`,
	Run:   getHosts,
}

var incidentsCmd = &cobra.Command{
	Use:   "incidents",
	Short: "Get incident summary metrics",
	Long:  `Retrieve incident summary metrics from CrowdStrike.`,
	Run:   getIncidents,
}

func getDetections(cmd *cobra.Command, args []string) {
	client, err := createAPIClient()
	if err != nil {
		fmt.Printf("Error creating API client: %v\n", err)
		return
	}

	data, err := client.GetDetectionSummary()
	if err != nil {
		fmt.Printf("Error getting detections: %v\n", err)
		return
	}

	printJSONResponse("Detection Summary", data)
}

func getHosts(cmd *cobra.Command, args []string) {
	client, err := createAPIClient()
	if err != nil {
		fmt.Printf("Error creating API client: %v\n", err)
		return
	}

	data, err := client.GetHostSummary()
	if err != nil {
		fmt.Printf("Error getting hosts: %v\n", err)
		return
	}

	printJSONResponse("Host Summary", data)
}

func getIncidents(cmd *cobra.Command, args []string) {
	client, err := createAPIClient()
	if err != nil {
		fmt.Printf("Error creating API client: %v\n", err)
		return
	}

	data, err := client.GetIncidentSummary()
	if err != nil {
		fmt.Printf("Error getting incidents: %v\n", err)
		return
	}

	printJSONResponse("Incident Summary", data)
}

// createAPIClient creates and returns a CrowdStrike API client
func createAPIClient() (*api.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("credentials not configured. Use 'cs creds set' to configure")
	}

	return api.NewClient(cfg), nil
}

// printJSONResponse pretty prints a JSON response
func printJSONResponse(title string, data []byte) {
	fmt.Printf("\n=== %s ===\n", title)
	
	// Pretty print the JSON
	var prettyJSON interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		fmt.Printf("Raw response: %s\n", string(data))
		return
	}

	prettyData, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		fmt.Printf("Raw response: %s\n", string(data))
		return
	}

	fmt.Println(string(prettyData))
}

func init() {
	// Add subcommands to metrics command
	metricsCmd.AddCommand(detectionsCmd)
	metricsCmd.AddCommand(hostsCmd)
	metricsCmd.AddCommand(incidentsCmd)
	
	// Add metrics command to root
	rootCmd.AddCommand(metricsCmd)
}
