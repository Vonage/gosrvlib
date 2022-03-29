package sqs

import (
	"context"
	"fmt"
	"io"

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
	sqs         SQS
	
	// queueUrl is the URL of the Amazon SQS queue to which a message is sent.
	// Queue URLs and names are case-sensitive.
	queueUrl string
	
	// messageGroupId is a tag that specifies that a message belongs to a specific message group.
	messageGroupId string
	
	// waitTimeSeconds is the duration (in seconds) for which the call waits for a message to arrive in the queue before returning.
	waitTimeSeconds int32
}

// New creates a new instance of the SQS client wrapper.
func New(ctx context.Context, bucketName string, opts ...Option) (*Client, error) {
	cfg, err := loadConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new sqs client: %w", err)
	}

	return &Client{
		sqs:         sqs.NewFromConfig(cfg),
		bucketName: bucketName,
	}, nil
}

// Message represents a message in the queue.
type Message struct {
	// The message's contents (not URL-encoded).
	Body string
	
	// A unique identifier for the message.
	MessageId string
}


// Send delivers a message to the queue.
func (c *Client) Send(ctx context.Context, message string) error {
	_, err := c.sqs.SendMessage(
	ctx, 
	&sqs.SendMessageInput{
			QueueUrl: c.queueUrl,
			MessageGroupId: c.messageGroupId,
			MessageBody: message,
		})
	if err != nil {
		return fmt.Errorf("cannot send message to the queue: %w", err)
	}

	return nil
}

// Receive retrieves one messages from the queue.
func (c *Client) Receive(ctx context.Context, key string) (*Message, error) {
	resp, err := c.sqs.ReceiveMessage(
		ctx, 
		&sqs.ReceiveMessageInput{
			QueueUrl: c.queueUrl,
			MessageGroupId: c.messageGroupId,
			WaitTimeSeconds: c.waitTimeSeconds,
		})
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve message from the queue: %w", err)
	}

	return &Message{
		Body: resp.Body,
		MessageId: resp.MessageId,
	}, nil
}



// Delete deletes the specified message from the queue.
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.sqs.DeleteMessage(
		ctx, 
		&sqs.DeleteMessageInput{
			QueueUrl: c.queueUrl,
			ReceiptHandle: ,
		})
	if err != nil {
		return fmt.Errorf("cannot delete message from the queue: %w", err)
	}

	return nil
}






