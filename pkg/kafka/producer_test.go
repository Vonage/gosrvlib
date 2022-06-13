package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestProducer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                       string
		urls                       []string
		topic                      string
		options                    []Option
		expectedTimeout            time.Duration
		expectedProduceChannelSize int
		expectErr                  bool
	}{
		{
			name: "success",
			urls: []string{"url1", "url2"},
			options: []Option{
				WithSessionTimeout(time.Second * 20),
			},
			expectedTimeout: time.Second * 20,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			producer, err := NewProducer(tt.urls, tt.topic, tt.options...)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NotNil(t, producer)
				require.Nil(t, err)

				require.Equal(t, tt.expectedTimeout, producer.cfg.sessionTimeout)

				err := producer.Close()
				require.NoError(t, err)
			}
		})
	}
}

type mockProducerClient struct{}

func (m mockProducerClient) WriteMessages(ctx context.Context, msg ...kafka.Message) error {
	return nil
}

func (m mockProducerClient) Close() error {
	return nil
}

func TestProduceMessageError(t *testing.T) {
	t.Parallel()

	producer, err := NewProducer([]string{"url"}, "test")
	require.Nil(t, err, "NewProducer() unexpected error = %v", err)

	err = producer.ProduceMessage(context.TODO(), nil)
	require.Error(t, err)

	producer.client = &mockProducerClient{}
	err = producer.ProduceMessage(context.TODO(), nil)
	require.NoError(t, err)
}
