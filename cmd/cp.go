package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func copyCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the cp command.")
}

var copyCommand = &cobra.Command{
	Use:   "cp",
	Short: "Copies a local file or S3 object to another location locally or in S3.",
	Run: copyCommandHandler,
}

func init() {
	RootCmd.AddCommand(copyCommand)
}
