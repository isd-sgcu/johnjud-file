package s3

import "github.com/isd-sgcu/johnjud-file/src/config"

type Client struct {
	conf config.S3
}

func NewClient(conf config.S3) *Client {
	return &Client{conf: conf}
}

func (c *Client) Upload() {

}

func (c *Client) Delete() {

}
