package service

import "backend-ping-pong-app/internal/repository"

// Service là struct gốc, chứa toàn bộ service của app
type Service struct {
	Player PlayerService
}

// NewService khởi tạo toàn bộ service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Player: NewPlayerService(repo.Player),
	}
}
