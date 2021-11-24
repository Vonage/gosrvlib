package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	var (
		expectedTimeout               = time.Second * 10
		expectedAutoOffsetResetPolicy = OffsetLatest
	)

	consumer, err := NewConsumer(
		[]string{"url1", "url2"},
		[]string{"topic1", "topic2"},
		"group1",
		WithTimeout(time.Second*10),
		WithAutoOffsetResetPolicy(OffsetLatest),
	)

	require.Nil(t, err, "NewConsumer() unexpected error = %v", err)
	require.Equal(t, expectedTimeout, consumer.cfg.timeout)
	require.Equal(t, expectedAutoOffsetResetPolicy, consumer.cfg.autoOffsetResetPolicy)

	require.Nil(t, consumer.Close())
}

func TestConsumerError(t *testing.T) {
	consumer, err := NewConsumer(
		[]string{"url1", "url2"},
		[]string{"topic1", "topic2"},
		"group1",
		WithAutoOffsetResetPolicy("badOffset"),
	)

	require.Error(t, err, "expects error but got %v", err)
	require.Nil(t, consumer, "consumer must be nil")
}

func TestConsumerTopicSubscribeError(t *testing.T) {
	consumer, err := NewConsumer(
		[]string{"url1", "url2"},
		nil,
		"group1",
	)

	require.Error(t, err, "expects error but got %v", err)
	require.Nil(t, consumer, "consumer must be nil")
}

func TestConsumerReadMessage(t *testing.T) {
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
