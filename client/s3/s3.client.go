package s3

import (
	"github.com/isd-sgcu/johnjud-file/cfgldr"
)

type Client struct {
	conf cfgldr.S3
	s3   AWSClient
}

type AWSClient interface {
	// Upload([]byte, string) error
	// GetSignedUrl(string) (string, error)
}

func NewClient(conf cfgldr.S3, s3 AWSClient) *Client {
	return &Client{conf: conf, s3: s3}
}

func (c *Client) Upload([]byte, string) error {
	return nil
}

func (c *Client) GetSignedUrl(string) (string, error) {
	return "", nil
}
