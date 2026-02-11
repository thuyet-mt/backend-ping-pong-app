package handlers

import (
	"github.com/gin-gonic/gin"

	"backend-ping-pong-app/internal/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.Service) {
	playerHandler := NewPlayerHandler(svc.Player)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/players", playerHandler.GetPlayersHandle)
		v1.GET("/players/search", playerHandler.SearchPlayersHandle)
		v1.POST("/players", playerHandler.CreatePlayerHandle)
	}
}
