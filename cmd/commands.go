package cmd

import (
	"github.com/block-wallet/campaigns-service/cmd/srv"
	"github.com/spf13/cobra"
)

func Cmds(version string) *cobra.Command {
	rootCmd := &cobra.Command{Version: version}

	runnable := srv.NewRunnable()
	runnableCmd := runnable.Cmd(version)

	rootCmd.AddCommand(runnableCmd)

	return rootCmd
}
