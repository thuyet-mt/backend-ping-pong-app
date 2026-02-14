package repository

import (
	"context"
	"database/sql"

	"backend-ping-pong-app/internal/models"
)

type TeamRepository interface {
	GetTeamsBySeasonIDRepo(ctx context.Context, seasonID string) ([]models.Team, error)
	GetTeamByIDRepo(ctx context.Context, id string) (*models.Team, error)
	CreateTeamRepo(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeamMembersRepo(ctx context.Context, teamID string) ([]models.TeamMember, error)
}

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) GetTeamsBySeasonIDRepo(ctx context.Context, seasonID string) ([]models.Team, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, season_id, name, avatar_url
		FROM teams
		WHERE season_id = $1
		ORDER BY name ASC
	`, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.ID,
			&team.SeasonID,
			&team.Name,
			&team.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *teamRepository) GetTeamByIDRepo(ctx context.Context, id string) (*models.Team, error) {
	var team models.Team
	err := r.db.QueryRowContext(ctx, `
		SELECT id, season_id, name, avatar_url
		FROM teams
		WHERE id = $1
	`, id).Scan(
		&team.ID,
		&team.SeasonID,
		&team.Name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *teamRepository) CreateTeamRepo(ctx context.Context, team *models.Team) (*models.Team, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO teams (season_id, name)
		VALUES ($1, $2)
		RETURNING id
	`, team.SeasonID, team.Name).Scan(&id)

	if err != nil {
		return nil, err
	}

	team.ID = id
	return team, nil
}

func (r *teamRepository) GetTeamMembersRepo(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT team_id, player_id, joined_round, left_round, transfer_type
		FROM team_members
		WHERE team_id = $1
		ORDER BY joined_round ASC
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		err := rows.Scan(
			&member.TeamID,
			&member.PlayerID,
			&member.JoinedRound,
			&member.LeftRound,
			&member.TransferType,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
