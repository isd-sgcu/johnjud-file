package s3

import (
	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/isd-sgcu/johnjud-file/internal/client/s3"
)

type Client interface {
	Upload([]byte, string) error
	GetSignedUrl(string) (string, error)
}

func NewClient(config cfgldr.S3) Client {
	return s3.NewClient(config)
}
