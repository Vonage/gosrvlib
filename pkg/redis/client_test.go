package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	libredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	srvOpts := &SrvOptions{
		Addr:     "test.redis.invalid:6379",
		Username: "test_user",
		Password: "test_password",
		DB:       0,
	}

	got, err := New(
		context.TODO(),
		srvOpts,
		WithMessageEncodeFunc(nil),
	)

	require.Error(t, err)
	require.Nil(t, got)

	got, err = New(
		context.TODO(),
		srvOpts,
	)

	require.NoError(t, err)
	require.NotNil(t, got)
}

type redisClientMock struct {
	closeFn     func() error
	delFn       func(ctx context.Context, keys ...string) *libredis.IntCmd
	getFn       func(ctx context.Context, key string) *libredis.StringCmd
	pingFn      func(ctx context.Context) *libredis.StatusCmd
	publishFn   func(ctx context.Context, channel string, message any) *libredis.IntCmd
	setFn       func(ctx context.Context, key string, value any, expiration time.Duration) *libredis.StatusCmd
	subscribeFn func(ctx context.Context, channels ...string) *libredis.PubSub
}

func (m redisClientMock) Close() error {
	return m.closeFn()
}

func (m redisClientMock) Del(ctx context.Context, keys ...string) *libredis.IntCmd {
	return m.delFn(ctx, keys...)
}

func (m redisClientMock) Get(ctx context.Context, key string) *libredis.StringCmd {
	return m.getFn(ctx, key)
}

func (m redisClientMock) Ping(ctx context.Context) *libredis.StatusCmd {
	return m.pingFn(ctx)
}

func (m redisClientMock) Publish(ctx context.Context, channel string, message any) *libredis.IntCmd {
	return m.publishFn(ctx, channel, message)
}

func (m redisClientMock) Set(ctx context.Context, key string, value any, expiration time.Duration) *libredis.StatusCmd {
	return m.setFn(ctx, key, value, expiration)
}

func (m redisClientMock) Subscribe(ctx context.Context, channels ...string) *libredis.PubSub {
	return m.subscribeFn(ctx, channels...)
}

func TestClose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{closeFn: func() error {
				return nil
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{closeFn: func() error {
				return errors.New("test error")
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			err = cli.Close()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{setFn: func(_ context.Context, _ string, _ any, _ time.Duration) *libredis.StatusCmd {
				return libredis.NewStatusResult("", nil)
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{setFn: func(_ context.Context, _ string, _ any, _ time.Duration) *libredis.StatusCmd {
				return libredis.NewStatusResult("", errors.New("test error"))
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			err = cli.Set(ctx, "key_1", "value_1", time.Second)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{getFn: func(_ context.Context, _ string) *libredis.StringCmd {
				return libredis.NewStringResult("value_2", nil)
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{getFn: func(_ context.Context, _ string) *libredis.StringCmd {
				return libredis.NewStringResult("", errors.New("test error"))
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			got, err := cli.Get(ctx, "key_2")
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "value_2", got)
		})
	}
}

func TestDel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{delFn: func(_ context.Context, _ ...string) *libredis.IntCmd {
				return libredis.NewIntResult(0, nil)
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{delFn: func(_ context.Context, _ ...string) *libredis.IntCmd {
				return libredis.NewIntResult(0, errors.New("test error"))
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			err = cli.Del(ctx, "key_3")
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{publishFn: func(_ context.Context, _ string, _ any) *libredis.IntCmd {
				return libredis.NewIntResult(0, nil)
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{publishFn: func(_ context.Context, _ string, _ any) *libredis.IntCmd {
				return libredis.NewIntResult(0, errors.New("test error"))
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			err = cli.Send(ctx, "channel_1", "message_1")
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestSetData(t *testing.T) {
	t.Parallel()

	srvOpts := &SrvOptions{
		Addr:     "test.redis.invalid:6379",
		Username: "test_user",
		Password: "test_password",
		DB:       0,
	}

	ctx := context.TODO()
	cli, err := New(ctx, srvOpts)
	require.NoError(t, err)
	require.NotNil(t, cli)

	cli.rdb = redisClientMock{setFn: func(_ context.Context, _ string, _ any, _ time.Duration) *libredis.StatusCmd {
		return libredis.NewStatusResult("", nil)
	}}

	type TestData struct {
		Alpha string
		Beta  int
	}

	err = cli.SetData(ctx, "key_4", TestData{Alpha: "abc123", Beta: -567}, time.Second)
	require.NoError(t, err)

	err = cli.SetData(ctx, "key_5", nil, time.Second)
	require.Error(t, err)
}

func TestGetData(t *testing.T) {
	t.Parallel()

	type TestData struct {
		Alpha string
		Beta  int
	}

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{getFn: func(_ context.Context, _ string) *libredis.StringCmd {
				return libredis.NewStringResult("Kf+BAwEBCFRlc3REYXRhAf+CAAECAQVBbHBoYQEMAAEEQmV0YQEEAAAAD/+CAQZhYmMxMjMB/gLtAA==", nil)
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{getFn: func(_ context.Context, _ string) *libredis.StringCmd {
				return libredis.NewStringResult("", errors.New("test error"))
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			var data TestData

			err = cli.GetData(ctx, "key_7", &data)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "abc123", data.Alpha)
			require.Equal(t, -375, data.Beta)
		})
	}
}

func TestSendData(t *testing.T) {
	t.Parallel()

	srvOpts := &SrvOptions{
		Addr:     "test.redis.invalid:6379",
		Username: "test_user",
		Password: "test_password",
		DB:       0,
	}

	ctx := context.TODO()
	cli, err := New(ctx, srvOpts)
	require.NoError(t, err)
	require.NotNil(t, cli)

	cli.rdb = redisClientMock{publishFn: func(_ context.Context, _ string, _ any) *libredis.IntCmd {
		return libredis.NewIntResult(0, nil)
	}}

	type TestData struct {
		Alpha string
		Beta  int
	}

	err = cli.SendData(ctx, "channel_2", TestData{Alpha: "abc345", Beta: -678})
	require.NoError(t, err)

	err = cli.SendData(ctx, "channel_3", nil)
	require.Error(t, err)
}

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    RClient
		wantErr bool
	}{
		{
			name: "success",
			mock: redisClientMock{pingFn: func(_ context.Context) *libredis.StatusCmd {
				return libredis.NewStatusResult("", nil)
			}},
			wantErr: false,
		},
		{
			name: "error",
			mock: redisClientMock{pingFn: func(_ context.Context) *libredis.StatusCmd {
				return libredis.NewStatusResult("", errors.New("test error"))
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srvOpts := &SrvOptions{
				Addr:     "test.redis.invalid:6379",
				Username: "test_user",
				Password: "test_password",
				DB:       0,
			}

			ctx := context.TODO()
			cli, err := New(ctx, srvOpts)
			require.NoError(t, err)
			require.NotNil(t, cli)

			cli.rdb = tt.mock

			err = cli.HealthCheck(ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
