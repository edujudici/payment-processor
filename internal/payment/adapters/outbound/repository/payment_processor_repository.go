package repository

import (
	"context"
	"database/sql"
	"errors"
	"payment-processor/internal/payment/domain"
	"payment-processor/internal/payment/ports"
)

type paymentProcessorRepositoryMySQL struct {
	db *sql.DB
}

// Garantia de implementação da interface
var _ ports.PaymentRepository = (*paymentProcessorRepositoryMySQL)(nil)

func NewPaymentProcessorRepositoryMySQL(db *sql.DB) ports.PaymentRepository {
	return &paymentProcessorRepositoryMySQL{
		db: db,
	}
}

func (r *paymentProcessorRepositoryMySQL) FindAll(ctx context.Context) ([]*domain.Payment, error) {
	query := `
		SELECT id, amount, status, created_at
		FROM payments
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment

	for rows.Next() {
		var payment domain.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Status,
			&payment.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		payments = append(payments, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *paymentProcessorRepositoryMySQL) FindByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `
		SELECT id, amount, status, created_at
		FROM payments
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var payment domain.Payment

	err := row.Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // padrão comum (não encontrado)
		}
		return nil, err
	}

	return &payment, nil
}

func (r *paymentProcessorRepositoryMySQL) Save(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, amount, status, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		payment.ID,
		payment.Amount,
		payment.Status,
		payment.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
