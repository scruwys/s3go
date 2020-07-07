package s3go

import (
    "path/filepath"
	"os"
)

func listFiles(rootPath string, recursive bool) (ch <-chan ObjectInfo, err error) {
    outputCh := make(chan ObjectInfo)

    go func() {
        defer close(outputCh)

        filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
            if info.IsDir() && !recursive && rootPath != path {
                return filepath.SkipDir
            }

            if !info.IsDir() {
                objectInfo := ObjectInfo{
                    Key:      &path,
                    IsPrefix: false,
                }
                outputCh <- objectInfo
            }

            return nil
        })
    }()

    return outputCh, nil
}

