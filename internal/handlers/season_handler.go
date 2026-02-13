package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/service"
)

type SeasonHandler struct {
	service service.SeasonService
}

func NewSeasonHandler(svc service.SeasonService) *SeasonHandler {
	return &SeasonHandler{service: svc}
}

// GetSeasonsHandle handles GET /api/v1/seasons
func (h *SeasonHandler) GetSeasonsHandle(c *gin.Context) {
	seasons, err := h.service.GetAllSeasons(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seasons)
}

// GetSeasonByIDHandle handles GET /api/v1/seasons/{seasonId}
func (h *SeasonHandler) GetSeasonByIDHandle(c *gin.Context) {
	seasonID := c.Param("seasonId")
	if seasonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seasonId is required"})
		return
	}

	season, err := h.service.GetSeasonByID(c.Request.Context(), seasonID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, season)
}

type CreateSeasonRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Status    string `json:"status"`
}

// CreateSeasonHandle handles POST /api/v1/seasons
func (h *SeasonHandler) CreateSeasonHandle(c *gin.Context) {
	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Parse StartDate and EndDate strings to time.Time
	season := &models.Season{
		Name:   req.Name,
		Status: req.Status,
	}

	if season.Status == "" {
		season.Status = "UPCOMING"
	}

	created, err := h.service.CreateSeason(c.Request.Context(), season)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}
