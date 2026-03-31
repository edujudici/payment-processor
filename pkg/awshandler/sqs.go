package awshandler

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSHandlerInterface interface {
	SendMessageToSQS(ctx context.Context, queueURL string, messageBody string) error
}

type SQSAPI interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type SQSHandler struct {
	svc SQSAPI
}

func NewSQSHandler(ctx context.Context) (*SQSHandler, error) {
	region := os.Getenv("REGION")
	if region == "" {
		return nil, fmt.Errorf("required environment variables not found to instance SQS")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*aws.String(region)))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := sqs.NewFromConfig(cfg)

	return &SQSHandler{svc: svc}, nil
}

func (h *SQSHandler) SendMessageToSQS(ctx context.Context, queueURL string, messageBody string) error {
	_, err := h.svc.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
