package s3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
)

func TestNewS3Client(t *testing.T) {
	t.Parallel()

	got, err := New(context.TODO(), "name")

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "name", got.bucketName)

	awsLoadDefaultConfigFn = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error) {
		return aws.Config{}, fmt.Errorf("config err")
	}
	defer func() {
		awsLoadDefaultConfigFn = config.LoadDefaultConfig
	}()

	got, err = New(context.TODO(), "name")
	require.Error(t, err)
	require.Nil(t, got)
}

type s3mock struct {
	getFn  func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	listFn func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	putFn  func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	delFn  func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func (s s3mock) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return s.getFn(ctx, params, optFns...)
}

func (s s3mock) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return s.listFn(ctx, params, optFns...)
}

func (s s3mock) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return s.putFn(ctx, params, optFns...)
}

func (s s3mock) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return s.delFn(ctx, params, optFns...)
}

func TestS3Client_GetObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		bucket  string
		mock    S3
		want    *Object
		wantErr bool
	}{

		{
			name:   "success",
			key:    "k1",
			bucket: "bucket",
			mock: s3mock{getFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader("test str")),
				}, nil
			}},
			want: &Object{
				bucket: "bucket",
				key:    "k1",
				body:   io.NopCloser(strings.NewReader("test str")),
			},
			wantErr: false,
		},

		{
			name:   "error",
			key:    "k1",
			bucket: "bucket",
			mock: s3mock{getFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, fmt.Errorf("some err")
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			cli, err := New(ctx, tt.bucket)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.s3 = tt.mock

			got, err := cli.Get(ctx, tt.key)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, got)

			expectedBytes, err := io.ReadAll(tt.want.body)
			require.NoError(t, err)
			gotBytes, err := io.ReadAll(got.body)
			require.NoError(t, err)

			require.Equal(t, string(expectedBytes), string(gotBytes))
		})
	}
}

func TestS3Client_PutObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		bucket  string
		mock    S3
		wantErr bool
	}{
		{
			name:   "success",
			key:    "k1",
			bucket: "bucket",
			mock: s3mock{putFn: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			}},
			wantErr: false,
		},
		{
			name:   "error",
			key:    "k1",
			bucket: "bucket",
			mock: s3mock{putFn: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return nil, fmt.Errorf("some err")
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			cli, err := New(ctx, tt.bucket)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.s3 = tt.mock

			err = cli.Put(ctx, tt.key, nil)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestS3Client_DeleteObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		bucket  string
		mock    S3
		wantErr bool
	}{
		{
			name:   "success",
			key:    "k1",
			bucket: "bucket",
			mock: s3mock{delFn: func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return &s3.DeleteObjectOutput{}, nil
			}},
			wantErr: false,
		},
		{
			name:   "error",
			key:    "k1",
			bucket: "bucket",
			mock: s3mock{delFn: func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				return nil, fmt.Errorf("some err")
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			cli, err := New(ctx, tt.bucket)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.s3 = tt.mock

			err = cli.Delete(ctx, tt.key)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestS3Client_ListObject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prefix  string
		bucket  string
		mock    S3
		want    []string
		wantErr bool
	}{
		{
			name:   "success - all",
			prefix: "",
			bucket: "bucket",
			mock: s3mock{listFn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return &s3.ListObjectsV2Output{
					Contents: []types.Object{
						{Key: aws.String("key1")},
						{Key: aws.String("another_key")},
					},
				}, nil
			}},
			want:    []string{"key1", "another_key"},
			wantErr: false,
		},
		{
			name:   "success - prefix",
			prefix: "ke",
			bucket: "bucket",
			mock: s3mock{listFn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return &s3.ListObjectsV2Output{
					Contents: []types.Object{
						{Key: aws.String("key1")},
					},
				}, nil
			}},
			want:    []string{"key1"},
			wantErr: false,
		},
		{
			name:   "error",
			prefix: "k1",
			bucket: "bucket",
			mock: s3mock{listFn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				return nil, fmt.Errorf("some err")
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()
			cli, err := New(ctx, tt.bucket)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.s3 = tt.mock

			got, err := cli.ListKeys(ctx, tt.prefix)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.want, got)
		})
	}
}
