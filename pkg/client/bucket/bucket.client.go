package bucket

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/isd-sgcu/johnjud-file/client/bucket"
)

type Client interface {
	Upload([]byte, string) error
	GetSignedUrl(string) (string, error)
}

func NewClient(config cfgldr.S3, awsClient *s3.Client) Client {
	return bucket.NewClient(config, awsClient)
}
