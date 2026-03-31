package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthHTTPClient struct {
	mock.Mock
}

func (m *MockAuthHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type TestCaseAuth struct {
	name          string
	response      *http.Response
	err           error
	expectedToken string
	expectedError string
}

func TestGetAccessToken(t *testing.T) {
	os.Setenv("ACCESS_TOKEN_URL", "http://example.com/token")
	os.Setenv("ACCESSTAGE_USERNAME", "test_id")
	os.Setenv("ACCESSTAGE_PASSWORD", "test_secret")

	testCases := []TestCaseAuth{
		{
			name: "Success",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"AuthenticationResult": {"IdToken": "mock_token"}}`))),
			},
			err:           nil,
			expectedToken: "mock_token",
			expectedError: "",
		},
		{
			name: "RequestError",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       nil,
			},
			err:           errors.New("request error"),
			expectedToken: "",
			expectedError: "error executing request: Post \"http://example.com/token\": request error",
		},
		{
			name: "HTTPError",
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(`Internal Server Error`))),
			},
			err:           nil,
			expectedToken: "",
			expectedError: "error response token. Code: 500 and ResponseBody: Internal Server Error",
		},
		{
			name: "JSONParseError",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(`invalid json`))),
			},
			err:           nil,
			expectedToken: "",
			expectedError: "error parsing JSON: invalid character 'i' looking for beginning of value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &http.Client{
				Transport: new(MockAuthHTTPClient),
			}
			authService := NewAuthService(mockClient)
			mockClient.Transport.(*MockAuthHTTPClient).On("RoundTrip", mock.Anything).Return(tc.response, tc.err)

			response, err := authService.GetAccessToken()

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tc.expectedToken, response)
			}
			mockClient.Transport.(*MockAuthHTTPClient).AssertExpectations(t)
		})
	}
}
