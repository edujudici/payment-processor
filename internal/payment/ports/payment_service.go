package ports

import "payment-processor/internal/payment/adapters/outbound/service"

type PaymentSvcInterface interface {
	PaymentCancel(token string, req *service.BoletoCancelInput) (int, *service.BoletoCancelOutput, error)
}
