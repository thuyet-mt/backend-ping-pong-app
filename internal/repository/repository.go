package repository

import "database/sql"

type Repository struct {
	Player PlayerRepository
	Season SeasonRepository
	Team   TeamRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Player: NewPlayerRepository(db),
		Season: NewSeasonRepository(db),
		Team:   NewTeamRepository(db),
	}
}
