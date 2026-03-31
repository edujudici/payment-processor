package usecase

import (
	"boleto-cancel/internal/boleto/adapters/outbound/service"
	"boleto-cancel/internal/boleto/domain"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuth struct {
	mock.Mock
}

type MockBoletoCancelService struct {
	mock.Mock
}

type MockSendMessageService struct {
	mock.Mock
}

type MockBoletoCancelRepository struct {
	mock.Mock
}

func (m *MockAuth) GetAccessToken() (*string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockBoletoCancelService) BoletoCancel(token string, req *service.BoletoCancelInput) (int, *service.BoletoCancelOutput, error) {
	args := m.Called(token, req)

	if args.Get(1) == nil {
		return args.Int(0), nil, args.Error(2)
	}

	return args.Int(0), args.Get(1).(*service.BoletoCancelOutput), args.Error(2)
}

func (m *MockSendMessageService) Send(ctx context.Context, messageBody string) error {
	args := m.Called(ctx, messageBody)
	return args.Error(0)
}

func (m *MockBoletoCancelRepository) GetBank(ctx context.Context, pixKey string) (*domain.Bank, error) {
	args := m.Called(ctx, pixKey)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Bank), args.Error(1)
}

func (m *MockBoletoCancelRepository) GetBankByAgreement(ctx context.Context, bankAgreement string) (*domain.Bank, error) {
	args := m.Called(ctx, bankAgreement)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Bank), args.Error(1)
}

func (m *MockBoletoCancelRepository) GetBusiness(ctx context.Context, businessId string) (*domain.Business, error) {
	args := m.Called(ctx, businessId)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Business), args.Error(1)
}

func TestBoletoCancel(t *testing.T) {
	tests := []struct {
		name           string
		input          BoletoRequest
		mockSetup      func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository)
		expectedCode   interface{}
		expectedOutput service.BoletoCancelOutput
		expectedError  error
	}{
		{
			name: "success",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "20260228235959|237|1234567|1000000000|12345678901234",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(ptr("valid_token"), nil)
				mockBoletoCancelService.On("BoletoCancel", "valid_token", mock.Anything).Return(http.StatusOK, &service.BoletoCancelOutput{
					Registrado: true,
				}, nil)
				mockBoletoCancelRepository.On("GetBankByAgreement", mock.Anything, "1234567").Return(&domain.Bank{
					Id:                      "CNH0MOTNIPIX",
					Company:                 "CNH0",
					BankAccountAbbreviation: "B",
					BankCode:                237,
					BankAgreement:           "12345",
					PixKey:                  "3ac72fa6-5061-4940-87e7-b15175d9377b",
					ModelOfGood:             "MOT",
					OriginCode:              "N",
					NameBank:                "BRADESCO",
					Agency:                  456,
					Account:                 789,
					DigitAccount:            0,
					PaymentMethod:           "PIX",
				}, nil)
				mockBoletoCancelRepository.On("GetBusiness", mock.Anything, "CNH0").Return(&domain.Business{
					Id:                "businessId",
					Company:           "CNH0",
					Name:              "COMPANY NAME",
					CNPJ:              "12345678000199",
					StateRegistration: "12345678000199",
				}, nil)
				mockSendMessageService.On("Send", mock.Anything, mock.Anything).Return(nil)

			},
			expectedCode: http.StatusOK,
			expectedOutput: service.BoletoCancelOutput{
				Registrado: true,
			},
			expectedError: nil,
		},
		{
			name: "error getting access token",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "20260228235959|237|1234567|1000000000|12345678901234",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(nil, errors.New("auth error"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedOutput: service.BoletoCancelOutput{},
			expectedError:  errors.New("auth error"),
		},
		{
			name: "error get bank data",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "20260228235959|237|1234567|1000000000|12345678901234",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(ptr("valid_token"), nil)
				mockBoletoCancelRepository.On("GetBankByAgreement", mock.Anything, "1234567").Return(nil, errors.New("bank not found"))
			},
			expectedError: errors.New("bank not found"),
		},
		{
			name: "error get business data",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "20260228235959|237|1234567|1000000000|12345678901234",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(ptr("valid_token"), nil)
				mockBoletoCancelRepository.On("GetBankByAgreement", mock.Anything, "1234567").Return(&domain.Bank{
					Id:                      "CNH0MOTNIPIX",
					Company:                 "CNH0",
					BankAccountAbbreviation: "B",
					BankCode:                237,
					BankAgreement:           "12345",
					PixKey:                  "3ac72fa6-5061-4940-87e7-b15175d9377b",
					ModelOfGood:             "MOT",
					OriginCode:              "N",
					NameBank:                "BRADESCO",
					Agency:                  456,
					Account:                 789,
					DigitAccount:            0,
					PaymentMethod:           "PIX",
				}, nil)
				mockBoletoCancelRepository.On("GetBusiness", mock.Anything, "CNH0").Return(nil, errors.New("business not found"))
			},
			expectedError: errors.New("business not found"),
		},
		{
			name: "error invalid DocCanal format",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "invalid_format",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(ptr("valid_token"), nil)
			},
			expectedError: errors.New("invalid text, does not contain all expected fields"),
		},
		{
			name: "error calceling boleto",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "20260228235959|237|1234567|1000000000|12345678901234",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(ptr("valid_token"), nil)
				mockBoletoCancelService.On("BoletoCancel", "valid_token", mock.Anything).Return(http.StatusInternalServerError, nil, errors.New("boleto error"))
				mockBoletoCancelRepository.On("GetBankByAgreement", mock.Anything, "1234567").Return(&domain.Bank{
					Id:                      "CNH0MOTNIPIX",
					Company:                 "CNH0",
					BankAccountAbbreviation: "B",
					BankCode:                237,
					BankAgreement:           "12345",
					PixKey:                  "3ac72fa6-5061-4940-87e7-b15175d9377b",
					ModelOfGood:             "MOT",
					OriginCode:              "N",
					NameBank:                "BRADESCO",
					Agency:                  456,
					Account:                 789,
					DigitAccount:            0,
					PaymentMethod:           "PIX",
				}, nil)
				mockBoletoCancelRepository.On("GetBusiness", mock.Anything, "CNH0").Return(&domain.Business{
					Id:                "businessId",
					Company:           "CNH0",
					Name:              "COMPANY NAME",
					CNPJ:              "12345678000199",
					StateRegistration: "12345678000199",
				}, nil)
			},
			expectedCode:   http.StatusInternalServerError,
			expectedOutput: service.BoletoCancelOutput{},
			expectedError:  errors.New("boleto error"),
		},
		{
			name: "error sending message",
			input: BoletoRequest{
				TxId:     "CNHaIm0000037875p0051525B4",
				DocCanal: "20260228235959|237|1234567|1000000000|12345678901234",
			},
			mockSetup: func(mockAuth *MockAuth, mockBoletoCancelService *MockBoletoCancelService, mockSendMessageService *MockSendMessageService, mockBoletoCancelRepository *MockBoletoCancelRepository) {
				mockAuth.On("GetAccessToken").Return(ptr("valid_token"), nil)
				mockBoletoCancelService.On("BoletoCancel", "valid_token", mock.Anything).Return(http.StatusOK, &service.BoletoCancelOutput{
					Registrado: true,
				}, nil)
				mockBoletoCancelRepository.On("GetBankByAgreement", mock.Anything, "1234567").Return(&domain.Bank{
					Id:                      "CNH0MOTNIPIX",
					Company:                 "CNH0",
					BankAccountAbbreviation: "B",
					BankCode:                237,
					BankAgreement:           "12345",
					PixKey:                  "3ac72fa6-5061-4940-87e7-b15175d9377b",
					ModelOfGood:             "MOT",
					OriginCode:              "N",
					NameBank:                "BRADESCO",
					Agency:                  456,
					Account:                 789,
					DigitAccount:            0,
					PaymentMethod:           "PIX",
				}, nil)
				mockBoletoCancelRepository.On("GetBusiness", mock.Anything, "CNH0").Return(&domain.Business{
					Id:                "businessId",
					Company:           "CNH0",
					Name:              "COMPANY NAME",
					CNPJ:              "12345678000199",
					StateRegistration: "12345678000199",
				}, nil)
				mockSendMessageService.On("Send", mock.Anything, mock.Anything).Return(errors.New("error sending message"))
			},
			expectedCode:   http.StatusInternalServerError,
			expectedOutput: service.BoletoCancelOutput{},
			expectedError:  errors.New("error to send boleto cancel message: error sending message"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(MockAuth)
			mockBoletoCancelService := new(MockBoletoCancelService)
			mockSendMessageService := new(MockSendMessageService)
			mockBoletoCancelRepository := new(MockBoletoCancelRepository)

			tt.mockSetup(mockAuth, mockBoletoCancelService, mockSendMessageService, mockBoletoCancelRepository)

			useCase := NewBoletoCancelUseCase(mockAuth, mockBoletoCancelService, mockSendMessageService, mockBoletoCancelRepository)

			err := useCase.Execute(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockAuth.AssertExpectations(t)
			mockBoletoCancelService.AssertExpectations(t)
			mockBoletoCancelRepository.AssertExpectations(t)
		})
	}
}

func ptr(s string) *string {
	return &s
}
