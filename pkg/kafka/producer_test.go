package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProducer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                       string
		urls                       []string
		options                    []Option
		expectedTimeout            time.Duration
		expectedProduceChannelSize int
		expectErr                  bool
	}{
		{
			name: "success",
			urls: []string{"url1", "url2"},
			options: []Option{
				WithTimeout(time.Second * 10),
				WithProduceChannelSize(1_000),
			},
			expectedTimeout:            time.Second * 10,
			expectedProduceChannelSize: 1_000,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			producer, err := NewProducer(tt.urls, tt.options...)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NotNil(t, producer)
				require.Nil(t, err)
				require.Equal(t, tt.expectedTimeout, producer.cfg.timeout)
				require.Equal(t, tt.expectedProduceChannelSize, producer.cfg.produceChannelSize)
				producer.Close()
			}
		})
	}
}

func TestProduceMessageError(t *testing.T) {
	t.Parallel()

	consumer, err := NewProducer([]string{"url"})
	require.Nil(t, err, "NewProducer() unexpected error = %v", err)

	err = consumer.ProduceMessage("", nil)
	require.Error(t, err)
}
