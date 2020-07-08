package main

import (
	"github.com/scruwys/s3go/cmd"
	"github.com/scruwys/s3go/internal"
)


func main() {
	err := cmd.RootCmd.Execute()

	if err != nil {
		s3go.ExitWithError(1, err)
	}
}
