package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func makeBucketCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the mb command.")
}

var makeBucketCommand = &cobra.Command{
	Use:   "mb",
	Short: "Creates an S3 bucket.",
	Run: makeBucketCommandHandler,
}

func init() {
	RootCmd.AddCommand(makeBucketCommand)
}
