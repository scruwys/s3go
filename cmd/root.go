package cmd

import (
	"fmt"

    "github.com/spf13/cobra"
    "github.com/scruwys/s3go/internal"
)

// Persistent Flags
var Debug bool
var Endpoint string
var VerifySSL bool
var Profile string
var Region string

// Local Flags
var flagACL string
var flagConcurrency int
var flagContentDisposition string
var flagContentEncoding string
var flagContentLanguage string
var flagContentType string
var flagDryRun bool
var flagExcludeFilter string
var flagExpiresIn int
var flagForce bool
var flagHumanReadable bool
var flagIncludeFilter string
var flagOnlyShowErrors bool
var flagQuiet bool
var flagRecursive bool
var flagRequestPayer string
var flagSummarize bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "s3go",
	Short: "golang cli for interacting with aws s3.",
}

// Make a new s3go.Client using the default persistent flags
func newClientWithPersistentFlags() *s3go.Client {
	if Endpoint == "" {
		Endpoint = s3go.DEFAULT_ENDPOINT_URL
	}

    return s3go.NewClient(&s3go.ClientOptions{
        Endpoint:   Endpoint,
        Debug:      Debug,
        Profile:    Profile,
        Region:     Region,
        DisableSSL: !VerifySSL,
    })
}

// Overrides some flags to pull the region from the bucket
func newClientWithRegionFromBucket(bucketName string) *s3go.Client {
    region, err := s3go.GetBucketRegion(bucketName)

    if err != nil {
        s3go.ExitWithError(1, err)
    }

    Region = region

    if Endpoint == "" {
		Endpoint = fmt.Sprintf("s3.%s.amazonaws.com", region)
    }

    return newClientWithPersistentFlags()
}

func init() {
	RootCmd.PersistentFlags().StringVar(
		&Endpoint,
		"endpoint-url",
		"",
		"Override command's default URL with the given URL.")

	RootCmd.PersistentFlags().BoolVar(
		&Debug,
		"debug",
		false,
		"Turn on debug logging.")

	RootCmd.PersistentFlags().StringVar(
		&Profile,
		"profile",
		"",
		"Use a specific profile from your credential file.")

	RootCmd.PersistentFlags().StringVar(
		&Region,
		"region",
		"us-east-1",
		"The region to use. Overrides config/env settings.")

	RootCmd.PersistentFlags().BoolVar(
		&VerifySSL,
		"no-verify-ssl",
		false,
		"This option overrides the default behavior of verifying SSL certificates.")
}
