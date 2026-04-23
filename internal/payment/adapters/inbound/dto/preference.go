package dto

import (
	"payment-processor/internal/payment/adapters/outbound/service"
	"payment-processor/internal/payment/domain"
	"time"
)

type CreatePreferenceInput struct {
	Items             []Item `json:"items"`
	Payer             Payer  `json:"payer"`
	ExternalReference string `json:"external_reference"`
}

type Item struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	CurrencyID string  `json:"currency_id"`
}

type Payer struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
}

type CreatePreferenceOutput struct {
	ID                string    `json:"id"`
	InitPoint         string    `json:"init_point"`
	SandboxInitPoint  string    `json:"sandbox_init_point"`
	ExternalReference string    `json:"external_reference"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ToPreference(pref, extRef, initPoint, sandboxInitPoint string) (*domain.Preference, error) {
	return domain.NewPreference(
		pref,
		extRef,
		initPoint,
		sandboxInitPoint,
	)
}

func FromPreference(preference *domain.Preference) *CreatePreferenceOutput {
	return &CreatePreferenceOutput{
		ID:                preference.ID,
		InitPoint:         preference.InitPoint,
		SandboxInitPoint:  preference.SandboxInitPoint,
		ExternalReference: preference.ExternalReference,
		CreatedAt:         preference.CreatedAt,
		UpdatedAt:         preference.UpdatedAt,
	}
}

func MapItems(dtoItems []Item) []service.Item {
	var items []service.Item

	for _, item := range dtoItems {
		items = append(items, service.Item{
			ID:         item.ID,
			Title:      item.Title,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			CurrencyID: item.CurrencyID,
		})
	}

	return items
}

func MapPayer(dtoPayer Payer) service.Payer {
	return service.Payer{
		Name:        dtoPayer.Name,
		Surname:     dtoPayer.Surname,
		Email:       dtoPayer.Email,
		DateCreated: dtoPayer.DateCreated,
	}
}
