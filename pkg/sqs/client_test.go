package sqs

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/require"
)

// nolint: paralleltest
func TestNew(t *testing.T) {
	opt := WithEndpoint("test", true)

	got, err := New(context.TODO(), "name", opt)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "name", got.bucketName)

	// make AWS lib to return an error
	t.Setenv("AWS_ENABLE_ENDPOINT_DISCOVERY", "ERROR")

	got, err = New(context.TODO(), "name")
	require.Error(t, err)
	require.Nil(t, got)
}

type sqsmock struct {
	delFn  func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	getFn  func(ctx context.Context, params *sqs.GetMessageInput, optFns ...func(*sqs.Options)) (*sqs.GetMessageOutput, error)
	listFn func(ctx context.Context, params *sqs.ListMessagesV2Input, optFns ...func(*sqs.Options)) (*sqs.ListMessagesV2Output, error)
	putFn  func(ctx context.Context, params *sqs.PutMessageInput, optFns ...func(*sqs.Options)) (*sqs.PutMessageOutput, error)
}

func (s sqsmock) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return s.delFn(ctx, params, optFns...)
}

func (s sqsmock) GetMessage(ctx context.Context, params *sqs.GetMessageInput, optFns ...func(*sqs.Options)) (*sqs.GetMessageOutput, error) {
	return s.getFn(ctx, params, optFns...)
}

func (s sqsmock) ListMessagesV2(ctx context.Context, params *sqs.ListMessagesV2Input, optFns ...func(*sqs.Options)) (*sqs.ListMessagesV2Output, error) {
	return s.listFn(ctx, params, optFns...)
}

func (s sqsmock) PutMessage(ctx context.Context, params *sqs.PutMessageInput, optFns ...func(*sqs.Options)) (*sqs.PutMessageOutput, error) {
	return s.putFn(ctx, params, optFns...)
}

func TestSQSClient_DeleteMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		bucket  string
		mock    SQS
		wantErr bool
	}{
		{
			name:   "success",
			key:    "k1",
			bucket: "bucket",
			mock: sqsmock{delFn: func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
				return &sqs.DeleteMessageOutput{}, nil
			}},
			wantErr: false,
		},
		{
			name:   "error",
			key:    "k1",
			bucket: "bucket",
			mock: sqsmock{delFn: func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
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

			cli.sqs = tt.mock

			err = cli.Delete(ctx, tt.key)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSQSClient_GetMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		bucket  string
		mock    SQS
		want    *Message
		wantErr bool
	}{

		{
			name:   "success",
			key:    "k1",
			bucket: "bucket",
			mock: sqsmock{getFn: func(ctx context.Context, params *sqs.GetMessageInput, optFns ...func(*sqs.Options)) (*sqs.GetMessageOutput, error) {
				return &sqs.GetMessageOutput{
					Body: io.NopCloser(strings.NewReader("test str")),
				}, nil
			}},
			want: &Message{
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
			mock: sqsmock{getFn: func(ctx context.Context, params *sqs.GetMessageInput, optFns ...func(*sqs.Options)) (*sqs.GetMessageOutput, error) {
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

			cli.sqs = tt.mock

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

func TestSQSClient_ListMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prefix  string
		bucket  string
		mock    SQS
		want    []string
		wantErr bool
	}{
		{
			name:   "success - all",
			prefix: "",
			bucket: "bucket",
			mock: sqsmock{listFn: func(ctx context.Context, params *sqs.ListMessagesV2Input, optFns ...func(*sqs.Options)) (*sqs.ListMessagesV2Output, error) {
				return &sqs.ListMessagesV2Output{
					Contents: []types.Message{
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
			mock: sqsmock{listFn: func(ctx context.Context, params *sqs.ListMessagesV2Input, optFns ...func(*sqs.Options)) (*sqs.ListMessagesV2Output, error) {
				return &sqs.ListMessagesV2Output{
					Contents: []types.Message{
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
			mock: sqsmock{listFn: func(ctx context.Context, params *sqs.ListMessagesV2Input, optFns ...func(*sqs.Options)) (*sqs.ListMessagesV2Output, error) {
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

			cli.sqs = tt.mock

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

func TestSQSClient_PutMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		bucket  string
		mock    SQS
		wantErr bool
	}{
		{
			name:   "success",
			key:    "k1",
			bucket: "bucket",
			mock: sqsmock{putFn: func(ctx context.Context, params *sqs.PutMessageInput, optFns ...func(*sqs.Options)) (*sqs.PutMessageOutput, error) {
				return &sqs.PutMessageOutput{}, nil
			}},
			wantErr: false,
		},
		{
			name:   "error",
			key:    "k1",
			bucket: "bucket",
			mock: sqsmock{putFn: func(ctx context.Context, params *sqs.PutMessageInput, optFns ...func(*sqs.Options)) (*sqs.PutMessageOutput, error) {
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

			cli.sqs = tt.mock

			err = cli.Put(ctx, tt.key, nil)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
