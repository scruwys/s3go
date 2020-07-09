package cmd

import (
    "github.com/spf13/cobra"
)

var moveCommand = &cobra.Command{
	Use:   "mv",
	Short: "Moves a local file or S3 object to another location locally or in S3.",
    Args:  cobra.ExactArgs(2),
	Run: moveCommandHandler,
}

func moveCommandHandler(cmd *cobra.Command, args []string) {
    transferCommandHandler(args, "move", true)
}

func init() {
    addTransferFlags(moveCommand)

	RootCmd.AddCommand(moveCommand)
}
