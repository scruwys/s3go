package cmd

import (
    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var moveCommand = &cobra.Command{
	Use:   "mv",
	Short: "Moves a local file or S3 object to another location locally or in S3.",
    Args:  cobra.ExactArgs(2),
	Run: moveCommandHandler,
}

func moveCommandHandler(cmd *cobra.Command, args []string) {
	s3go.Echo("Executing the mv command.")
}

func init() {
	RootCmd.AddCommand(moveCommand)
}
