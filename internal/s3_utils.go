package s3go

import (
	"errors"
	"fmt"
	"net/http"
    "path/filepath"
	"net/url"
	"strings"
)

const (
    DEFAULT_ENDPOINT_URL = "s3.amazonaws.com"
)

type S3Url struct {
	// The proto scheme of the URL. Should usually be only S3 or blank.
	Scheme  string

    // The name of the Amazon S3 bucket
	Bucket  string

    // The prefix of the S3 object / path
	Prefix  string
}

func (u *S3Url) IsLocal() bool {
	return u.Scheme != "s3"
}

func ParseUrl(input string) (*S3Url, error) {
	u, err := url.Parse(input)

	if err != nil {
		return nil, err
	}

	prefix := ""

	if u.Path != "/" {
		prefix = u.Path

		if u.Scheme == "s3" {
			prefix = strings.TrimLeft(u.Path, "/")
		}
	}

	bucket := " "

	if u.Host != "" {
		bucket = u.Host
	}

	output := &S3Url{
		Scheme: u.Scheme,
		Bucket: bucket,
		Prefix: prefix,
	}

	return output, nil
}

func GetBucketRegion(bucket string) (string, error) {
    url := fmt.Sprintf("https://%s.%s", bucket, DEFAULT_ENDPOINT_URL)
    res, err := http.Head(url)
    if err != nil {
        return "", err
    }
    if res.StatusCode == 404 {
        return "", errors.New("(NoSuchBucket) Provided bucket does not exist")
    }
    return res.Header.Get("X-Amz-Bucket-Region"), nil
}

func buildTargetPrefix(targetPrefix, sourcePrefix, objectKey string, recursive bool) string {
    dir, fname := filepath.Split(objectKey)

    if targetPrefix == "" {
    	return fname
    }

    if recursive {
	    substr := IntMin(len(sourcePrefix), len(dir))
        targetPrefix = strings.TrimRight(targetPrefix, "/") + "/" + dir[substr:]
    }

    if recursive || (len(targetPrefix) > 0 && targetPrefix[len(targetPrefix)-1:] == "/") {
        targetPrefix = targetPrefix + fname
    }

    return targetPrefix
}
