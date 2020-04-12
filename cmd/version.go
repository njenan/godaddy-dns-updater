package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GoDaddy DNS Updater",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
