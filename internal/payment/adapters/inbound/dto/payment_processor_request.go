package dto

import (
	"fmt"
)

type Request struct {
	Amount float64 `json:"amount"`
}

func (r Request) Validate() error {
	if r.Amount == 0 {
		return fmt.Errorf("amount is required")
	}

	return nil
}
