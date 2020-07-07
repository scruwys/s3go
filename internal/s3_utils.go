package s3go

import (
	"net/url"
	"strings"
)

type S3Url struct {
	// ok
	Scheme  string

	// ok
	Bucket  string

	// ok
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
		prefix = strings.TrimLeft(u.Path, "/")
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
