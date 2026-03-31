package main

import (
	"boleto-cancel/internal/boleto/adapters/inbound/handler"
	"boleto-cancel/internal/boleto/adapters/outbound/repository"
	"boleto-cancel/internal/boleto/adapters/outbound/service"
	"boleto-cancel/internal/boleto/usecase"
	"boleto-cancel/pkg/awshandler"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Error   string `json:"Erro"`
	Message string `json:"Mensagem"`
}

func main() {
	ctx := context.Background()

	dynamoClient, err := awshandler.NewDynamoDBClient(ctx)
	if err != nil {
		log.Fatalf("error create client DynamoDB: %v", err)
	}
	boletoCancelRepository, err := repository.NewDynamoDBBoletoCancelRepository(dynamoClient)
	if err != nil {
		log.Fatal(err)
	}
	sqsHandler, err := awshandler.NewSQSHandler(ctx)
	if err != nil {
		log.Fatal(err)
	}

	authService := service.NewAuthService(nil)
	boletoCancelService := service.NewBoletoCancelService(nil)
	sendMessageService := service.NewSendMessageService(sqsHandler)

	boletoCancelUseCase := usecase.NewBoletoCancelUseCase(authService, boletoCancelService, sendMessageService, boletoCancelRepository)

	lambda.Start(handler.NewHandler(boletoCancelUseCase))
}
