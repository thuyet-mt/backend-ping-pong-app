package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/service"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{service: svc}
}

// GetTeamsBySeasonHandle handles GET /api/v1/seasons/{seasonId}/teams
func (h *TeamHandler) GetTeamsBySeasonHandle(c *gin.Context) {
	seasonID := c.Param("seasonId")
	if seasonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seasonId is required"})
		return
	}

	teams, err := h.service.GetTeamsBySeasonIDService(c.Request.Context(), seasonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeamByIDHandle handles GET /api/v1/teams/{teamId}
func (h *TeamHandler) GetTeamByIDHandle(c *gin.Context) {
	teamID := c.Param("teamId")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is required"})
		return
	}

	team, err := h.service.GetTeamByIDService(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

type CreateTeamRequest struct {
	SeasonID string `json:"season_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// CreateTeamHandle handles POST /api/v1/seasons/{seasonId}/teams
func (h *TeamHandler) CreateTeamHandle(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team := &models.Team{
		SeasonID: req.SeasonID,
		Name:     req.Name,
	}

	created, err := h.service.CreateTeamService(c.Request.Context(), team)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetTeamMembersHandle handles GET /api/v1/teams/{teamId}/members
func (h *TeamHandler) GetTeamMembersHandle(c *gin.Context) {
	teamID := c.Param("teamId")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is required"})
		return
	}

	members, err := h.service.GetTeamMembersService(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}
