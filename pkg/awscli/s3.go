package awscli

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3 represents the mockable functions in the AWS SDK S3 client.
type S3 interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// S3Client is a wrapper for the S3 client in the AWS SDK.
type S3Client struct {
	s3         S3
	bucketName string
}

// NewS3Client creates a new instance of the S3 client wrapper.
func NewS3Client(ctx context.Context, bucketName string, opts ...Option) (*S3Client, error) {
	cfg, err := loadConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("s3client: %w", err)
	}

	return &S3Client{
		s3:         s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}, nil
}

// Object represents object retrieved from S3.
type Object struct {
	bucket string
	key    string
	body   io.ReadCloser
}

// GetObject returns *Object.
func (c *S3Client) GetObject(ctx context.Context, key string) (*Object, error) {
	resp, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3client: %w", err)
	}

	return &Object{bucket: c.bucketName, key: key, body: resp.Body}, nil
}

// PutObject uploads data from reader to S3 Bucket.
func (c *S3Client) PutObject(ctx context.Context, key string, reader io.Reader) error {
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(c.bucketName), Key: aws.String(key), Body: reader})
	if err != nil {
		return fmt.Errorf("s3client: %w", err)
	}

	return nil
}

// DeleteObject deletes an object from S3 Bucket by key.
func (c *S3Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(c.bucketName), Key: aws.String(key)})
	if err != nil {
		return fmt.Errorf("s3client: %w", err)
	}

	return nil
}

// ListObjects search keys by prefix; returns all keys if prefix is empty string.
func (c *S3Client) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	l, err := c.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(c.bucketName), Prefix: aws.String(prefix)})
	if err != nil {
		return nil, fmt.Errorf("s3client: %w", err)
	}

	var keysList = make([]string, 0, len(l.Contents))
	for _, key := range l.Contents {
		keysList = append(keysList, aws.ToString(key.Key))
	}

	return keysList, nil
}
