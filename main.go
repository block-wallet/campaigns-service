package main

import (
	"os"

	"github.com/block-wallet/campaigns-service/cmd"
)

var Version string

func main() {
	if err := cmd.Cmds(Version).Execute(); err != nil {
		os.Exit(1)
	}
}
