package bucket

import (
	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/isd-sgcu/johnjud-file/client/bucket"
	"github.com/minio/minio-go/v7"
)

type Client interface {
	Upload([]byte, string) (string, string, error)
	Delete(string) error
}

// func NewClient(config cfgldr.Bucket, awsClient *s3.Client) Client {
func NewClient(config cfgldr.Bucket, minioClient *minio.Client) Client {
	return bucket.NewClient(config, minioClient)
}
