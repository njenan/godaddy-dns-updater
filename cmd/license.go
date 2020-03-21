package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

var License string
var FullLicense string
var printFull bool

func init() {
	licenseCmd.Flags().BoolVarP(&printFull, "print-full-licence", "f", false, "Print the full text of the license? Default: false")

	rootCmd.AddCommand(licenseCmd)
}

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Print the license terms of GoDaddy DNS Updater",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !printFull {
			fmt.Println(License)
		} else {
			statikFS, err := fs.New()
			if err != nil {
				log.Fatal(err)
			}

			f, err := statikFS.Open("/LICENSE")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			license, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(license)

		}
	},
}
