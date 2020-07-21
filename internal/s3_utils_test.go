package s3go

import (
    "testing"
)

func TestS3Url(t *testing.T) {
	s3Url := S3Url{"s3", "s3go-testing", "config/at/this/path.json"}

	if s3Url.IsLocal() {
		t.Errorf("A valid S3 url should not be considered local")
	}

	fsUrl := S3Url{"", "", "./config/at/this/path.json"}

	if !fsUrl.IsLocal() {
		t.Errorf("A FS object should always be considered local")
	}
}

func TestParseUrl(t *testing.T) {
    var tests = []struct {
        path, scheme, bucket, prefix string
    }{
        {"./s3go/README", "", " ", "./s3go/README"},
        {"./s3go/internal/fs_test.go", "", " ", "./s3go/internal/fs_test.go"},
        {"/Users/computron/dev/src/s3go/README", "", " ", "/Users/computron/dev/src/s3go/README"},
        {"/ok.go", "", " ", "/ok.go"},
        {"s3://s3go-testing/config/at/this/path.json", "s3", "s3go-testing", "config/at/this/path.json"},
        {"s3://s3go-testing/folder/structure/honored/", "s3", "s3go-testing", "folder/structure/honored/"},
        {"s3://s3go-nopath/", "s3", "s3go-nopath", ""},
    }

    for _, tt := range tests {
        t.Run(tt.path, func(t *testing.T) {
        	url, err := ParseUrl(tt.path)

        	if err != nil {
        		t.Errorf("%v", err)
        	}

        	if url.Scheme != tt.scheme {
                t.Errorf("got %s, want %s", url.Scheme, tt.scheme)
        	}

        	if url.Bucket != tt.bucket {
                t.Errorf("got %s, want %s", url.Bucket, tt.bucket)
        	}

        	if url.Prefix != tt.prefix {
                t.Errorf("got %s, want %s", url.Prefix, tt.prefix)
        	}
        })
    }
}

func TestGetBucketRegion(t *testing.T) {

}

func TestBuildTargetPrefix(t *testing.T) {

}
