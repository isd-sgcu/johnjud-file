package s3

import "github.com/isd-sgcu/johnjud-file/cfgldr"

type Client struct {
	conf cfgldr.S3
}

func NewClient(conf cfgldr.S3) *Client {
	return &Client{conf: conf}
}

func (c *Client) Upload() {

}

func (c *Client) Delete() {

}
