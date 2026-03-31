package service

import (
	"boleto-cancel/pkg/awshandler"
	"context"
	"fmt"
	"log"
	"os"
)

type SendMessage struct {
	sqsHandler awshandler.SQSHandlerInterface
}

func NewSendMessageService(sqsHandler awshandler.SQSHandlerInterface) *SendMessage {
	return &SendMessage{
		sqsHandler: sqsHandler,
	}
}

func (sm *SendMessage) Send(ctx context.Context, messageBody string) error {
	queueURL := os.Getenv("QUEUE_URL")
	if queueURL == "" {
		return fmt.Errorf("QUEUE_URL is not set")
	}

	log.Printf("Send messageBody: %s", messageBody)

	err := sm.sqsHandler.SendMessageToSQS(ctx, queueURL, messageBody)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
