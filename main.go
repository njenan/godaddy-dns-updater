package main

import "git.nathanjenan.me/njenan/godaddy-dns-updater/cmd"

var Version string
var License string
var FullLicense string

func main() {
	cmd.Version = Version
	cmd.License = License
	cmd.FullLicense = FullLicense
	cmd.Execute()
}
