package cmd

import (
	"git.nathanjenan.me/njenan/godaddy-dns-updater/updater"
	"github.com/spf13/cobra"
)

var endpoint string
var recordNames []string
var dryRun bool

func init() {
	updateCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "API endpoint to use")
	updateCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry-run (report what will be updated)")
	updateCmd.Flags().StringSliceVarP(&recordNames, "record-names", "r", []string{}, "List of A record names to update (if not specified, all found A records will be updated)")

	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [DOMAIN] [IP]",
	Short: "Check and update A records for the specified domain",
	Long:  "Check and update A records for the specified domain to the specified IP address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		ip := args[1]

		updater := updater.Updater{}
		updater.CheckAndUpdate(domain, ip)
	},
}
