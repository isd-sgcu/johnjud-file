package s3

import (
	"github.com/isd-sgcu/johnjud-file/cfgldr"
)

type clientImpl struct {
	conf cfgldr.S3
}

func NewClient(conf cfgldr.S3) *clientImpl {
	return &clientImpl{conf: conf}
}

func (c *clientImpl) Upload([]byte, string) error {
	return nil
}

func (c *clientImpl) GetSignedUrl(string) (string, error) {
	return "", nil
}
