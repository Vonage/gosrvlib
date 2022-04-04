package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQS represents the mockable functions in the AWS SDK SQS client.
type SQS interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// Client is a wrapper for the SQS client in the AWS SDK.
type Client struct {
	// sqs is the interface for the upstream functions.
	sqs SQS

	// Queue URLs and names are case-sensitive and limited up to 80 chars.
	queueURL *string

	// messageGroupId is a tag that specifies that a message belongs to a specific message group.
	messageGroupID *string

	// waitTimeSeconds is the duration (in seconds) for which the call waits for a message to arrive in the queue before returning.
	waitTimeSeconds int32
}

// New creates a new instance of the SQS client wrapper.
func New(ctx context.Context, queueURL, msgGroupID string, opts ...Option) (*Client, error) {
	cfg, err := loadConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new sqs client: %w", err)
	}

	return &Client{
		sqs:            sqs.NewFromConfig(cfg),
		queueURL:       aws.String(queueURL),
		messageGroupID: aws.String(msgGroupID),
	}, nil
}

// Message represents a message in the queue.
type Message struct {
	// can contain: JSON, XML, plain text.
	Body string

	// id is the identifier used to delete the message.
	id *string
}

// Send delivers a message to the queue.
func (c *Client) Send(ctx context.Context, message string) error {
	_, err := c.sqs.SendMessage(
		ctx,
		&sqs.SendMessageInput{
			QueueUrl:       c.queueURL,
			MessageGroupId: c.messageGroupID,
			MessageBody:    aws.String(message),
		})
	if err != nil {
		return fmt.Errorf("cannot send message to the queue: %w", err)
	}

	return nil
}

// Receive retrieves a message from the queue.
func (c *Client) Receive(ctx context.Context) (*Message, error) {
	resp, err := c.sqs.ReceiveMessage(
		ctx,
		&sqs.ReceiveMessageInput{
			QueueUrl:        c.queueURL,
			WaitTimeSeconds: c.waitTimeSeconds,
		})
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve message from the queue: %w", err)
	}

	if len(resp.Messages) < 1 {
		return &Message{}, nil
	}

	return &Message{
		Body: aws.ToString(resp.Messages[0].Body),
		id:   resp.Messages[0].ReceiptHandle,
	}, nil
}

// Delete deletes the specified message from the queue.
func (c *Client) Delete(ctx context.Context, msg *Message) error {
	_, err := c.sqs.DeleteMessage(
		ctx,
		&sqs.DeleteMessageInput{
			QueueUrl:      c.queueURL,
			ReceiptHandle: msg.id,
		})
	if err != nil {
		return fmt.Errorf("cannot delete message from the queue: %w", err)
	}

	return nil
}
