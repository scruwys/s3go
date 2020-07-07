package cmd

import (
    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var makeBucketCommand = &cobra.Command{
	Use:   "mb",
	Short: "Creates an S3 bucket.",
    Args:  cobra.ExactArgs(1),
	Run:   makeBucketCommandHandler,
}

func makeBucketCommandHandler(cmd *cobra.Command, args []string) {
    uri, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    client := newClientWithPersistentFlags()

    if err = client.MakeBucket(uri.Bucket, Region); err != nil {
        s3go.ExitWithError(1, err)
    }

    s3go.Echo("make_bucket: %s://%s", uri.Scheme, uri.Bucket)
}

func init() {
	RootCmd.AddCommand(makeBucketCommand)
}
