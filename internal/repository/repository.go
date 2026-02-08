package repository

import "database/sql"

type Repository struct {
	Player PlayerRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Player: &playerRepository{db: db},
	}
}
