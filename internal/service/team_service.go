package service

import (
	"context"
	"errors"

	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/repository"
	"backend-ping-pong-app/internal/utils"
)

var ErrTeamNotFound = errors.New("team not found")

type TeamService interface {
	GetTeamsBySeasonIDService(ctx context.Context, seasonID string) ([]models.TeamListResponse, error)
	GetTeamByIDService(ctx context.Context, id string) (*models.Team, error)
	CreateTeamService(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeamMembersService(ctx context.Context, teamID string) ([]models.TeamMember, error)
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{repo: repo}
}

func (s *teamService) GetTeamsBySeasonIDService(ctx context.Context, seasonID string) ([]models.TeamListResponse, error) {
	teams, err := s.repo.GetTeamsBySeasonIDRepo(ctx, seasonID)
	if err != nil {
		return nil, err
	}

	response := make([]models.TeamListResponse, 0, len(teams))
	for _, team := range teams {
		var avatarPath *string
		if team.AvatarURL != nil {
			url := utils.BuildCDNURL(*team.AvatarURL)
			avatarPath = &url
		}
		response = append(response, models.TeamListResponse{
			ID:        team.ID,
			Name:      team.Name,
			AvatarURL: avatarPath,
		})
	}
	return response, nil
}

func (s *teamService) GetTeamByIDService(ctx context.Context, id string) (*models.Team, error) {
	team, err := s.repo.GetTeamByIDRepo(ctx, id)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

func (s *teamService) CreateTeamService(ctx context.Context, team *models.Team) (*models.Team, error) {
	return s.repo.CreateTeamRepo(ctx, team)
}

func (s *teamService) GetTeamMembersService(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	return s.repo.GetTeamMembersRepo(ctx, teamID)
}
