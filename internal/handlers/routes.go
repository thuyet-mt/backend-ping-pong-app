package handlers

import (
	"github.com/gin-gonic/gin"

	"backend-ping-pong-app/internal/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.Service) {
	playerHandler := NewPlayerHandler(svc.Player)
	seasonHandler := NewSeasonHandler(svc.Season)
	teamHandler := NewTeamHandler(svc.Team)

	v1 := r.Group("/api/v1")
	{
		// Player routes
		v1.GET("/players", playerHandler.GetPlayersHandle)
		v1.GET("/players/search", playerHandler.SearchPlayersHandle)
		v1.POST("/players", playerHandler.CreatePlayerHandle)

		// Season routes
		v1.GET("/seasons", seasonHandler.GetSeasonsHandle)
		v1.GET("/seasons/:seasonId", seasonHandler.GetSeasonByIDHandle)
		v1.POST("/seasons", seasonHandler.CreateSeasonHandle)

		// Team routes
		v1.GET("/seasons/:seasonId/teams", teamHandler.GetTeamsBySeasonHandle)
		v1.GET("/teams/:teamId", teamHandler.GetTeamByIDHandle)
		v1.POST("/seasons/:seasonId/teams", teamHandler.CreateTeamHandle)
		v1.GET("/teams/:teamId/members", teamHandler.GetTeamMembersHandle)
	}
}
