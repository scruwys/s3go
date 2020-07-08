package cmd

import (
    "github.com/spf13/cobra"
)

var moveCommand = &cobra.Command{
	Use:   "mv",
	Short: "Moves a local file or S3 object to another location locally or in S3.",
    Args:  cobra.ExactArgs(2),
	Run: moveCommandHandler,
}

func moveCommandHandler(cmd *cobra.Command, args []string) {
    transferCommandHandler(args, "move", true)
}

func init() {
    moveCommand.Flags().BoolVar(
        &flagDryRun,
        "dryrun",
        false,
        "Displays the operations that would be performed using the specified command without actually running them.")

    moveCommand.Flags().BoolVar(
        &flagQuiet,
        "quiet",
        false,
        "Does not display the operations performed from the specified command.")

    moveCommand.Flags().BoolVar(
        &flagRecursive,
        "recursive",
        false,
        "Command is performed on all files or objects under the specified directory or prefix.")

    moveCommand.Flags().StringVar(
        &flagIncludeFilter,
        "include",
        "",
        "Don't exclude files or objects in the command that match the specified pattern.")

    moveCommand.Flags().StringVar(
        &flagExcludeFilter,
        "exclude",
        "",
        "Exclude all files or objects from the command that matches the specified pattern.")

    moveCommand.Flags().StringVar(
        &flagRequestPayer,
        "request-payer",
        "",
        "Confirms that the requester knows that she or he will be charged for the request.")

    moveCommand.Flags().IntVar(
        &flagConcurrency,
        "concurrency",
        1,
        "Number of concurrent workers (e.g., goroutines) to spin up.")

	RootCmd.AddCommand(moveCommand)
}
