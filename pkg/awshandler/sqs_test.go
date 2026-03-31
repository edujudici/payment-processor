package awshandler

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQSClient struct {
	mock.Mock
}

func (m *MockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	args := m.Called(ctx, params)
	return nil, args.Error(1)
}

func TestNewSQSHandler(t *testing.T) {
	t.Run("successfully creates SQSHandler", func(t *testing.T) {
		os.Setenv("REGION", "sa-east-1")
		ctx := context.TODO()

		handler, err := NewSQSHandler(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})

	t.Run("fails when REGION is not set", func(t *testing.T) {
		os.Unsetenv("REGION")
		ctx := context.TODO()

		_, err := NewSQSHandler(ctx)

		assert.Error(t, err)
		assert.Equal(t, "required environment variables not found to instance SQS", err.Error())
	})
}

func TestSQSHandlerSendMessage(t *testing.T) {

	ctx := context.TODO()
	queueURL := "https://example.com/queue"
	messageBody := "Test Message"

	t.Run("successfully sends message", func(t *testing.T) {
		os.Setenv("REGION", "sa-east-1")

		mockSQSClient := new(MockSQSClient)
		mockSQSClient.On("SendMessage", ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(messageBody),
		}).Return(&sqs.SendMessageOutput{}, nil)

		handler := &SQSHandler{svc: mockSQSClient}
		err := handler.SendMessageToSQS(ctx, queueURL, messageBody)

		assert.NoError(t, err)

		mockSQSClient.AssertExpectations(t)
	})

	t.Run("fails to send message due to SQS error", func(t *testing.T) {
		os.Setenv("REGION", "sa-east-1")

		mockSQSClient := new(MockSQSClient)
		mockSQSClient.On("SendMessage", ctx, &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(messageBody),
		}).Return(nil, fmt.Errorf("SQS error"))

		handler := &SQSHandler{svc: mockSQSClient}
		err := handler.SendMessageToSQS(ctx, queueURL, messageBody)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to send message: SQS error")

		mockSQSClient.AssertExpectations(t)
	})
}
