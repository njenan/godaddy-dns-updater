package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godaddy-dns-updater",
	Short: "Keep Godaddy dns updated from the command line",
	Long:  `A CLI designed to keep A records in GoDaddy updated quickly and easily.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
