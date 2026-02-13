package repository

import (
	"context"
	"database/sql"
	"time"

	"backend-ping-pong-app/internal/models"
)

type SeasonRepository interface {
	GetAllSeasons(ctx context.Context) ([]models.Season, error)
	GetSeasonByID(ctx context.Context, id string) (*models.Season, error)
	CreateSeason(ctx context.Context, season *models.Season) (*models.Season, error)
	UpdateSeason(ctx context.Context, season *models.Season) (*models.Season, error)
}

type seasonRepository struct {
	db *sql.DB
}

func NewSeasonRepository(db *sql.DB) SeasonRepository {
	return &seasonRepository{db: db}
}

func (r *seasonRepository) GetAllSeasons(ctx context.Context) ([]models.Season, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, start_date, end_date, status, created_at, update_at
		FROM seasons
		ORDER BY start_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seasons []models.Season
	for rows.Next() {
		var season models.Season
		err := rows.Scan(
			&season.ID,
			&season.Name,
			&season.StartDate,
			&season.EndDate,
			&season.Status,
			&season.CreatedAt,
			&season.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		seasons = append(seasons, season)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seasons, nil
}

func (r *seasonRepository) GetSeasonByID(ctx context.Context, id string) (*models.Season, error) {
	var season models.Season
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, start_date, end_date, status, created_at, update_at
		FROM seasons
		WHERE id = $1
	`, id).Scan(
		&season.ID,
		&season.Name,
		&season.StartDate,
		&season.EndDate,
		&season.Status,
		&season.CreatedAt,
		&season.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &season, nil
}

func (r *seasonRepository) CreateSeason(ctx context.Context, season *models.Season) (*models.Season, error) {
	now := time.Now()
	season.CreatedAt = now
	season.UpdatedAt = now

	if season.Status == "" {
		season.Status = "UPCOMING"
	}

	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO seasons (name, start_date, end_date, status, created_at, update_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, season.Name, season.StartDate, season.EndDate, season.Status, season.CreatedAt, season.UpdatedAt).Scan(&id)

	if err != nil {
		return nil, err
	}

	season.ID = id
	return season, nil
}

func (r *seasonRepository) UpdateSeason(ctx context.Context, season *models.Season) (*models.Season, error) {
	season.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, `
		UPDATE seasons
		SET name = $1, start_date = $2, end_date = $3, status = $4, update_at = $5
		WHERE id = $6
		RETURNING id, created_at
	`, season.Name, season.StartDate, season.EndDate, season.Status, season.UpdatedAt, season.ID).Scan(&season.ID, &season.CreatedAt)

	if err != nil {
		return nil, err
	}

	return season, nil
}
