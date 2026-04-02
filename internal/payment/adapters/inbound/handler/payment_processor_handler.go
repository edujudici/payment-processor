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
	result, err := h.payGetUseCase.Execute(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *PaymentProcessorHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req dto.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err := h.payCreateUseCase.Execute(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Payment created successfully"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
