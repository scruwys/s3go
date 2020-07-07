package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var presignCommand = &cobra.Command{
    Use:   "presign",
    Short: "Generate a pre-signed URL for an Amazon S3 object.",
    Args:  cobra.ExactArgs(1),
    Run:   presignCommandHandler,
}

func presignCommandHandler(cmd *cobra.Command, args []string) {
    uri, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    client := newClientWithPersistentFlags()

    urlStr, err := client.Presign(uri.Bucket, uri.Prefix, flagExpiresIn)

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    fmt.Println(urlStr)
}

func init() {
    presignCommand.Flags().IntVar(
        &flagExpiresIn,
        "expires-in",
        3600,
        "Number of seconds until the pre-signed URL expires.")

    RootCmd.AddCommand(presignCommand)
}
