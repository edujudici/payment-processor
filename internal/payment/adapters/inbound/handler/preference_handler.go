package handler

import (
	"encoding/json"
	"net/http"
	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/usecase"
)

type PreferenceHandler struct {
	preferenceCreateUseCase usecase.PreferenceCreateUseCaseInterface
}

func NewPreferenceHandler(preferenceCreateUseCase usecase.PreferenceCreateUseCaseInterface) *PreferenceHandler {
	return &PreferenceHandler{
		preferenceCreateUseCase: preferenceCreateUseCase,
	}
}

func (h *PreferenceHandler) CreatePreference(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePreferenceInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.preferenceCreateUseCase.Execute(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, output)
}
