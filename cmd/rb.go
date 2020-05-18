package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func removeBucketCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the copy command.")
}

var removeBucketCommand = &cobra.Command{
	Use:   "rb",
	Short: "Deletes an empty S3 bucket.",
	Run: removeBucketCommandHandler,
}

func init() {
	RootCmd.AddCommand(removeBucketCommand)
}
