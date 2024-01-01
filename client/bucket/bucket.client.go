package bucket

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Client struct {
	conf cfgldr.S3
	s3   *s3.Client
}

func NewClient(conf cfgldr.S3, awsClient *s3.Client) *Client {
	return &Client{conf: conf, s3: awsClient}
}

func (c *Client) Upload(file []byte, objectKey string) (string, string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	buffer := bytes.NewReader(file)
	var partMiBs int64 = 10
	uploader := manager.NewUploader(c.s3, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	uploadOutput, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.conf.BucketName),
		Key:    aws.String(objectKey),
		Body:   buffer,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "file").
			Str("module", "bucket client").
			Msgf("Couldn't upload object to %v:%v.", c.conf.BucketName, objectKey)

		return "", "", errors.Wrap(err, "Error while uploading the object")
	}
	urlObjectKey := strings.Replace(objectKey, " ", "+", -1)

	return fmt.Sprintf("https://%v.s3.%v.amazonaws.com/%v", c.conf.BucketName, c.conf.Region, urlObjectKey), *uploadOutput.Key, nil
}

func (c *Client) Delete(objectKey string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	var objectIds []types.ObjectIdentifier
	objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(objectKey)})

	_, err := c.s3.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(c.conf.BucketName),
		Delete: &types.Delete{Objects: objectIds},
	})

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
