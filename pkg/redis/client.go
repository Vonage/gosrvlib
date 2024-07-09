package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Vonage/gosrvlib/pkg/encode"
	libredis "github.com/redis/go-redis/v9"
)

// TEncodeFunc is the type of function used to replace the default message encoding function used by SendData().
type TEncodeFunc func(ctx context.Context, data any) (string, error)

// TDecodeFunc is the type of function used to replace the default message decoding function used by ReceiveData().
type TDecodeFunc func(ctx context.Context, msg string, data any) error

// SrvOptions is an alias for the parent service client options.
type SrvOptions = libredis.Options

// RClient represents the mockable functions in the parent Redis Client.
type RClient interface {
	Close() error
	Del(ctx context.Context, keys ...string) *libredis.IntCmd
	Get(ctx context.Context, key string) *libredis.StringCmd
	Ping(ctx context.Context) *libredis.StatusCmd // this function is used by the HealthCheck
	Publish(ctx context.Context, channel string, message any) *libredis.IntCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *libredis.StatusCmd
	Subscribe(ctx context.Context, channels ...string) *libredis.PubSub
}

// Client is a wrapper for the Redis Client.
type Client struct {
	// rdb is the interface for the upstream Client functions.
	rdb RClient

	// messageEncodeFunc is the function used by SendData()
	// to encode and serialize the input data to a string compatible with Redis.
	messageEncodeFunc TEncodeFunc

	// messageDecodeFunc is the function used by ReceiveData()
	// to decode a message encoded with messageEncodeFunc to the provided data object.
	// The value underlying data must be a pointer to the correct type for the next data item received.
	messageDecodeFunc TDecodeFunc
}

// New creates a new instance of the Redis client wrapper.
func New(ctx context.Context, srvopt *SrvOptions, opts ...Option) (*Client, error) {
	cfg, err := loadConfig(ctx, srvopt, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new redis client: %w", err)
	}

	return &Client{
		rdb:               libredis.NewClient(cfg.srvOpts),
		messageEncodeFunc: cfg.messageEncodeFunc,
		messageDecodeFunc: cfg.messageDecodeFunc,
	}, nil
}

// Close closes the client, releasing any open resources.
func (c *Client) Close() error {
	err := c.rdb.Close()
	if err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}

	return nil
}

// Set a raw value for the specified key with an expiration time.
// Zero expiration means the key has no expiration time.
func (c *Client) Set(ctx context.Context, key string, value any, exp time.Duration) error {
	err := c.rdb.Set(ctx, key, value, exp).Err()
	if err != nil {
		return fmt.Errorf("cannot set key %s: %w", key, err)
	}

	return nil
}

// Get retrieves the raw value of the specified key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("cannot retrieve key %s: %w", key, err)
	}

	return val, nil
}

// Del deletes the specified key from the datastore.
func (c *Client) Del(ctx context.Context, key string) error {
	err := c.rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("cannot delete key: %s %w", key, err)
	}

	return nil
}

// Send publish a raw value to the specified channel.
func (c *Client) Send(ctx context.Context, channel string, message any) error {
	err := c.rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("cannot send message to %s channel: %w", channel, err)
	}

	return nil
}

// MessageEncode encodes and serialize the input data to a string.
func MessageEncode(data any) (string, error) {
	return encode.Encode(data) //nolint:wrapcheck
}

// MessageDecode decodes a message encoded with MessageEncode to the provided data object.
// The value underlying data must be a pointer to the correct type for the next data item received.
func MessageDecode(msg string, data any) error {
	return encode.Decode(msg, data) //nolint:wrapcheck
}

// DefaultMessageEncodeFunc is the default function to encode and serialize the input data for SendData().
func DefaultMessageEncodeFunc(_ context.Context, data any) (string, error) {
	return MessageEncode(data)
}

// DefaultMessageDecodeFunc is the default function to decode a message for ReceiveData().
// The value underlying data must be a pointer to the correct type for the next data item received.
func DefaultMessageDecodeFunc(_ context.Context, msg string, data any) error {
	return MessageDecode(msg, data)
}

// SetData sets an encoded value for the specified key with an expiration time.
// Zero expiration means the key has no expiration time.
func (c *Client) SetData(ctx context.Context, key string, data any, exp time.Duration) error {
	value, err := c.messageEncodeFunc(ctx, data)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, value, exp)
}

// GetData retrieves an encoded value of the specified key.
func (c *Client) GetData(ctx context.Context, key string, data any) error {
	value, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	return c.messageDecodeFunc(ctx, value, data)
}

// SendData publish an encoded value to the specified channel.
func (c *Client) SendData(ctx context.Context, channel string, data any) error {
	message, err := c.messageEncodeFunc(ctx, data)
	if err != nil {
		return err
	}

	return c.Send(ctx, channel, message)
}

// HealthCheck checks if the current data-store is alive.
func (c *Client) HealthCheck(ctx context.Context) error {
	err := c.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("unable to connect to Redis: %w", err)
	}

	return nil
}
