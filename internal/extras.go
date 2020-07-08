package s3go

import (
	"fmt"
	"os"
	"sync"
    "strconv"
    "regexp"
)

func Echo(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func ExitWithError(code int, err error) {
	fmt.Println(err)
	os.Exit(code)
}

func HumanizeBytes(b int64) string {
    const unit = 1024
    if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func ByteSizeToString(size int64, isflagHumanReadable bool) string {
    sz := strconv.FormatInt(size, 10)
    if isflagHumanReadable {
        sz = HumanizeBytes(size)
    }
    return sz
}

// https://medium.com/justforfunc/why-are-there-nil-channels-in-go-9877cc0b2308
// https://github.com/jakewright/tutorials/blob/master/go/02-go-concurrency/05-fib.go
func MergeWaitWithObjectInfo(done <-chan bool, channels ...<-chan ObjectInfo) <-chan ObjectInfo {
    var wg sync.WaitGroup
    out := make(chan ObjectInfo)

    multiplex := func(c <-chan ObjectInfo) {
        defer wg.Done()
        for i := range c {
            select {
                case <-done:
                    return
                case out <- i:
            }
        }
    }

    wg.Add(len(channels))

    for _, c := range channels {
        go multiplex(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

func regexpCompile(pattern, defaultTo string) (*regexp.Regexp, error) {
    if pattern == "" {
        pattern = defaultTo
    }
    return regexp.Compile(pattern)
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
