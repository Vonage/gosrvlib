package sqs

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/require"
)

func Test_loadConfig(t *testing.T) {
	t.Parallel()

	var (
		wt  int32 = 13
		vt  int32 = 17
		reg       = "eu-central-1"
	)

	got, err := loadConfig(
		context.TODO(),
		WithEndpoint("test", true),
		withAWSOption(config.WithRegion(reg)),
		WithWaitTimeSeconds(wt),
		WithVisibilityTimeout(vt),
	)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, reg, got.awsConfig.Region)
	require.Equal(t, wt, got.waitTimeSeconds)
	require.Equal(t, vt, got.visibilityTimeout)

	got, err = loadConfig(
		context.TODO(),
		WithWaitTimeSeconds(-1),
	)

	require.Error(t, err)
	require.Nil(t, got)

	got, err = loadConfig(
		context.TODO(),
		WithVisibilityTimeout(-1),
	)

	require.Error(t, err)
	require.Nil(t, got)

	got, err = loadConfig(
		context.TODO(),
		WithVisibilityTimeout(43201),
	)

	require.Error(t, err)
	require.Nil(t, got)
}
