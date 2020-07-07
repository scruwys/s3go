package cmd

import (
    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var copyCommand = &cobra.Command{
	Use:   "cp",
	Short: "Copies a local file or S3 object to another location locally or in S3.",
    Args:  cobra.ExactArgs(2),
	Run: copyCommandHandler,
}


func copyCommandHandler(cmd *cobra.Command, args []string) {
    source, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    target, err := s3go.ParseUrl(args[1])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    client := newClientWithPersistentFlags()
    doneCh := make(chan bool)

    objectCh, err := client.ListSourceObjects(&s3go.ListSourceObjectsInput{
        SourceUrl:     source,
        Recursive:     flagRecursive,
        IncludeFilter: flagIncludeFilter,
        ExcludeFilter: flagExcludeFilter,
    })

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    workerInput := &copyCommandWorkerInput{
        objectCh: objectCh,
        doneCh:   doneCh,
        source:   source,
        target:   target,
    }

    workers := make([]<-chan s3go.ObjectInfo, flagConcurrency)

    for i := 0; i < flagConcurrency; i++ {
        workers[i] = copyCommandWorker(client, workerInput)
    }

    for range s3go.MergeWaitWithObjectInfo(doneCh, workers...) {
        continue
    }
}

type copyCommandWorkerInput struct {
    // Channel
    objectCh <-chan s3go.ObjectInfo

    // Channel
    doneCh <-chan bool

    // Docs
    source *s3go.S3Url

    // Docs
    target *s3go.S3Url
}

func copyCommandWorker(client *s3go.Client, input *copyCommandWorkerInput) <-chan s3go.ObjectInfo {
    resultCh := make(chan s3go.ObjectInfo)

    go func() {
        defer close(resultCh)
        for item := range input.objectCh {
            select {
                case <-input.doneCh:
                    return
                case resultCh <- item:
                    if !item.IsPrefix {
                        client.MoveObject(item, input.source, input.target, true)
                    }
            }
        }
    }()

    return resultCh
}

func init() {
    copyCommand.Flags().BoolVar(
        &flagDryRun,
        "dryrun",
        false,
        "Displays the operations that would be performed using the specified command without actually running them.")

    copyCommand.Flags().BoolVar(
        &flagQuiet,
        "quiet",
        false,
        "Does not display the operations performed from the specified command.")

    copyCommand.Flags().BoolVar(
        &flagRecursive,
        "recursive",
        false,
        "Command is performed on all files or objects under the specified directory or prefix.")

    copyCommand.Flags().StringVar(
        &flagIncludeFilter,
        "include",
        "",
        "Don't exclude files or objects in the command that match the specified pattern.")

    copyCommand.Flags().StringVar(
        &flagExcludeFilter,
        "exclude",
        "",
        "Exclude all files or objects from the command that matches the specified pattern.")

    copyCommand.Flags().StringVar(
        &flagRequestPayer,
        "request-payer",
        "",
        "Confirms that the requester knows that she or he will be charged for the request.")

    copyCommand.Flags().IntVar(
        &flagConcurrency,
        "concurrency",
        1,
        "Number of concurrent workers (e.g., goroutines) to spin up.")

	RootCmd.AddCommand(copyCommand)
}
