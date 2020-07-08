package cmd

import (

    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var removeBucketCommand = &cobra.Command{
	Use:   "rb",
	Short: "Deletes an empty S3 bucket.",
    Args:  cobra.ExactArgs(1),
	Run:   removeBucketCommandHandler,
}

func removeBucketCommandHandler(cmd *cobra.Command, args []string) {
    uri, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    client := newClientWithRegionFromBucket(uri.Bucket)

    if err = client.RemoveBucket(uri.Bucket, flagForce); err != nil {
        s3go.ExitWithError(1, err)
    }

    s3go.Echo("remove_bucket: %s://%s", uri.Scheme, uri.Bucket)
}

func init() {
	removeBucketCommand.Flags().BoolVar(
		&flagForce,
		"force",
		false,
		"Deletes all objects in the bucket including the bucket itself. Does not delete versions of objects.")

	RootCmd.AddCommand(removeBucketCommand)
}
