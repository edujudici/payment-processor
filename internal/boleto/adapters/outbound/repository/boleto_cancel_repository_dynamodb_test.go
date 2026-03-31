package repository

import (
	"boleto-cancel/internal/boleto/domain"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDBClient struct {
	mock.Mock
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func TestGetBank(t *testing.T) {

	t.Run("should return error if TABLE_BANK_NAME is not set", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Unsetenv("TABLE_BANK_NAME")

		_, err = repo.GetBank(context.Background(), "123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "TABLE_BANK_NAME is not set")
	})

	t.Run("should return error if GetItem fails", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BANK_NAME", "BankTable")

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(nil, errors.New("dynamodb error"))

		_, err = repo.GetBank(context.Background(), "123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error get bank")
		mockDynamoClient.AssertExpectations(t)
	})

	t.Run("should return error if bank is not found", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BANK_NAME", "BankTable")

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: nil,
		}, nil)

		_, err = repo.GetBank(context.Background(), "123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error bank not found with Id")
		mockDynamoClient.AssertExpectations(t)
	})

	t.Run("should return error if unmarshal fails", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BANK_NAME", "BankTable")

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"bankCode": &types.AttributeValueMemberS{Value: "123"},
			},
		}, nil)

		_, err = repo.GetBank(context.Background(), "123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error to map data to bank")
		mockDynamoClient.AssertExpectations(t)
	})

	t.Run("should return bank on success", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BANK_NAME", "BankTable")

		mockBank := domain.Bank{
			Id: "123",
		}
		mockItem, _ := attributevalue.MarshalMap(mockBank)

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: mockItem,
		}, nil)

		bank, err := repo.GetBank(context.Background(), "123")
		assert.Nil(t, err)
		assert.Equal(t, &mockBank, bank)
		mockDynamoClient.AssertExpectations(t)
	})
}

func TestGetBusiness(t *testing.T) {

	t.Run("should return error if TABLE_BUSINESS_NAME is not set", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Unsetenv("TABLE_BUSINESS_NAME")

		_, err = repo.GetBusiness(context.Background(), "tx123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "TABLE_BUSINESS_NAME is not set")
	})

	t.Run("should return error if GetItem fails", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BUSINESS_NAME", "BusinessTable")

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(nil, errors.New("dynamodb error"))

		_, err = repo.GetBusiness(context.Background(), "tx123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error get business")
		mockDynamoClient.AssertExpectations(t)
	})

	t.Run("should return error if business is not found", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BUSINESS_NAME", "BusinessTable")

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: nil,
		}, nil)

		_, err = repo.GetBusiness(context.Background(), "tx123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error business not found with Id")
		mockDynamoClient.AssertExpectations(t)
	})

	t.Run("should return error if unmarshal fails", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BUSINESS_NAME", "BusinessTable")

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "tx123"},
				"name": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "not a string"},
					},
				},
			},
		}, nil)

		_, err = repo.GetBusiness(context.Background(), "tx123")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error to map data to business")
		mockDynamoClient.AssertExpectations(t)
	})

	t.Run("should return business on success", func(t *testing.T) {
		mockDynamoClient := new(MockDynamoDBClient)
		repo, err := NewDynamoDBBoletoCancelRepository(mockDynamoClient)
		assert.NoError(t, err, "Repository initialization should not return an error")

		os.Setenv("TABLE_BUSINESS_NAME", "BusinessTable")

		mockBusiness := domain.Business{
			Id:                "businessId",
			Company:           "CNH0",
			Name:              "COMPANY NAME",
			CNPJ:              "12345678000199",
			StateRegistration: "12345678000199",
		}
		mockItem, _ := attributevalue.MarshalMap(mockBusiness)

		mockDynamoClient.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
			Item: mockItem,
		}, nil)

		business, err := repo.GetBusiness(context.Background(), "tx123")
		assert.Nil(t, err)
		assert.Equal(t, &mockBusiness, business)
		mockDynamoClient.AssertExpectations(t)
	})
}
