package awscli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ErrObjectNotFound = errors.New("s3 object not found")
)

// S3Client is a wrapper for the S3 client in the AWS SDK.
type S3Client struct {
	s3 *s3.Client
}

// NewS3Client creates a new instance of the S3 client wrapper.
func NewS3Client(ctx context.Context, opts ...Option) (*S3Client, error) {
	cfg, err := loadConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("configure S3 client: %w", err)
	}

	return &S3Client{
		s3: s3.NewFromConfig(cfg),
	}, nil
}

// GetObject returns
func (c *S3Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	resp, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *S3Client) PutObject() {

}

func (c *S3Client) DeleteObject() {

}

func (c *S3Client) ListObjects() {

}
