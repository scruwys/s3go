package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func removeCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the rm command.")
}

var removeCommand = &cobra.Command{
	Use:   "rm",
	Short: "Deletes an S3 object.",
	Run: removeCommandHandler,
}

func init() {
	RootCmd.AddCommand(removeCommand)
}
