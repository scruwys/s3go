package cmd

import (
	"github.com/spf13/cobra"
)

const (
    defaultEndpointUrl = "s3.amazonaws.com"
)

// Flags
var Debug bool
var Endpoint string
var VerifySSL bool
var Profile string
var Region string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "s3go",
	Short: "golang cli for interacting with aws s3.",
}

func init() {
	RootCmd.PersistentFlags().BoolVar(
		&Debug,
		"debug",
		false,
		"Turn on debug logging.")

	RootCmd.PersistentFlags().StringVar(
		&Endpoint,
		"endpoint-url",
		defaultEndpointUrl,
		"Override command's default URL with the given URL.")

	RootCmd.PersistentFlags().BoolVar(
		&VerifySSL,
		"no-verify-ssl",
		false,
		"This option overrides the default behavior of verifying SSL certificates.")

	RootCmd.PersistentFlags().StringVar(
		&Profile,
		"profile",
		"",
		"Use a specific profile from your credential file.")

	RootCmd.PersistentFlags().StringVar(
		&Region,
		"region",
		"",
		"The region to use. Overrides config/env settings.")

}
