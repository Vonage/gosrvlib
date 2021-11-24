package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProducer(t *testing.T) {
	t.Parallel()

	var (
		expectedTimeout            = time.Second * 10
		expectedProduceChannelSize = 1_000
	)

	consumer, err := NewProducer(
		[]string{"url1", "url2"},
		WithTimeout(time.Second*10),
		WithProduceChannelSize(1_000),
	)

	require.Nil(t, err, "NewProducer() unexpected error = %v", err)
	require.Equal(t, expectedTimeout, consumer.cfg.timeout)
	require.Equal(t, expectedProduceChannelSize, consumer.cfg.produceChannelSize)

	consumer.Close()
}

func TestProduceMessageError(t *testing.T) {
	consumer, err := NewProducer([]string{"url"})
	require.Nil(t, err, "NewProducer() unexpected error = %v", err)

	err = consumer.ProduceMessage("", nil)
	require.Error(t, err)
}
