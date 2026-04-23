package usecase

import (
	"context"
	"fmt"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/adapters/outbound/service"
	"payment-processor/internal/payment/ports"
)

type PreferenceCreateUseCaseInterface interface {
	Execute(ctx context.Context, input dto.CreatePreferenceInput) (*dto.CreatePreferenceOutput, error)
}

type PreferenceCreateUseCase struct {
	preferenceService    ports.PreferenceSvcInterface
	preferenceRepository ports.PreferenceRepository
}

func NewPreferenceCreateUseCase(
	preferenceService ports.PreferenceSvcInterface,
	preferenceRepository ports.PreferenceRepository,
) *PreferenceCreateUseCase {
	return &PreferenceCreateUseCase{
		preferenceService:    preferenceService,
		preferenceRepository: preferenceRepository,
	}
}

func (uc *PreferenceCreateUseCase) Execute(ctx context.Context, input dto.CreatePreferenceInput) (*dto.CreatePreferenceOutput, error) {

	pref, err := uc.preferenceService.PreferenceCreate(ctx, &service.PreferenceCreateInput{
		Items: dto.MapItems(input.Items),
		Payer: dto.MapPayer(input.Payer),
		BackURLs: service.BackURLs{
			Success: "https://example.com/success",
			Failure: "https://example.com/failure",
			Pending: "https://example.com/pending",
		},
		AutoReturn:        "approved",
		NotificationURL:   "https://example.com/notification",
		ExternalReference: input.ExternalReference,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	preference, err := dto.ToPreference(pref.ID, pref.ExternalReference, pref.InitPoint, pref.SandboxInitPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	err = uc.preferenceRepository.Save(ctx, preference)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	return dto.FromPreference(preference), nil
}
