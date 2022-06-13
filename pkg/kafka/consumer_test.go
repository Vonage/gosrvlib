package kafka

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		urls      []string
		topic     string
		groupID   string
		options   []Option
		expectErr bool
	}{
		{
			name:      "success",
			urls:      []string{"url1", "url2"},
			topic:     "topic1",
			groupID:   "one",
			expectErr: false,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			consumer, err := NewConsumer(tt.urls, tt.topic, tt.groupID)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.NotNil(t, consumer, "consumerClient is nil")

				require.Nil(t, consumer.Close())
			}
		})
	}
}

type mockConsumerClient struct{}

func (m mockConsumerClient) ReadMessage(_ context.Context) (kafka.Message, error) {
	return kafka.Message{Value: []byte{1}}, nil
}

func (m mockConsumerClient) Close() error {
	return nil
}

func TestConsumerReadMessage(t *testing.T) {
	t.Parallel()

	consumer, err := NewConsumer(
		[]string{"url1", "url2"},
		"topic1",
		"group1",
	)
	require.Nil(t, err, "NewConsumer() unexpected error = %v", err)

	ctx := context.TODO()

	msg, err := consumer.ReadMessage(ctx)
	require.Error(t, err)
	require.Nil(t, msg)

	consumer.client = mockConsumerClient{}
	msg, err = consumer.ReadMessage(ctx)
	require.NoError(t, err)
	require.NotNil(t, msg)
}
