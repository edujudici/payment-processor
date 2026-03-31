package ports

import "context"

type SendMessageInterface interface {
	Send(ctx context.Context, messageBody string) error
}
