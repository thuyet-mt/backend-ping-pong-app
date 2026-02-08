package service

import (
	"context"
	"errors"

	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/repository"
)

var ErrInvalidPlayerName = errors.New("full_name is required")

type PlayerService interface {
	GetAll(ctx context.Context) ([]models.Player, error)
	GetByID(ctx context.Context, id string) (*models.Player, error)
	Create(ctx context.Context, p *models.Player) (*models.Player, error)
}

type playerService struct {
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) PlayerService {
	return &playerService{repo: repo}
}

func (s *playerService) GetAll(ctx context.Context) ([]models.Player, error) {
	return s.repo.GetAll(ctx)
}

func (s *playerService) GetByID(ctx context.Context, id string) (*models.Player, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *playerService) Create(ctx context.Context, p *models.Player) (*models.Player, error) {
	if p.FullName == "" {
		return nil, ErrInvalidPlayerName
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
