# s3go

`s3go` is a simple Golang port of the [AWS S3 CLI](https://docs.aws.amazon.com/cli/latest/reference/s3/). It was built as a Hack Week project with the primary objective being to learn some basic features of Go.

The CLI aims to use the concurrency features of Golang to improve performance for simple S3 tasks, such as copying files or restoring files from Glacier.

## Installation

```
$ go get -u github.com/scruwys/s3go
```

## Usage

```
$ s3go --help

golang cli for interacting with aws s3.

Usage:
  s3go [command]

Available Commands:
  cp          Copies a local file or S3 object to another location locally or in S3.
  help        Help about any command
  ls          List S3 objects and common prefixes under a prefix or all S3 buckets.
  mb          Creates an S3 bucket.
  mv          Moves a local file or S3 object to another location locally or in S3.
  presign     Generate a pre-signed URL for an Amazon S3 object.
  rb          Deletes an empty S3 bucket.
  restore     Restores S3 object(s) stored in Glacier.
  rm          Deletes an S3 object.

Flags:
      --debug                 Turn on debug logging.
      --endpoint-url string   Override command's default URL with the given URL.
  -h, --help                  help for s3go
      --no-verify-ssl         This option overrides the default behavior of verifying SSL certificates.
      --profile string        Use a specific profile from your credential file.
      --region string         The region to use. Overrides config/env settings. (default "us-east-1")

Use "s3go [command] --help" for more information about a command.
```

Every command has been implemented except for `sync`. Intial benchmarks show that most do not perform better than the existing [awscli](https://github.com/aws/aws-cli) tool. More on that soon.

The primary difference is the addition of the `--concurrency` flag, which allows you to control how many goroutines are spun up to handle command executation:

```
s3go rm s3://my-test-bucket/20200101/tmp/ --recursive --concurrency 5
```

## License

Released under the [MIT license](LICENSE).
