package handlers

import (
	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreatePlayerRequest struct {
	FullName  string  `json:"full_name" binding:"required"`
	BirthYear *int    `json:"birth_year"`
	Phone     *string `json:"phone"`
	CCCD      *string `json:"cccd"`
}

type PlayerHandler struct {
	service service.PlayerService
}

func NewPlayerHandler(svc service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: svc}
}

func (h *PlayerHandler) GetPlayers(c *gin.Context) {
	players, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player := &models.Player{
		FullName:  req.FullName,
		BirthYear: req.BirthYear,
		Phone:     req.Phone,
		CCCD:      req.CCCD,
	}

	created, err := h.service.Create(c.Request.Context(), player)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}
