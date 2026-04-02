package handler

// import (
// 	"context"
// 	"errors"
// 	"payment-processor/internal/payment/usecase"
// 	"testing"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type MockBoletoCancelHandler struct {
// 	mock.Mock
// }

// func (m *MockBoletoCancelHandler) Execute(ctx context.Context, request usecase.BoletoRequest) error {
// 	args := m.Called(ctx, request)
// 	return args.Error(0)
// }

// func TestHandleRequest_Provider(t *testing.T) {
// 	ctx := context.Background()

// 	tests := []struct {
// 		name        string
// 		sqsEvent    events.SQSEvent
// 		mockSetup   func(mockHandler *MockBoletoCancelHandler)
// 		expectedErr bool
// 	}{
// 		{
// 			name: "Success",
// 			sqsEvent: events.SQSEvent{
// 				Records: []events.SQSMessage{
// 					{
// 						ReceiptHandle: "test-receipt-handle",
// 						Body:          `{"endToEndId":"E12345678202009091221kkkkkkkkkkk","txid":"CNHaIm0000037465b001234240","valor":"9.99","horario":"2024-01-09T00:00:00Z","infoPagador":"0123456789","chave":"dbb56a88-9160-434c-a801-e881c8f59f9f","docCanal":"20260228235959|237|4261504261500018100|1000000000|12345678901234"}`,
// 					},
// 				},
// 			},
// 			mockSetup: func(mockHandler *MockBoletoCancelHandler) {
// 				mockHandler.On("Execute", ctx, mock.Anything).Return(nil)
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "JSON Unmarshal Error",
// 			sqsEvent: events.SQSEvent{
// 				Records: []events.SQSMessage{
// 					{
// 						ReceiptHandle: "test-receipt-handle",
// 						Body:          "invalid json", // JSON inválido
// 					},
// 				},
// 			},
// 			mockSetup:   func(mockHandler *MockBoletoCancelHandler) {},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Execute Error",
// 			sqsEvent: events.SQSEvent{
// 				Records: []events.SQSMessage{
// 					{
// 						ReceiptHandle: "test-receipt-handle",
// 						Body:          `{"endToEndId":"E12345678202009091221kkkkkkkkkkk","txid":"CNHaIm0000037465b001234240","valor":"9.99","horario":"2024-01-09T00:00:00Z","infoPagador":"0123456789","chave":"dbb56a88-9160-434c-a801-e881c8f59f9f","docCanal":"20260228235959|237|4261504261500018100|1000000000|12345678901234"}`,
// 					},
// 				},
// 			},
// 			mockSetup: func(mockHandler *MockBoletoCancelHandler) {
// 				mockHandler.On("Execute", ctx, mock.Anything).Return(errors.New("boleto cancel error"))
// 			},
// 			expectedErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockHandler := new(MockBoletoCancelHandler)
// 			tt.mockSetup(mockHandler)

// 			handler := NewPaymentProcessorHandler(mockHandler)

// 			err := handler.Process(ctx, tt.sqsEvent)

// 			if tt.expectedErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			mockHandler.AssertExpectations(t)
// 		})
// 	}
// }
