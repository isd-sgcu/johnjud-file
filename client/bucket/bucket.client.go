package bucket

import (
	"bytes"
	"context"
	"time"

	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Client struct {
	conf  cfgldr.Bucket
	minio *minio.Client
}

func NewClient(conf cfgldr.Bucket, minioClient *minio.Client) *Client {
	return &Client{conf: conf, minio: minioClient}
}

func (c *Client) Upload(file []byte, objectKey string) (string, string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	buffer := bytes.NewReader(file)

	uploadOutput, err := c.minio.PutObject(context.Background(), c.conf.BucketName, objectKey, buffer,
		buffer.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "file").
			Str("module", "bucket client").
			Msgf("Couldn't upload object to %v:%v.", c.conf.BucketName, objectKey)

		return "", "", errors.Wrap(err, "Error while uploading the object")
	}

	return c.getURL(objectKey), uploadOutput.Key, nil
}

func (c *Client) Delete(objectKey string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err := c.minio.RemoveObject(context.Background(), c.conf.BucketName, objectKey, opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "file").
			Str("module", "bucket client").
			Msgf("Couldn't delete object from bucket %v:%v.", c.conf.BucketName, objectKey)

		return errors.Wrap(err, "Error while deleting the object")
	}

	return nil
}

func (c *Client) getURL(objectKey string) string {
	return "https://" + c.conf.Endpoint + "/" + c.conf.BucketName + "/" + objectKey
}
