package cmd

import (
    "github.com/spf13/cobra"
)

var copyCommand = &cobra.Command{
	Use:   "cp",
	Short: "Copies a local file or S3 object to another location locally or in S3.",
    Args:  cobra.ExactArgs(2),
	Run: copyCommandHandler,
}


func copyCommandHandler(cmd *cobra.Command, args []string) {
    transferCommandHandler(args, "copy", false)
}

func init() {
    addTransferFlags(copyCommand)

	RootCmd.AddCommand(copyCommand)
}
