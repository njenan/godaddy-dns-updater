package main

import "git.nathanjenan.me/njenan/godaddy-dns-updater/cmd"

var Version string

func main() {
	cmd.Version = Version
	cmd.Execute()
}
