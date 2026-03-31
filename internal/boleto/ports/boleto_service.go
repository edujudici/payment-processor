package ports

import "boleto-cancel/internal/boleto/adapters/outbound/service"

type BoletoCancelInterface interface {
	BoletoCancel(token string, req *service.BoletoCancelInput) (int, *service.BoletoCancelOutput, error)
}
