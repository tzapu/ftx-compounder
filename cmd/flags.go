package cmd

import "github.com/spf13/cobra"

func addRunFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("ftx.key", "xxx", "ftx key")
	cmd.PersistentFlags().String("ftx.secret", "xxx", "ftx secret")
}
