package service

import "backend-ping-pong-app/internal/repository"

// Service là struct gốc, chứa toàn bộ service của app
type Service struct {
	Player PlayerService
	Season SeasonService
	Team   TeamService
}

// NewService khởi tạo toàn bộ service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Player: NewPlayerService(repo.Player),
		Season: NewSeasonService(repo.Season),
		Team:   NewTeamService(repo.Team),
	}
}
