package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Auth struct {
	HTTPClient *http.Client
}

type AuthResult struct {
	IdToken string `json:"IdToken"`
}

type AuthResponse struct {
	AuthenticationResult AuthResult `json:"AuthenticationResult"`
}

func NewAuthService(client *http.Client) *Auth {
	if client == nil {
		client = &http.Client{}
	}
	return &Auth{HTTPClient: client}
}

func (a *Auth) GetAccessToken() (*string, error) {
	url := os.Getenv("ACCESS_TOKEN_URL")
	username := os.Getenv("ACCESSTAGE_USERNAME")
	password := os.Getenv("ACCESSTAGE_PASSWORD")
	if url == "" || username == "" || password == "" {
		return nil, fmt.Errorf("required environment variables not found to access token")
	}

	body := []byte(`{
		"AuthParameters": {
			"USERNAME": "` + username + `",
			"PASSWORD": "` + password + `"
		}
	}`)

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := a.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if httpResp.StatusCode == http.StatusOK {
		var authResponse AuthResponse
		if err := json.Unmarshal(responseBody, &authResponse); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}

		return &authResponse.AuthenticationResult.IdToken, nil
	}

	log.Printf("token response error code: %d and body: %s", httpResp.StatusCode, string(responseBody))
	return nil, fmt.Errorf("error response token. Code: %d and ResponseBody: %s", httpResp.StatusCode, string(responseBody))
}
