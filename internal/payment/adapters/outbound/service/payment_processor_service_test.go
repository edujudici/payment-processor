package service

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type MockBoletoHTTPClient struct {
// 	mock.Mock
// }

// func (m *MockBoletoHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
// 	args := m.Called(req)
// 	return args.Get(0).(*http.Response), args.Error(1)
// }

// type TestCase struct {
// 	name         string
// 	envSetup     func()
// 	request      string
// 	mockResponse *http.Response
// 	mockError    error
// 	expectedErr  string
// 	expectedCode interface{}
// 	expectedResp interface{}
// 	expectCall   bool // flag to indicate if RoundTrip should be called
// }

// func TestBoletoCancel(t *testing.T) {
// 	testCases := []TestCase{
// 		{
// 			name:     "Success",
// 			envSetup: func() { os.Setenv("CANCEL_BOLETO_URL", "http://example.com/boleto") },
// 			request:  `{"codBanco":"BANCO_BRASIL","cnpjClienteAccesstage":"45441789000154","numeroConvenio":3128557,"nossoNumero":"7725860","operacao":"BAIXAR"}`,
// 			mockResponse: &http.Response{
// 				StatusCode: http.StatusOK,
// 				Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"registrado":false,"dataHoraRetorno":"30072025150818"}`))),
// 			},
// 			mockError:    nil,
// 			expectedErr:  "",
// 			expectedCode: http.StatusOK,
// 			expectedResp: &BoletoCancelOutput{Registrado: false, DataHoraRetorno: "30072025150818"},
// 			expectCall:   true,
// 		},
// 		{
// 			name:         "Environment_Error",
// 			envSetup:     func() { os.Unsetenv("CANCEL_BOLETO_URL") },
// 			request:      `{"codBanco":"BANCO_BRASIL","cnpjClienteAccesstage":"45441789000154","numeroConvenio":3128557,"nossoNumero":"7725860","operacao":"BAIXAR"}`,
// 			mockResponse: nil,
// 			mockError:    nil,
// 			expectedErr:  "required environment variables not found to cancel boleto",
// 			expectedCode: http.StatusInternalServerError,
// 			expectedResp: nil,
// 			expectCall:   false,
// 		},
// 		{
// 			name:     "RequestError",
// 			envSetup: func() { os.Setenv("CANCEL_BOLETO_URL", "http://example.com/boleto") },
// 			request:  `{"codBanco":"BANCO_BRASIL","cnpjClienteAccesstage":"45441789000154","numeroConvenio":3128557,"nossoNumero":"7725860","operacao":"BAIXAR"}`,
// 			mockResponse: &http.Response{
// 				StatusCode: http.StatusInternalServerError,
// 				Body:       nil,
// 			},
// 			mockError:    errors.New("request error"),
// 			expectedErr:  "error executing request: Put \"http://example.com/boleto\": request error",
// 			expectedCode: http.StatusInternalServerError,
// 			expectedResp: nil,
// 			expectCall:   true,
// 		},
// 		{
// 			name:     "HTTPError",
// 			envSetup: func() { os.Setenv("CANCEL_BOLETO_URL", "http://example.com/boleto") },
// 			request:  `{"codBanco":"BANCO_BRASIL","cnpjClienteAccesstage":"45441789000154","numeroConvenio":3128557,"nossoNumero":"7725860","operacao":"BAIXAR"}`,
// 			mockResponse: &http.Response{
// 				StatusCode: http.StatusInternalServerError,
// 				Body:       io.NopCloser(bytes.NewBuffer([]byte(`Internal Server Error`))),
// 			},
// 			mockError:    nil,
// 			expectedErr:  "error boleto response. Code: 500 and ResponseBody: Internal Server Error",
// 			expectedCode: http.StatusInternalServerError,
// 			expectedResp: nil,
// 			expectCall:   true,
// 		},
// 		{
// 			name:     "JSONParseError",
// 			envSetup: func() { os.Setenv("CANCEL_BOLETO_URL", "http://example.com/boleto") },
// 			request:  `{"codBanco":"BANCO_BRASIL","cnpjClienteAccesstage":"45441789000154","numeroConvenio":3128557,"nossoNumero":"7725860","operacao":"BAIXAR"}`,
// 			mockResponse: &http.Response{
// 				StatusCode: http.StatusOK,
// 				Body:       io.NopCloser(bytes.NewBuffer([]byte(`invalid json`))),
// 			},
// 			mockError:    nil,
// 			expectedErr:  "error parsing JSON: invalid character 'i' looking for beginning of value",
// 			expectedCode: http.StatusInternalServerError,
// 			expectedResp: nil,
// 			expectCall:   true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.envSetup()

// 			mockClient := &MockBoletoHTTPClient{}
// 			httpClient := &http.Client{
// 				Transport: mockClient,
// 			}
// 			boletoService := NewPaymentCancelService(httpClient)

// 			if tc.expectCall {
// 				if tc.mockError != nil {
// 					mockClient.On("RoundTrip", mock.Anything).Return(tc.mockResponse, tc.mockError)
// 				} else {
// 					mockClient.On("RoundTrip", mock.Anything).Return(tc.mockResponse, nil)
// 				}
// 			}

// 			var mockRequest BoletoCancelInput
// 			err := json.Unmarshal([]byte(tc.request), &mockRequest)
// 			if err != nil {
// 				t.Fatalf("Error unmarshalling mock request: %v", err)
// 				return
// 			}

// 			statusCode, response, err := boletoService.PaymentCancel("token_valid", &mockRequest)

// 			if tc.expectedErr != "" {
// 				assert.EqualError(t, err, tc.expectedErr)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tc.expectedCode, statusCode)
// 				assert.Equal(t, tc.expectedResp, response)
// 			}

// 			if tc.expectCall {
// 				mockClient.AssertExpectations(t)
// 			} else {
// 				mockClient.AssertNotCalled(t, "RoundTrip", mock.Anything)
// 			}
// 		})
// 	}
// }
