package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		urls                          []string
		topics                        []string
		groupID                       string
		options                       []Option
		expectedTimeout               time.Duration
		expectedAutoOffsetResetPolicy Offset
		expectErr                     bool
	}{
		"success": {
			urls:    []string{"url1", "url2"},
			topics:  []string{"topic1", "topic2"},
			groupID: "one",
			options: []Option{
				WithTimeout(time.Second * 10),
				WithAutoOffsetResetPolicy(OffsetLatest),
			},
			expectedTimeout:               time.Second * 10,
			expectedAutoOffsetResetPolicy: OffsetLatest,
		},
		"bad offset": {
			urls:    []string{"url1", "url2"},
			topics:  []string{"topic1", "topic2"},
			groupID: "one",
			options: []Option{
				WithAutoOffsetResetPolicy("bad offset"),
			},
			expectErr: true,
		},
		"empty topics": {
			urls:      []string{"url1", "url2"},
			topics:    nil,
			groupID:   "one",
			expectErr: true,
		},
	}

	for name, tt := range testCases {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			consumer, err := NewConsumer(tt.urls, tt.topics, tt.groupID, tt.options...)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NotNil(t, consumer)
				require.Nil(t, err)
				require.Equal(t, tt.expectedTimeout, consumer.cfg.timeout)
				require.Equal(t, tt.expectedAutoOffsetResetPolicy, consumer.cfg.autoOffsetResetPolicy)
				require.Nil(t, consumer.Close())
			}
		})
	}
}

func TestConsumerReadMessage(t *testing.T) {
	t.Parallel()

	consumer, err := NewConsumer(
		[]string{"url1", "url2"},
		[]string{"topic1", "topic2"},
		"group1",
	)
	require.Nil(t, err, "NewConsumer() unexpected error = %v", err)

	msg, err := consumer.ReadMessage()
	require.Error(t, err)
	require.Nil(t, msg)
}
