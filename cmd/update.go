package cmd

import (
	"encoding/json"
	"fmt"

	"git.nathanjenan.me/njenan/godaddy-dns-updater/updater"
	"github.com/spf13/cobra"
)

var endpoint string
var authKey string
var authSecret string
var recordNames []string
var dryRun bool

func init() {
	updateCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "https://api.godaddy.com", "API endpoint to use")
	updateCmd.Flags().StringVarP(&authKey, "auth-key", "a", "", "Auth key to use when authenticating")
	updateCmd.Flags().StringVarP(&authSecret, "auth-secret", "s", "", "Auth secret to use when authenticating")
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

		var withRecordName []updater.WithRecordName
		for _, v := range recordNames {
			withRecordName = append(withRecordName, updater.WithRecordName(v))
		}

		// TODO pass in record names

		updateClient := updater.Updater{}
		records, err := updateClient.CheckAndUpdate(domain, ip, updater.WithEndpoint(endpoint), updater.WithDryRun(dryRun), updater.WithAuthKey(authKey), updater.WithAuthSecret(authSecret))
		if err != nil {
			fmt.Println(err)
			return
		}

		b, err := json.Marshal(records)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	},
}
