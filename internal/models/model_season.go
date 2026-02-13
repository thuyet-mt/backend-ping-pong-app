package models

import "time"

// Season represents a league season
type Season struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"` // UPCOMING, ACTIVE, FINISHED
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PlayerRating represents player's rating in a season
type PlayerRating struct {
	PlayerID          string    `json:"player_id"`
	SeasonID          string    `json:"season_id"`
	Points            int       `json:"points"`
	Rank              string    `json:"rank"`
	AccumulatedPoints int       `json:"accumulated_points"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// SeasonListResponse for GET /seasons
type SeasonListResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
}

// TeamListResponse for GET /seasons/{seasonId}/teams
type TeamListResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}
