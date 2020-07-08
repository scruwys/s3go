package s3go

import (
    "time"
    "fmt"
    "errors"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Client struct {
    // Authenticated S3 client to make API requests
    svc *s3.S3

    // TBD...
    sess *session.Session

    // Original options used to configure the client
    options *ClientOptions

    uploader *s3manager.Uploader

    downloader *s3manager.Downloader
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
    sess, _ := session.NewSessionWithOptions(session.Options{
        Profile: options.Profile,
    })

    svc := s3.New(sess, NewConfig(options))

    uploader := s3manager.NewUploaderWithClient(svc)

    downloader := s3manager.NewDownloaderWithClient(svc)

    return &Client{svc, sess, options, uploader, downloader}
}

func NewConfig(options *ClientOptions) *aws.Config {
    cfg := &aws.Config{
        Endpoint:   &options.Endpoint,
        DisableSSL: &options.DisableSSL,
    }

    if options.Region != "" {
        cfg = cfg.WithRegion(options.Region)
    }

    return cfg
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
        fmt.Println("Emptying the bucket. This might take a few minutes.")
        return c.EmptyBucket(bucketName)
    }

    // We override the session since it's best to delete buckets from the main S3 endpoint.
    svc := s3.New(c.sess, NewConfig(&ClientOptions{
        Endpoint: fmt.Sprintf(DEFAULT_ENDPOINT_URL),
        Region:   c.options.Region,
    }))

    _, err := svc.DeleteBucket(&s3.DeleteBucketInput{
        Bucket: aws.String(bucketName),
    })

    if err != nil {
        return err
    }

    err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
        Bucket: aws.String(bucketName),
    })

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
func(c *Client) DownloadObject(object ObjectInfo, targetPrefix string, options *MoveObjectOptions) (string, error) {
    sourcePath := fmt.Sprintf("s3://%s/%s", options.Source.Bucket, *object.Key)
    targetPath := fmt.Sprintf("%s", targetPrefix)
    logMessage := fmt.Sprintf("%s to %s", sourcePath, targetPath)

    if options.DryRun {
        return logMessage, nil
    }

    file, err := createFile(targetPrefix)

    if err != nil {
        return "", err
    }

    defer file.Close()

    _, err = c.downloader.Download(file, &s3.GetObjectInput{
        Bucket: aws.String(options.Source.Bucket),
        Key:    aws.String(*object.Key),
    })

    if err != nil {
        return "", err
    }

    if options.DeleteAfter {
        err = c.DeleteObject(options.Source.Bucket, *object.Key, options.RequestPayer)

        if err != nil {
            return "", err
        }
    }

    return logMessage, nil
}

// TODO(scruwys)
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/s3/s3_upload_object.go
func(c *Client) UploadObject(object ObjectInfo, targetPrefix string, options *MoveObjectOptions) (string, error) {
    sourcePath := fmt.Sprintf("%s", object.Path())
    targetPath := fmt.Sprintf("s3://%s/%s", options.Target.Bucket, targetPrefix)
    logMessage := fmt.Sprintf("%s to %s", sourcePath, targetPath)

    if options.DryRun {
        return logMessage, nil
    }

    file, err := os.Open(sourcePath)

    if err != nil {
        return "", err
    }

    defer file.Close()

    _, err = c.uploader.Upload(&s3manager.UploadInput{
        Bucket:     aws.String(options.Target.Bucket),
        Key:        aws.String(targetPrefix),
        Body:       file,
    })

    if err != nil {
        return "", err
    }

    if options.DeleteAfter {
        err = os.Remove(sourcePath)

        if err != nil {
            return "", err
        }
    }

    return logMessage, nil
}

// TODO(scruwys)
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/s3/s3_copy_object.go
func(c *Client) CopyObject(object ObjectInfo, targetPrefix string, options *MoveObjectOptions) (string, error) {
    sourcePath := fmt.Sprintf("s3://%s", object.Path())
    targetPath := fmt.Sprintf("s3://%s/%s", options.Target.Bucket, targetPrefix)

    logMessage := fmt.Sprintf("%s to %s", sourcePath, targetPath)

    if options.DryRun {
        return logMessage, nil
    }

    _, err := c.svc.CopyObject(&s3.CopyObjectInput{
        CopySource: aws.String(object.Path()),
        Bucket:     aws.String(options.Target.Bucket),
        Key:        aws.String(targetPrefix),
    })

    if err != nil {
        return "", err
    }

    err = c.svc.WaitUntilObjectExists(&s3.HeadObjectInput{
        Bucket: aws.String(options.Target.Bucket),
        Key:    aws.String(targetPrefix),
    })

    if err != nil {
        return "", err
    }

    if options.DeleteAfter {
        err = c.DeleteObject(*object.Bucket, *object.Key, options.RequestPayer)

        if err != nil {
            return "", err
        }
    }

    return logMessage, nil
}

type MoveObjectOptions struct {
    Target       *S3Url
    Source       *S3Url
    DryRun       bool
    Recursive    bool
    RequestPayer string
    DeleteAfter  bool
}

func(c *Client) MoveObject(object ObjectInfo, options *MoveObjectOptions) (string, error) {
    targetPrefix := buildTargetPrefix(options.Target.Prefix, options.Source.Prefix, *object.Key, options.Recursive)

    if !options.Source.IsLocal() && options.Target.IsLocal() {
        return c.DownloadObject(object, targetPrefix, options)
    }

    if !options.Source.IsLocal() && !options.Target.IsLocal() {
        return c.CopyObject(object, targetPrefix, options)
    }

    if options.Source.IsLocal() && !options.Target.IsLocal() {
        return c.UploadObject(object, targetPrefix, options)
    }

    return "", errors.New("This action is not supported.")
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

func (o *ObjectInfo) Path() string {
    delimiter := "/"
    if *o.Bucket == "" {
        delimiter = ""
    }
    return fmt.Sprintf("%s%s%s", *o.Bucket, delimiter, *o.Key)
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

    region, err := GetBucketRegion(options.Bucket)

    if err != nil {
        return nil, err
    }

    // We override the session to list objects from whatever target bucket.
    svc := s3.New(c.sess, NewConfig(&ClientOptions{
        Endpoint: fmt.Sprintf("s3.%s.amazonaws.com", region),
        Region:   region,
    }))

    if !options.Recursive {
        emptyCh := make(chan ObjectInfo)
        defer close(emptyCh)

        if options.Prefix == "" {
            return emptyCh, nil
        }

        _, err := svc.HeadObject(&s3.HeadObjectInput{
            Bucket: aws.String(options.Bucket),
            Key:    aws.String(options.Prefix),
        })

        if err != nil {
            return nil, err
        }
    }

    input := &s3.ListObjectsV2Input{
        Bucket:    aws.String(options.Bucket),
        Prefix:    aws.String(options.Prefix),
        Delimiter: aws.String(delimiter),
    }

    excludeRe, err := regexpCompile(options.ExcludeFilter, "$^")
    if err != nil {
        return nil, err
    }

    includeRe, err := regexpCompile(options.IncludeFilter, ".*")
    if err != nil {
        return nil, err
    }

    outputCh := make(chan ObjectInfo)

    go func() {
        defer close(outputCh)

        svc.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
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
                if excludeRe.MatchString(*object.Key) || !includeRe.MatchString(*object.Key) {
                    continue
                }
                objectKey := *object.Key
                objectInfo := ObjectInfo{
                    Bucket:       &options.Bucket,
                    Key:          object.Key,
                    IsPrefix:     objectKey[len(objectKey)-1:] == "/",
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
        return listFiles(options.SourceUrl.Prefix, options.Recursive, options.ExcludeFilter, options.IncludeFilter)
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
