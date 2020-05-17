package s3go

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Client struct{
	// svc *minio.Client
}

func New() {

}

func (c *Client) Get() string {
    endpoint := "s3.dualstack.us-west-2.amazonaws.com"
    accessKeyID := ""
    secretAccessKey := ""
    useSSL := true

    // Initialize minio client object.
    minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("%#v\n", minioClient) // minioClient is now setup

	return "Ok"
}
