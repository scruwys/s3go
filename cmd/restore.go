package cmd

import (
    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var restoreCommand = &cobra.Command{
	Use:   "restore",
	Short: "Restores S3 object(s) stored in Glacier.",
    Args:  cobra.ExactArgs(1),
	Run: restoreCommandHandler,
}

func restoreCommandHandler(cmd *cobra.Command, args []string) {
    uri, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    client := newClientWithRegionFromBucket(uri.Bucket)
    doneCh := make(chan bool)

    objectCh, err := client.ListObjectsV2(&s3go.ListObjectsV2Input{
        Bucket:        uri.Bucket,
        Prefix:        uri.Prefix,
        Recursive:     flagRecursive,
        IncludeFilter: flagIncludeFilter,
        ExcludeFilter: flagExcludeFilter,
    })

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    workerInput := &restoreCommandWorkerInput{objectCh, doneCh}
    workers := make([]<-chan s3go.ObjectInfo, flagConcurrency)

    for i := 0; i < flagConcurrency; i++ {
        workers[i] = restoreCommandWorker(client, workerInput)
    }

    for range s3go.MergeWaitWithObjectInfo(doneCh, workers...) {
        continue
    }
}

type restoreCommandWorkerInput struct {
    // Channel
    objectCh <-chan s3go.ObjectInfo

    // Channel
    doneCh <-chan bool
}

func restoreCommandWorker(client *s3go.Client, input *restoreCommandWorkerInput) <-chan s3go.ObjectInfo {
    resultCh := make(chan s3go.ObjectInfo)

    // We append this to output when we are doing a dry run.
    dryRunPrefix := ""

    if flagDryRun {
        dryRunPrefix = "(dryrun) "
    }

    go func() {
        defer close(resultCh)
        for item := range input.objectCh {
            select {
                case <-input.doneCh:
                    return
                case resultCh <- item:
                    err := error(nil)

                    if !flagDryRun {
                        err = client.RestoreObject(*item.Bucket, *item.Key, flagRequestPayer)

                        if err != nil {
                            s3go.Echo("%v", err)
                        }
                    }

                    if err == nil && !flagQuiet {
                        s3go.Echo("%srestore: s3://%s/%s", dryRunPrefix, *item.Bucket, *item.Key)
                    }
            }
        }
    }()
    return resultCh
}

func init() {
    restoreCommand.Flags().BoolVar(
        &flagDryRun,
        "dryrun",
        false,
        "Displays the operations that would be performed using the specified command without actually running them.")

    restoreCommand.Flags().BoolVar(
        &flagQuiet,
        "quiet",
        false,
        "Does not display the operations performed from the specified command.")

    restoreCommand.Flags().BoolVar(
        &flagRecursive,
        "recursive",
        false,
        "Command is performed on all files or objects under the specified directory or prefix.")

    restoreCommand.Flags().StringVar(
        &flagIncludeFilter,
        "include",
        "",
        "Don't exclude files or objects in the command that match the specified pattern.")

    restoreCommand.Flags().StringVar(
        &flagExcludeFilter,
        "exclude",
        "",
        "Exclude all files or objects from the command that matches the specified pattern.")

    restoreCommand.Flags().StringVar(
        &flagRequestPayer,
        "request-payer",
        "",
        "Confirms that the requester knows that she or he will be charged for the request.")

    restoreCommand.Flags().IntVar(
        &flagConcurrency,
        "concurrency",
        1,
        "Number of concurrent workers (e.g., goroutines) to spin up.")

	RootCmd.AddCommand(restoreCommand)
}
