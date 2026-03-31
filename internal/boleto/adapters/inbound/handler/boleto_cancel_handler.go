package handler

import (
	"boleto-cancel/internal/boleto/adapters/inbound/dto"
	"boleto-cancel/internal/boleto/usecase"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	boletoCancelUseCase usecase.BoletoCancelInterface
}

func NewHandler(
	boletoCancelUseCase usecase.BoletoCancelInterface,
) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	handler := &Handler{
		boletoCancelUseCase: boletoCancelUseCase,
	}
	return handler.HandleRequest
}

func (h *Handler) HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {

		receiptHandle := message.ReceiptHandle
		log.Printf("processing message with receiptHandle: %s\n", receiptHandle)

		log.Println("Received message =============================:", message.Body)

		// transform json data to object
		var boletoRequest dto.Request
		if err := json.Unmarshal([]byte(message.Body), &boletoRequest); err != nil {
			return fmt.Errorf("error parsing JSON request: %v", err)
		}

		input := boletoRequest.ToUsecaseInput()

		// boleto cancel use case
		err := h.boletoCancelUseCase.Execute(ctx, input)
		if err != nil {
			return fmt.Errorf("error boleto cancel: %v", err)
		}
	}
	return nil
}
