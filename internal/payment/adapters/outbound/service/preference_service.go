package service

import (
	"context"
	"fmt"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PreferenceCreateInput struct {
	Items             []Item   `json:"items"`
	Payer             Payer    `json:"payer"`
	BackURLs          BackURLs `json:"back_urls"`
	AutoReturn        string   `json:"auto_return"`
	NotificationURL   string   `json:"notification_url"`
	ExternalReference string   `json:"external_reference"`
}

type Item struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	CurrencyID string  `json:"currency_id"`
}

type BackURLs struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Pending string `json:"pending"`
}

type Payer struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
}

type PreferenceCreateOutput struct {
	ID                string `json:"id"`
	ExternalReference string `json:"externalReference"`
	InitPoint         string `json:"initPoint"`
	SandboxInitPoint  string `json:"sandboxInitPoint"`
}

// type Preference struct {
// 	HTTPClient *http.Client
// }

// func NewPreferenceCreateService(client *http.Client) *Preference {
// 	if client == nil {
// 		client = &http.Client{}
// 	}
// 	return &Preference{HTTPClient: client}
// }

// func (p *Preference) PreferenceCreate(token string, req *PreferenceCreateInput) (*PreferenceCreateOutput, error) {
// 	url := os.Getenv("PREFERENCE_CREATE_URL")
// 	if url == "" {
// 		return nil, fmt.Errorf("PREFERENCE_CREATE_URL not set")
// 	}

// 	body, err := json.Marshal(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf(
// 		"preference create request | request=%s",
// 		body,
// 	)

// 	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating request: %v", err)
// 	}

// 	httpReq.Header.Set("Authorization", "Bearer "+token)
// 	httpReq.Header.Set("Content-Type", "application/json")

// 	httpResp, err := p.HTTPClient.Do(httpReq)
// 	if err != nil {
// 		return nil, fmt.Errorf("error executing request: %v", err)
// 	}
// 	defer httpResp.Body.Close()

// 	responseBody, err := io.ReadAll(httpResp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response body: %v", err)
// 	}

// 	log.Printf("preference create response code: %d and body: %s", httpResp.StatusCode, string(responseBody))

// 	if httpResp.StatusCode == http.StatusOK || httpResp.StatusCode == http.StatusBadRequest {
// 		var preferenceResponse PreferenceCreateOutput
// 		if err := json.Unmarshal(responseBody, &preferenceResponse); err != nil {
// 			return nil, fmt.Errorf("error parsing JSON: %v", err)
// 		}

// 		return &preferenceResponse, nil
// 	}

// 	return nil, fmt.Errorf("error preference response. Code: %d and ResponseBody: %s", httpResp.StatusCode, string(responseBody))
// }

type Preference struct {
	client preference.Client
}

func NewPreferenceCreateService(accessToken string) (*Preference, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return nil, err
	}

	return &Preference{
		client: preference.NewClient(cfg),
	}, nil
}

func (p *Preference) PreferenceCreate(ctx context.Context, req *PreferenceCreateInput) (*PreferenceCreateOutput, error) {

	items := make([]preference.ItemRequest, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, preference.ItemRequest{
			ID:         item.ID,
			Title:      item.Title,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			CurrencyID: item.CurrencyID,
		})
	}

	request := preference.Request{
		Items: items,
		BackURLs: &preference.BackURLsRequest{
			Success: req.BackURLs.Success,
			Failure: req.BackURLs.Failure,
			Pending: req.BackURLs.Pending,
		},
		AutoReturn:        req.AutoReturn,
		NotificationURL:   req.NotificationURL,
		ExternalReference: req.ExternalReference,
		Payer: &preference.PayerRequest{
			Name:    req.Payer.Name,
			Surname: req.Payer.Surname,
			Email:   req.Payer.Email,
		},
	}

	// 🚀 chamada SDK
	resource, err := p.client.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error creating preference: %w", err)
	}

	// 🔄 Mapeando SDK -> seu output
	return &PreferenceCreateOutput{
		ID:                resource.ID,
		ExternalReference: resource.ExternalReference,
		InitPoint:         resource.InitPoint,
		SandboxInitPoint:  resource.SandboxInitPoint,
	}, nil
}
