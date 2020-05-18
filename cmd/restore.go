package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func restoreCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the restore command.")
}

var restoreCommand = &cobra.Command{
	Use:   "restore",
	Short: "Restores S3 object(s) stored in Glacier.",
	Run: restoreCommandHandler,
}

func init() {
	RootCmd.AddCommand(restoreCommand)
}
