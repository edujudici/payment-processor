package repository

import (
	"boleto-cancel/internal/boleto/domain"
	"boleto-cancel/internal/boleto/ports"
	"boleto-cancel/pkg/awshandler"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoDBBoletoCancelRepository struct {
	db awshandler.DynamoDBClient
}

func NewDynamoDBBoletoCancelRepository(dynamoClient awshandler.DynamoDBClient) (ports.BoletoCancelRepository, error) {
	return &dynamoDBBoletoCancelRepository{
		db: dynamoClient,
	}, nil
}

// GetPayment get payment data.
func (r *dynamoDBBoletoCancelRepository) GetBank(ctx context.Context, bankId string) (*domain.Bank, error) {
	bankTable := os.Getenv("TABLE_BANK_NAME")
	if bankTable == "" {
		return nil, fmt.Errorf("TABLE_BANK_NAME is not set")
	}

	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &bankTable,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: bankId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error get bank: %w", err)
	}
	if result.Item == nil {
		return nil, fmt.Errorf("error bank not found with Id: %s", bankId)
	}

	var bank domain.Bank
	err = attributevalue.UnmarshalMap(result.Item, &bank)
	if err != nil {
		return nil, fmt.Errorf("error to map data to bank: %w", err)
	}

	return &bank, nil
}

// GetBankByAgreement get bank data by agreement.
func (r *dynamoDBBoletoCancelRepository) GetBankByAgreement(ctx context.Context, bankAgreement string) (*domain.Bank, error) {
	bankTable := os.Getenv("TABLE_BANK_NAME")
	if bankTable == "" {
		return nil, fmt.Errorf("TABLE_BANK_NAME is not set")
	}

	result, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              &bankTable,
		IndexName:              aws.String("BankAgreementIndex"),
		KeyConditionExpression: aws.String("bankAgreement = :bankAgreement"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":bankAgreement": &types.AttributeValueMemberS{Value: bankAgreement},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, fmt.Errorf("error query bank: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("bank not found with bankAgreement: %s", bankAgreement)
	}

	var bank domain.Bank
	err = attributevalue.UnmarshalMap(result.Items[0], &bank)
	if err != nil {
		return nil, fmt.Errorf("error to map data to bank: %w", err)
	}

	return &bank, nil
}

// GetBusiness get business data.
func (r *dynamoDBBoletoCancelRepository) GetBusiness(ctx context.Context, businessId string) (*domain.Business, error) {
	businessTable := os.Getenv("TABLE_BUSINESS_NAME")
	if businessTable == "" {
		return nil, fmt.Errorf("TABLE_BUSINESS_NAME is not set")
	}

	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &businessTable,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: businessId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error get business: %w", err)
	}
	if result.Item == nil {
		return nil, fmt.Errorf("error business not found with Id: %s", businessId)
	}

	var business domain.Business
	err = attributevalue.UnmarshalMap(result.Item, &business)
	if err != nil {
		return nil, fmt.Errorf("error to map data to business: %w", err)
	}

	return &business, nil
}
