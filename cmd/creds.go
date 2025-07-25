package cmd

import (
	"fmt"
	"syscall"

	"crowdstrike-cli/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var credsCmd = &cobra.Command{
	Use:   "creds",
	Short: "Manage CrowdStrike API credentials",
	Long:  `Configure and manage your CrowdStrike API credentials for authentication.`,
}

var setCredsCmd = &cobra.Command{
	Use:   "set",
	Short: "Set CrowdStrike API credentials",
	Long:  `Set your CrowdStrike API client ID and secret for authentication.`,
	Run:   setCredentials,
}

var showCredsCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current credentials (client ID only)",
	Long:  `Display the currently configured client ID (secret will be hidden for security).`,
	Run:   showCredentials,
}

func setCredentials(cmd *cobra.Command, args []string) {
	fmt.Print("Enter Client ID: ")
	var clientID string
	fmt.Scanln(&clientID)

	fmt.Print("Enter Client Secret: ")
	secretBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		return
	}
	fmt.Println() // New line after password input

	clientSecret := string(secretBytes)

	// Load existing config or create new one
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Update credentials
	cfg.ClientID = clientID
	cfg.ClientSecret = clientSecret

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Println("Credentials saved successfully!")
}

func showCredentials(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if cfg.ClientID == "" {
		fmt.Println("No credentials configured. Use 'cs creds set' to configure.")
		return
	}

	fmt.Printf("Client ID: %s\n", cfg.ClientID)
	fmt.Printf("Base URL: %s\n", cfg.BaseURL)
	fmt.Println("Client Secret: [HIDDEN]")
}

func init() {
	// Add subcommands to creds command
	credsCmd.AddCommand(setCredsCmd)
	credsCmd.AddCommand(showCredsCmd)
	
	// Add creds command to root
	rootCmd.AddCommand(credsCmd)
}
