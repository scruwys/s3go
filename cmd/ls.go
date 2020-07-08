package cmd

import (
	"strings"
	"path/filepath"

    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

var listCommand = &cobra.Command{
	Use:   "ls",
	Short: "List S3 objects and common prefixes under a prefix or all S3 buckets.",
    Args:  cobra.ExactArgs(1),
	Run:   listCommandHandler,
}

func listCommandHandler(cmd *cobra.Command, args []string) {
    uri, err := s3go.ParseUrl(args[0])

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    client := newClientWithPersistentFlags()

    objectCt := 0
    objectSz := *new(int64)

    objectCh, err := client.ListObjectsV2(&s3go.ListObjectsV2Input{
        Bucket:    uri.Bucket,
        Prefix:    uri.Prefix,
        Recursive: flagRecursive,
    })

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    for object := range objectCh {
    	if object.IsPrefix {
    		prefix := *object.Key

    		if strings.Trim(prefix, "/") != strings.Trim(uri.Prefix, "/") {
    			prefix = strings.Replace(prefix, uri.Prefix, "", -1)
    		}

	        s3go.Echo("%28sPRE %s", "", prefix)
	        continue
    	}

    	ts := object.LastModified.Format("2006-01-02 15:04:05")
    	sz := s3go.ByteSizeToString(*object.Size, flagHumanReadable)
    	ok := *object.Key

    	if !flagRecursive {
    		ok = filepath.Base(*object.Key)
    	}

	    s3go.Echo("%s %11s %s", ts, sz, ok)

	    objectSz += *object.Size
	    objectCt += 1
    }

    if flagSummarize {
	    s3go.Echo("\nTotal Objects: %v", objectCt)
	    s3go.Echo("   Total Size: %s", s3go.ByteSizeToString(objectSz, flagHumanReadable))
    }
}

func init() {
	listCommand.Flags().BoolVar(
		&flagRecursive,
		"recursive",
		false,
		"Command is performed on all files or objects under the specified directory or prefix.")

	listCommand.Flags().BoolVar(
		&flagHumanReadable,
		"human-readable",
		false,
		"Displays file sizes in human readable format.")

	listCommand.Flags().BoolVar(
		&flagSummarize,
		"summarize",
		false,
		"Displays summary information (number of objects, total size).")

	listCommand.Flags().StringVar(
		&flagRequestPayer,
		"request-payer",
		"",
		"Confirms that the requester knows that she or he will be charged for the request.")

	RootCmd.AddCommand(listCommand)
}
