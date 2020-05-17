package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/scruwys/s3go/internal"
)

func listCommandHandler(cmd *cobra.Command, args []string) {
	client := s3go.Client{}

	fmt.Println(
		fmt.Sprintf("ls command executed: %s", client.Get()),
	)
}

var listCommand = &cobra.Command{
	Use:   "ls",
	Short: "TBD",
	Run: listCommandHandler,
}

func init() {
	RootCmd.AddCommand(listCommand)
}
