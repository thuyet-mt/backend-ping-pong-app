package repository

import (
	"context"
	"database/sql"

	"backend-ping-pong-app/internal/models"
)

type playerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) PlayerRepository {
	return &playerRepository{db: db}
}

type PlayerRepository interface {
	GetAllPlayerRepo(ctx context.Context) ([]models.PlayerListResponse, error)
	SearchByNameRepo(ctx context.Context, name string) ([]models.PlayerListResponse, error)
	CreatePlayerRepo(ctx context.Context, p *models.Player) error
}

func (r *playerRepository) GetAllPlayerRepo(ctx context.Context) ([]models.PlayerListResponse, error) {
	queryGetAllPlayer := `
		SELECT
			full_name,
			date_of_birth,
			avatar_url
		FROM players
		ORDER BY created_at DESC;`
	rows, err := r.db.QueryContext(ctx, queryGetAllPlayer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var PlayerListResponses []models.PlayerListResponse

	for rows.Next() {
		var p models.PlayerListResponse
		if err := rows.Scan(
			&p.FullName,
			&p.DateOfBirth,
			&p.AvatarPath,
		); err != nil {
			return nil, err
		}
		PlayerListResponses = append(PlayerListResponses, p)
	}

	return PlayerListResponses, nil
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

func (r *playerRepository) SearchByNameRepo(ctx context.Context, name string) ([]models.PlayerListResponse, error) {
	var PlayerListResponses []models.PlayerListResponse

	querySearchByName := `
		SELECT 
			full_name,
			date_of_birth,
			avatar_url
		FROM players 
		WHERE full_name ILIKE '%' || $1 || '%'
	`

	rows, err := r.db.QueryContext(ctx, querySearchByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PlayerListResponse
		if err := rows.Scan(
			&p.FullName,
			&p.DateOfBirth,
			&p.AvatarPath,
		); err != nil {
			return nil, err
		}
		PlayerListResponses = append(PlayerListResponses, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return PlayerListResponses, nil
}

func (r *playerRepository) CreatePlayerRepo(ctx context.Context, p *models.Player) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO players (
			full_name,
			date_of_birth,
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
