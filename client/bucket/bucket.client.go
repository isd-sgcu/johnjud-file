package bucket

import (
	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/minio/minio-go/v7"
)

type Client struct {
	conf cfgldr.Bucket
	// s3   *s3.Client
	minio *minio.Client
}

// func NewClient(conf cfgldr.Bucket, awsClient *s3.Client) *Client {
func NewClient(conf cfgldr.Bucket, minioClient *minio.Client) *Client {
	return &Client{conf: conf, minio: minioClient}
}

func (c *Client) Upload(file []byte, objectKey string) (string, string, error) {
	return "", "", nil
}

func (c *Client) Delete(objectKey string) error {

	return nil
}

// func (c *Client) Upload(file []byte, objectKey string) (string, string, error) {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
// 	defer cancel()

// 	buffer := bytes.NewReader(file)
// 	var partMiBs int64 = 10
// 	uploader := manager.NewUploader(c.s3, func(u *manager.Uploader) {
// 		u.PartSize = partMiBs * 1024 * 1024
// 	})

// 	uploadOutput, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
// 		Bucket: aws.String(c.conf.BucketName),
// 		Key:    aws.String(objectKey),
// 		Body:   buffer,
// 	})

// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Str("service", "file").
// 			Str("module", "bucket client").
// 			Msgf("Couldn't upload object to %v:%v.", c.conf.BucketName, objectKey)

// 		return "", "", errors.Wrap(err, "Error while uploading the object")
// 	}

// 	return uploadOutput.Location, *uploadOutput.Key, nil
// }

// func (c *Client) Delete(objectKey string) error {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
// 	defer cancel()

// 	input := &s3.DeleteObjectInput{
// 		Bucket: aws.String(c.conf.BucketName),
// 		Key:    aws.String(objectKey),
// 	}

// 	_, err := c.s3.DeleteObject(context.TODO(), input)

// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Str("service", "file").
// 			Str("module", "bucket client").
// 			Msgf("Couldn't delete object from bucket %v:%v.", c.conf.BucketName, objectKey)

// 		return errors.Wrap(err, "Error while deleting the object")
// 	}

// 	return nil
// }
