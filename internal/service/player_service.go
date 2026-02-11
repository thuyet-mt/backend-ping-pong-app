package service

import (
	"context"
	"errors"

	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/repository"
	"backend-ping-pong-app/internal/utils"
)

var ErrInvalidPlayerName = errors.New("full_name is required")

type PlayerService interface {
	GetAllPlayerService(ctx context.Context) ([]models.PlayerListResponse, error)
	SearchPlayerByNameService(ctx context.Context, name string) ([]models.PlayerListResponse, error)
	CreatePlayerService(ctx context.Context, p *models.Player) (*models.Player, error)
}

type playerService struct {
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) PlayerService {
	return &playerService{repo: repo}
}

func (s *playerService) GetAllPlayerService(ctx context.Context) ([]models.PlayerListResponse, error) {
	items, err := s.repo.GetAllPlayerRepo(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.PlayerListResponse, 0, len(items))
	for _, p := range items {
		var avatarPath *string
		if p.AvatarPath != nil {
			url := utils.BuildCDNURL(*p.AvatarPath)
			avatarPath = &url
		}
		res = append(res, models.PlayerListResponse{
			FullName:   p.FullName,
			BirthYear:  p.BirthYear,
			AvatarPath: avatarPath,
		})
	}
	return res, nil
}

func (s *playerService) SearchPlayerByNameService(ctx context.Context, name string) ([]models.PlayerListResponse, error) {
	if name == "" {
		return nil, errors.New("Name is required")
	}
	items, err := s.repo.SearchByNameRepo(ctx, name)
	if err != nil {
		return nil, err
	}

	res := make([]models.PlayerListResponse, 0, len(items))
	for _, p := range items {
		var avatarPath *string
		if p.AvatarPath != nil {
			url := utils.BuildCDNURL(*p.AvatarPath)
			avatarPath = &url
		}
		res = append(res, models.PlayerListResponse{
			FullName:   p.FullName,
			BirthYear:  p.BirthYear,
			AvatarPath: avatarPath,
		})
	}
	return res, nil
}

func (s *playerService) CreatePlayerService(ctx context.Context, p *models.Player) (*models.Player, error) {
	if p.FullName == "" {
		return nil, ErrInvalidPlayerName
	}

	if err := s.repo.CreatePlayerRepo(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}
