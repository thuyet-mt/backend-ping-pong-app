package repository

import (
	"context"
	"database/sql"

	"backend-ping-pong-app/internal/models"
)

type playerRepository struct {
	db *sql.DB
}

type PlayerRepository interface {
	GetAll(ctx context.Context) ([]models.Player, error)
	GetByID(ctx context.Context, id string) (*models.Player, error)
	Create(ctx context.Context, p *models.Player) error
}

func (r *playerRepository) GetAll(ctx context.Context) ([]models.Player, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			id,
			full_name,
			birth_year,
			phone,
			avatar_url,
			created_at
		FROM players
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.Player

	for rows.Next() {
		var p models.Player
		if err := rows.Scan(
			&p.ID,
			&p.FullName,
			&p.BirthYear,
			&p.Phone,
			&p.AvatarURL,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func (r *playerRepository) GetByID(ctx context.Context, id string) (*models.Player, error) {
	var p models.Player

	err := r.db.QueryRowContext(ctx, `
		SELECT 
			id,
			full_name,
			birth_year,
			phone,
			cccd,
			avatar_url,
			is_active,
			created_at
		FROM players
		WHERE id = $1
	`, id).Scan(
		&p.ID,
		&p.FullName,
		&p.BirthYear,
		&p.Phone,
		&p.CCCD,
		&p.AvatarURL,
		&p.IsActive,
		&p.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *playerRepository) Create(ctx context.Context, p *models.Player) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO players (
			full_name,
			birth_year,
			phone,
			cccd,
			avatar_url
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, is_active, created_at
	`,
		p.FullName,
		p.BirthYear,
		p.Phone,
		p.CCCD,
		p.AvatarURL,
	).Scan(
		&p.ID,
		&p.IsActive,
		&p.CreatedAt,
	)
}
