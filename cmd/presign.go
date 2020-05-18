package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func presignCommandHandler(cmd *cobra.Command, args []string) {
	fmt.Printf("Executing the presign command.")
}

var presignCommand = &cobra.Command{
	Use:   "presign",
	Short: "Generate a pre-signed URL for an Amazon S3 object.",
	Run: presignCommandHandler,
}

func init() {
	RootCmd.AddCommand(presignCommand)
}
