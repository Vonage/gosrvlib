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
}
