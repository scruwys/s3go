package s3go

import (
    "time"
    "fmt"
    "errors"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client struct {
    // Authenticated S3 client to make API requests
    svc *s3.S3

    // Turn on debug logging
    debug bool

    // Specifies the AWS Region where the API actions will be scoped to
    region string
}

type ClientOptions struct {
    // Endpoint to make S3 API calls against
    Endpoint string

    // Turn on debug logging
    Debug bool

    // This option overrides the default behavior of verifying SSL certificates
    DisableSSL bool

    // Use a specific profile from your credential file
    Profile string

    // Specifies the AWS Region where the API actions will be scoped to
    Region string
}

func NewClient(options *ClientOptions) *Client {
    cfg := &aws.Config{
        Endpoint:   &options.Endpoint,
        DisableSSL: &options.DisableSSL,
    }

    if options.Region != "" {
        cfg = cfg.WithRegion(options.Region)
    }

    sess, _ := session.NewSessionWithOptions(session.Options{
        Profile: options.Profile,
    })

    svc := s3.New(sess, cfg)

    return &Client{
        svc:    svc,
        debug:  options.Debug,
        region: options.Region,
    }
}

// Generates a pre-signed URL for an Amazon S3 object
func (c *Client) Presign(bucketName string, key string, expireInSeconds int) (string, error) {
    req, _ := c.svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
    })

    presignedUrl, err := req.Presign(time.Duration(expireInSeconds) * time.Second)

    if err != nil {
        return "", err
    }

    return presignedUrl, nil
}

// Create an S3 bucket if one does not exist with that name.
func(c *Client) MakeBucket(bucketName, region string) error {
    _, err := c.svc.CreateBucket(&s3.CreateBucketInput{
        Bucket: aws.String(bucketName),
    })

    if err != nil {
        return err
    }

    err = c.svc.WaitUntilBucketExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucketName),
    })

    return err
}

// Removes an S3 bucket if it  exists within the context of the current session.
func(c *Client) RemoveBucket(bucketName string, forceFlag bool) error {
    if (forceFlag) {
        return c.EmptyBucket(bucketName)
    }

    _, err := c.svc.DeleteBucket(&s3.DeleteBucketInput{
        Bucket: aws.String(bucketName),
    })

    if err != nil {
        return err
    }

    err = c.svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucketName),
    })

    return err
}

// Checks to see if an S3 bucket exists within the context of the current session.
func(c *Client) BucketExists(bucketName string) error {
    _, err := c.svc.HeadBucket(&s3.HeadBucketInput{Bucket: aws.String(bucketName)})
    return err
}

// Empty an entire S3 bucket using batch DeleteObjects API operations.
func(c *Client) EmptyBucket(bucketName string) error {
    iter := s3manager.NewDeleteListIterator(c.svc, &s3.ListObjectsInput{
        Bucket: aws.String(bucketName),
    })

    return s3manager.NewBatchDeleteWithClient(c.svc).Delete(aws.BackgroundContext(), iter)
}

// Executes RestoreObject API operation on a single S3 key.
func(c *Client) RestoreObject(bucketName, key, requestPayer string) error {
    _, err := c.svc.RestoreObject(&s3.RestoreObjectInput{
        Bucket:       aws.String(bucketName),
        Key:          aws.String(key),
        RequestPayer: aws.String(requestPayer),
    })

    return err
}

// Executes DeleteObject API operation on a single S3 key.
func(c *Client) DeleteObject(bucketName, key, requestPayer string) error {
    _, err := c.svc.DeleteObject(&s3.DeleteObjectInput{
        Bucket:       aws.String(bucketName),
        Key:          aws.String(key),
        RequestPayer: aws.String(requestPayer),
    })

    return err
}

// TODO(scruwys)
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/s3/s3_download_object.go
func(c *Client) DownloadObject(object ObjectInfo, target *S3Url, keepSourceFlag bool) error {
    fmt.Println("DownloadObject: ", *object.Key)
    return nil
}

// TODO(scruwys)
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/s3/s3_upload_object.go
func(c *Client) UploadObject(object ObjectInfo, target *S3Url, keepSourceFlag bool) error {
    fmt.Println("UploadObject: ", *object.Key)
    return nil
}

// TODO(scruwys)
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/s3/s3_copy_object.go
func(c *Client) CopyObject(object ObjectInfo, target *S3Url, keepSourceFlag bool) error {
    fmt.Println("CopyObject: ", *object.Key)
    return nil
}

// TODO(scruwys)
func(c *Client) MoveObject(object ObjectInfo, source *S3Url, target *S3Url, keepSourceFlag bool) error {
    if !source.IsLocal() && target.IsLocal() {
        return c.DownloadObject(object, target, keepSourceFlag)
    }

    if !source.IsLocal() && !target.IsLocal() {
        return c.CopyObject(object, target, keepSourceFlag)
    }

    if source.IsLocal() && !target.IsLocal() {
        return c.UploadObject(object, target, keepSourceFlag)
    }

    return errors.New("This action is not supported.")
}

type ObjectInfo struct {
    // The entity tag is an MD5 hash of the object. ETag reflects only changes to
    // the contents of an object, not its metadata.
    ETag *string `type:"string"`

    // The name that you assign to an object. You use the object key to retrieve
    // the object.
    Key *string `min:"1" type:"string"`

    // The date the Object was Last Modified
    LastModified *time.Time `type:"timestamp"`

    // Size in bytes of the object
    Size *int64 `type:"integer"`

    // The class of storage used to store the object.
    StorageClass *string `type:"string"`

    // Indicates whether or not the ObjectInfo represents a prefix.
    IsPrefix bool

    // The name of the Amazon S3 bucket
    Bucket *string `type:"string"`
}

type ListObjectsV2Input struct {
    // The name of the Amazon S3 bucket you want to list
    Bucket string

    // Limits the response to keys that begin with the specified prefix
    Prefix string

    // List all files or objects under the specified directory or prefix
    Recursive bool

    // Don't exclude files or objects in the command that match the specified pattern
    ExcludeFilter string

    // Exclude all files or objects from the command that matches the specified pattern
    IncludeFilter string

}

// List objects in an S3 bucket and prefix using the list-objects-v2 API method.
func(c *Client) ListObjectsV2(options *ListObjectsV2Input) (ch <-chan ObjectInfo, err error) {
    delimiter := ""
    if !options.Recursive {
        delimiter = "/"
    }

    if err := c.BucketExists(options.Bucket); err != nil {
        return nil, err
    }

    input := &s3.ListObjectsV2Input{
        Bucket:    aws.String(options.Bucket),
        Prefix:    aws.String(options.Prefix),
        Delimiter: aws.String(delimiter),
    }

    outputCh := make(chan ObjectInfo)

    go func() {
        defer close(outputCh)

        c.svc.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
            for _, prefix := range page.CommonPrefixes {
                objectInfo := ObjectInfo{
                    Bucket:       &options.Bucket,
                    Key:          prefix.Prefix,
                    IsPrefix:     true,
                    Size:         new(int64),
                    LastModified: &time.Time{},
                }
                outputCh <- objectInfo
            }

            for _, object := range page.Contents {
                objectInfo := ObjectInfo{
                    Bucket:       &options.Bucket,
                    Key:          object.Key,
                    IsPrefix:     false,
                    Size:         object.Size,
                    LastModified: object.LastModified,
                }
                outputCh <- objectInfo
            }

            return !lastPage
        })
    }()

    return outputCh, nil
}

type ListSourceObjectsInput struct {
    // Object representation of the source path
    SourceUrl *S3Url

    // List all files or objects under the specified directory or prefix
    Recursive bool

    // Don't exclude files or objects in the command that match the specified pattern
    ExcludeFilter string

    // Exclude all files or objects from the command that matches the specified pattern
    IncludeFilter string
}

// List source objects from either S3 or the local file system. Used for "cp" and "mv" commands, etc.
func(c *Client) ListSourceObjects(options *ListSourceObjectsInput) (ch <-chan ObjectInfo, err error) {
    if options.SourceUrl.IsLocal() {
        return listFiles(options.SourceUrl.Prefix, options.Recursive)
    }

    input := &ListObjectsV2Input{
        Bucket:        options.SourceUrl.Bucket,
        Prefix:        options.SourceUrl.Prefix,
        Recursive:     options.Recursive,
        ExcludeFilter: options.ExcludeFilter,
        IncludeFilter: options.IncludeFilter,
    }

    return c.ListObjectsV2(input)
}
