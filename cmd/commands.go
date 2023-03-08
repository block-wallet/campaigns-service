package cmd

import (
	"github.com/block-wallet/golang-service-template/cmd/server"
	"github.com/spf13/cobra"
)

func Cmds(version string) *cobra.Command {
	rootCmd := &cobra.Command{Version: version}

	runnable := server.NewRunnable()
	runnableCmd := runnable.Cmd(version)

	rootCmd.AddCommand(runnableCmd)

	return rootCmd
}
