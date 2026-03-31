package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQSHandler struct {
	mock.Mock
}

func (m *MockSQSHandler) SendMessageToSQS(ctx context.Context, queueURL string, messageBody string) error {
	args := m.Called(ctx, queueURL, messageBody)
	return args.Error(0)
}

func TestSendMessageSuccess(t *testing.T) {
	queueURL := "https://example.com/queue"

	os.Setenv("QUEUE_URL", queueURL)

	ctx := context.TODO()
	messageBody := "Test Message"

	mockSQSHandler := new(MockSQSHandler)
	mockSQSHandler.On("SendMessageToSQS", ctx, queueURL, messageBody).Return(nil)

	useCase := NewSendMessageService(mockSQSHandler)

	err := useCase.Send(ctx, messageBody)

	assert.NoError(t, err)

	mockSQSHandler.AssertExpectations(t)
}

func TestMissingEnvironmentVariable(t *testing.T) {
	os.Unsetenv("QUEUE_URL")

	ctx := context.TODO()
	messageBody := "Test Message"

	mockSQSHandler := new(MockSQSHandler)

	useCase := NewSendMessageService(mockSQSHandler)

	err := useCase.Send(ctx, messageBody)

	assert.Error(t, err)
	assert.EqualError(t, err, "QUEUE_URL is not set")

	mockSQSHandler.AssertExpectations(t)
}

func TestSendMessageFailure(t *testing.T) {
	queueURL := "https://example.com/queue"

	os.Setenv("QUEUE_URL", queueURL)

	ctx := context.TODO()
	messageBody := "Test Message"

	mockSQSHandler := new(MockSQSHandler)
	mockSQSHandler.On("SendMessageToSQS", ctx, queueURL, messageBody).Return(fmt.Errorf("SQS error"))

	useCase := NewSendMessageService(mockSQSHandler)

	err := useCase.Send(ctx, messageBody)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to send message: SQS error")

	mockSQSHandler.AssertExpectations(t)
}
