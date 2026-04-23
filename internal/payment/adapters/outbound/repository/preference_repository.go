package repository

import (
	"context"
	"database/sql"
	"errors"
	"payment-processor/internal/payment/domain"
	"payment-processor/internal/payment/ports"
)

type preferenceRepositoryMySQL struct {
	db *sql.DB
}

// Garantia de implementação da interface
var _ ports.PreferenceRepository = (*preferenceRepositoryMySQL)(nil)

func NewPreferenceRepositoryMySQL(db *sql.DB) ports.PreferenceRepository {
	return &preferenceRepositoryMySQL{
		db: db,
	}
}

func (r *preferenceRepositoryMySQL) FindAll(ctx context.Context) ([]*domain.Preference, error) {
	query := `
		SELECT id, amount, status, created_at, updated_at, preference_id, description, preference_type
		FROM preferences
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []*domain.Preference

	for rows.Next() {
		var preference domain.Preference

		err := rows.Scan(
			&preference.ID,
			&preference.CreatedAt,
			&preference.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		preferences = append(preferences, &preference)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return preferences, nil
}

func (r *preferenceRepositoryMySQL) FindByID(ctx context.Context, id string) (*domain.Preference, error) {
	query := `
		SELECT id, amount, status, created_at, updated_at, preference_id, description, preference_type
		FROM preferences
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var preference domain.Preference

	err := row.Scan(
		&preference.ID,
		&preference.CreatedAt,
		&preference.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // padrão comum (não encontrado)
		}
		return nil, err
	}

	return &preference, nil
}

func (r *preferenceRepositoryMySQL) Save(ctx context.Context, preference *domain.Preference) error {
	query := `
		INSERT INTO preferences (id, amount, status, created_at, updated_at, preference_id, description, preference_type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		preference.ID,
		preference.CreatedAt,
		preference.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
