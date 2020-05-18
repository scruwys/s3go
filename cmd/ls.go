package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the ls command.")
}

var listCommand = &cobra.Command{
	Use:   "ls",
	Short: "List S3 objects and common prefixes under a prefix or all S3 buckets.",
	Run: listCommandHandler,
}

func init() {
	RootCmd.AddCommand(listCommand)
}
