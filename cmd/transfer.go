package cmd

import (
    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

func transferCommandHandler(args []string, desc string, deleteAfter bool) {
    source, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    target, err := s3go.ParseUrl(args[1])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    bucket := target.Bucket

    if target.IsLocal() {
        bucket = source.Bucket
    }

    client := newClientWithRegionFromBucket(bucket)
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

    workerInput := &transferCommandWorkerInput{objectCh, doneCh, source, target, deleteAfter, desc}
    workers := make([]<-chan s3go.ObjectInfo, flagConcurrency)

    for i := 0; i < flagConcurrency; i++ {
        workers[i] = transferCommandWorker(client, workerInput)
    }

    for range s3go.MergeWaitWithObjectInfo(doneCh, workers...) {
        continue
    }
}

type transferCommandWorkerInput struct {
    // Channel
    objectCh <-chan s3go.ObjectInfo

    // Channel
    doneCh <-chan bool

    // Docs
    source *s3go.S3Url

    // Docs
    target *s3go.S3Url

    deleteAfter bool

    desc string
}

func transferCommandWorker(client *s3go.Client, input *transferCommandWorkerInput) <-chan s3go.ObjectInfo {
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
                    if item.IsPrefix {
                        continue
                    }

                    logMessage, err := client.MoveObject(item, &s3go.MoveObjectOptions{
                        Target:             input.target,
                        Source:             input.source,
                        DeleteAfter:        input.deleteAfter,
                        ACL:                flagACL,
                        ContentDisposition: flagContentDisposition,
                        ContentEncoding:    flagContentEncoding,
                        ContentLanguage:    flagContentLanguage,
                        ContentType:        flagContentType,
                        DryRun:             flagDryRun,
                        Recursive:          flagRecursive,
                        RequestPayer:       flagRequestPayer,
                    })

                    if err != nil {
                        s3go.Echo("%v", err)
                    }

                    if  err == nil && !flagQuiet && !flagOnlyShowErrors {
                        s3go.Echo("%s%s: %s", dryRunPrefix, input.desc, logMessage)
                    }
            }
        }
    }()

    return resultCh
}

func addTransferFlags(command *cobra.Command) {
    command.Flags().BoolVar(
        &flagDryRun,
        "dryrun",
        false,
        "Displays the operations that would be performed using the specified command without actually running them.")

    command.Flags().BoolVar(
        &flagQuiet,
        "quiet",
        false,
        "Does not display the operations performed from the specified command.")

    command.Flags().BoolVar(
        &flagRecursive,
        "recursive",
        false,
        "Command is performed on all files or objects under the specified directory or prefix.")

    command.Flags().StringVar(
        &flagIncludeFilter,
        "include",
        "",
        "Don't exclude files or objects in the command that match the specified pattern.")

    command.Flags().StringVar(
        &flagExcludeFilter,
        "exclude",
        "",
        "Exclude all files or objects from the command that matches the specified pattern.")

    command.Flags().StringVar(
        &flagRequestPayer,
        "request-payer",
        "",
        "Confirms that the requester knows that she or he will be charged for the request.")

    command.Flags().StringVar(
        &flagACL,
        "acl",
        "",
        "Sets the ACL for the object when the command is performed.")

    command.Flags().BoolVar(
        &flagOnlyShowErrors,
        "only-show-errors",
        false,
        "Only errors and warnings are displayed. All other output is suppressed.")

    command.Flags().StringVar(
        &flagContentDisposition,
        "content-disposition",
        "",
        "Specifies presentational information for the object.")

    command.Flags().StringVar(
        &flagContentEncoding,
        "content-encoding",
        "",
        "Specifies what content encodings have been applied to the object.")

    command.Flags().StringVar(
        &flagContentLanguage,
        "content-language",
        "",
        "The language the content is in.")

    command.Flags().StringVar(
        &flagContentType,
        "content-type",
        "",
        "Specify an explicit content type for this operation. This value overrides any guessed mime types.")

    command.Flags().IntVar(
        &flagConcurrency,
        "concurrency",
        1,
        "Number of concurrent workers (e.g., goroutines) to spin up.")
}
