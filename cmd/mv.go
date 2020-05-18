package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func moveCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the mv command.")
}

var moveCommand = &cobra.Command{
	Use:   "mv",
	Short: "Moves a local file or S3 object to another location locally or in S3.",
	Run: moveCommandHandler,
}

func init() {
	RootCmd.AddCommand(moveCommand)
}
