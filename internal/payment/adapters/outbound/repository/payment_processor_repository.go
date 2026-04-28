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

var _ ports.PaymentRepository = (*paymentProcessorRepositoryMySQL)(nil)

func NewPaymentProcessorRepositoryMySQL(db *sql.DB) ports.PaymentRepository {
	return &paymentProcessorRepositoryMySQL{
		db: db,
	}
}

func (r *paymentProcessorRepositoryMySQL) FindAll(ctx context.Context) ([]*domain.Payment, error) {
	query := `
		SELECT id, protocol, username, surname, status, quantity, total, subtotal, description, preference_id, external_reference, preference_init_point, preference_sandbox_init_point, created_at, updated_at
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
			&payment.Protocol,
			&payment.Username,
			&payment.Surname,
			&payment.Status,
			&payment.Quantity,
			&payment.Total,
			&payment.Subtotal,
			&payment.Description,
			&payment.PreferenceID,
			&payment.ExternalReference,
			&payment.InitPoint,
			&payment.SandboxInitPoint,
			&payment.CreatedAt,
			&payment.UpdatedAt,
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
		SELECT id, protocol, username, surname, status, quantity, total, subtotal, description, preference_id, external_reference, preference_init_point, preference_sandbox_init_point, created_at, updated_at
		FROM payments
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var payment domain.Payment

	err := row.Scan(
		&payment.ID,
		&payment.Protocol,
		&payment.Username,
		&payment.Surname,
		&payment.Status,
		&payment.Quantity,
		&payment.Total,
		&payment.Subtotal,
		&payment.Description,
		&payment.PreferenceID,
		&payment.ExternalReference,
		&payment.InitPoint,
		&payment.SandboxInitPoint,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

func (r *paymentProcessorRepositoryMySQL) Save(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, protocol, username, surname, email, status, quantity, total, subtotal, description, preference_id, external_reference, preference_init_point, preference_sandbox_init_point, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		payment.ID,
		payment.Protocol,
		payment.Username,
		payment.Surname,
		payment.Email,
		payment.Status,
		payment.Quantity,
		payment.Total,
		payment.Subtotal,
		payment.Description,
		payment.PreferenceID,
		payment.ExternalReference,
		payment.InitPoint,
		payment.SandboxInitPoint,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
