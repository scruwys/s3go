package cmd

import (

    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var removeCommand = &cobra.Command{
    Use:   "rm",
    Short: "Deletes an S3 object.",
    Args:  cobra.ExactArgs(1),
    Run:   removeCommandHandler,
}

func removeCommandHandler(cmd *cobra.Command, args []string) {
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

    workerInput := &removeCommandWorkerInput{objectCh, doneCh}
    workers := make([]<-chan s3go.ObjectInfo, flagConcurrency)

    for i := 0; i < flagConcurrency; i++ {
        workers[i] = removeCommandWorker(client, workerInput)
    }

    for range s3go.MergeWaitWithObjectInfo(doneCh, workers...) {
        continue
    }
}

type removeCommandWorkerInput struct {
    // Channel
    objectCh <-chan s3go.ObjectInfo

    // Channel
    doneCh <-chan bool
}

func removeCommandWorker(client *s3go.Client, input *removeCommandWorkerInput) <-chan s3go.ObjectInfo {
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
                        err = client.DeleteObject(*item.Bucket, *item.Key, flagRequestPayer)

                        if err != nil {
                            s3go.Echo("%v", err)
                        }
                    }

                    if err == nil && !flagQuiet && !flagOnlyShowErrors {
                        s3go.Echo("%sdelete: s3://%s/%s", dryRunPrefix, *item.Bucket, *item.Key)
                    }
            }
        }
    }()
    return resultCh
}

func init() {
    removeCommand.Flags().BoolVar(
        &flagDryRun,
        "dryrun",
        false,
        "Displays the operations that would be performed using the specified command without actually running them.")

    removeCommand.Flags().BoolVar(
        &flagQuiet,
        "quiet",
        false,
        "Does not display the operations performed from the specified command.")

    removeCommand.Flags().BoolVar(
        &flagRecursive,
        "recursive",
        false,
        "Command is performed on all files or objects under the specified directory or prefix.")

    removeCommand.Flags().StringVar(
        &flagIncludeFilter,
        "include",
        "",
        "Don't exclude files or objects in the command that match the specified pattern.")

    removeCommand.Flags().StringVar(
        &flagExcludeFilter,
        "exclude",
        "",
        "Exclude all files or objects from the command that matches the specified pattern.")

    removeCommand.Flags().StringVar(
        &flagRequestPayer,
        "request-payer",
        "",
        "Confirms that the requester knows that she or he will be charged for the request.")

    removeCommand.Flags().BoolVar(
        &flagOnlyShowErrors,
        "only-show-errors",
        false,
        "Only errors and warnings are displayed. All other output is suppressed.")

    removeCommand.Flags().IntVar(
        &flagConcurrency,
        "concurrency",
        1,
        "Number of concurrent workers (e.g., goroutines) to spin up.")

    RootCmd.AddCommand(removeCommand)
}
