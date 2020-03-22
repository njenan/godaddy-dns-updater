package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check and update A records for the specified domain",
	Long:  "Check and update A records for the specified domain to the specified IP address",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
