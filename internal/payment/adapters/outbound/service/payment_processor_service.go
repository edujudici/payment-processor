package service

import (
	"context"
	"fmt"
	"log"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PaymentCreateInput struct {
	Item              Item     `json:"items"`
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

type PaymentCreateOutput struct {
	ID                string `json:"id"`
	ExternalReference string `json:"externalReference"`
	InitPoint         string `json:"initPoint"`
	SandboxInitPoint  string `json:"sandboxInitPoint"`
}

type Payment struct {
	client preference.Client
}

func NewPaymentCreateService(accessToken string) (*Payment, error) {
	cfg, err := config.New(accessToken)
	if err != nil {
		return nil, err
	}

	return &Payment{
		client: preference.NewClient(cfg),
	}, nil
}

func (p *Payment) PaymentCreate(ctx context.Context, req *PaymentCreateInput) (*PaymentCreateOutput, error) {

	items := make([]preference.ItemRequest, 0, 1)
	items = append(items, preference.ItemRequest{
		ID:         req.Item.ID,
		Title:      req.Item.Title,
		Quantity:   req.Item.Quantity,
		UnitPrice:  req.Item.UnitPrice,
		CurrencyID: req.Item.CurrencyID,
	})

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

	log.Printf("Creating payment with request: %+v", request)

	resource, err := p.client.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	log.Printf("Payment created successfully: %+v", resource)

	return &PaymentCreateOutput{
		ID:                resource.ID,
		ExternalReference: resource.ExternalReference,
		InitPoint:         resource.InitPoint,
		SandboxInitPoint:  resource.SandboxInitPoint,
	}, nil
}
