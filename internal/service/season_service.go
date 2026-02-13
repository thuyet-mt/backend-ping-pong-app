package service

import (
	"context"
	"errors"

	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/repository"
)

var ErrSeasonNotFound = errors.New("season not found")

type SeasonService interface {
	GetAllSeasons(ctx context.Context) ([]models.SeasonListResponse, error)
	GetSeasonByID(ctx context.Context, id string) (*models.Season, error)
	CreateSeason(ctx context.Context, season *models.Season) (*models.Season, error)
	UpdateSeason(ctx context.Context, season *models.Season) (*models.Season, error)
}

type seasonService struct {
	repo repository.SeasonRepository
}

func NewSeasonService(repo repository.SeasonRepository) SeasonService {
	return &seasonService{repo: repo}
}

func (s *seasonService) GetAllSeasons(ctx context.Context) ([]models.SeasonListResponse, error) {
	seasons, err := s.repo.GetAllSeasons(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]models.SeasonListResponse, 0, len(seasons))
	for _, season := range seasons {
		response = append(response, models.SeasonListResponse{
			ID:        season.ID,
			Name:      season.Name,
			StartDate: season.StartDate,
			EndDate:   season.EndDate,
			Status:    season.Status,
		})
	}
	return response, nil
}

func (s *seasonService) GetSeasonByID(ctx context.Context, id string) (*models.Season, error) {
	season, err := s.repo.GetSeasonByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if season == nil {
		return nil, ErrSeasonNotFound
	}

	return season, nil
}

func (s *seasonService) CreateSeason(ctx context.Context, season *models.Season) (*models.Season, error) {
	return s.repo.CreateSeason(ctx, season)
}

func (s *seasonService) UpdateSeason(ctx context.Context, season *models.Season) (*models.Season, error) {
	return s.repo.UpdateSeason(ctx, season)
}
