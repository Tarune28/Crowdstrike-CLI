package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cs",
	Short: "A CLI tool for CrowdStrike API interactions",
	Long: `CrowdStrike CLI is a command-line tool that allows you to interact with 
CrowdStrike Falcon APIs to retrieve metrics, manage endpoints, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to CrowdStrike CLI!")
		fmt.Println("Use 'cs --help' to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
