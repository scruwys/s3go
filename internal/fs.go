package s3go

import (
    "path/filepath"
    "errors"
    "fmt"
	"os"
)

func pathExists(path string) (error) {
    _, err := os.Stat(path)
    if err == nil { return nil }
    if os.IsNotExist(err) {
        return errors.New(fmt.Sprintf("The user-provided path %s does not exist.", path))
    }
    return err
}

func ensureDir(dirName string) error {
    err := os.Mkdir(dirName, 0700)

    if err == nil || os.IsExist(err) {
        return nil
    } else {
        return err
    }
}


func createFile(path string) (*os.File, error) {
    dir, _ := filepath.Split(path)

    if err := ensureDir(dir); err != nil {
        return nil, err
    }

    return os.Create(path)
}


func listFiles(rootPath string, recursive bool, excludeFilter string, includeFilter string) (ch <-chan ObjectInfo, err error) {
    outputCh := make(chan ObjectInfo)

    excludeRe, err := regexpCompile(excludeFilter, "$^")
    if err != nil {
        return nil, err
    }

    includeRe, err := regexpCompile(includeFilter, ".*")
    if err != nil {
        return nil, err
    }

    if err := pathExists(rootPath); err != nil {
        return nil, err
    }

    go func() {
        defer close(outputCh)

        filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
            // We ignore anything that doesn't pass the filter checks
            if excludeRe.MatchString(path) || !includeRe.MatchString(path) {
                return nil
            }

            // This ignores the root directory, but only recurses if commanded
            if info.IsDir() && !recursive && rootPath != path {
                return filepath.SkipDir
            }

            if !recursive && rootPath != path {
                return nil
            }

            // If it's not a directory, it's a file, so let's process it
            if !info.IsDir() {
                emptyStr := ""
                fmt.Println(path)

                objectInfo := ObjectInfo{
                    Bucket:   &emptyStr,
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

