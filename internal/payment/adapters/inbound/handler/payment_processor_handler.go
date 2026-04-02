package handler

import (
	"encoding/json"
	"net/http"

	"payment-processor/internal/payment/adapters/inbound/dto"
	"payment-processor/internal/payment/usecase"
)

type PaymentProcessorHandler struct {
	payCreateUseCase usecase.PaymentCreateUseCaseInterface
	payGetUseCase    usecase.PaymentGetUseCaseInterface
}

func NewPaymentProcessorHandler(payCreateUseCase usecase.PaymentCreateUseCaseInterface, payGetUseCase usecase.PaymentGetUseCaseInterface) *PaymentProcessorHandler {
	return &PaymentProcessorHandler{
		payCreateUseCase: payCreateUseCase,
		payGetUseCase:    payGetUseCase,
	}
}

func (h *PaymentProcessorHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	output, err := h.payGetUseCase.Execute(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, output)
}

func (h *PaymentProcessorHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePaymentInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.payCreateUseCase.Execute(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, output)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
