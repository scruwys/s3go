package cmd

import (
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
                        Target:       input.target,
                        Source:       input.source,
                        Recursive:    flagRecursive,
                        RequestPayer: flagRequestPayer,
                        DryRun:       flagDryRun,
                        DeleteAfter:  input.deleteAfter,
                    })

                    if err != nil {
                        s3go.Echo("%v", err)
                    }

                    if  err == nil && !flagQuiet {
                        s3go.Echo("%s%s: %s", dryRunPrefix, input.desc, logMessage)
                    }
            }
        }
    }()

    return resultCh
}
