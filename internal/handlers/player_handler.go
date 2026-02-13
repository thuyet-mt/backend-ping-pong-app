package handlers

import (
	"backend-ping-pong-app/internal/models"
	"backend-ping-pong-app/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreatePlayerRequest struct {
	FullName  string  `json:"full_name" binding:"required"`
	BirthYear *int    `json:"birth_year"binding:"required"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"binding:"required"`
}

type PlayerHandler struct {
	service service.PlayerService
}

func NewPlayerHandler(svc service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: svc}
}

func (h *PlayerHandler) GetPlayersHandle(c *gin.Context) {
	players, err := h.service.GetAllPlayerService(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) SearchPlayersHandle(c *gin.Context) {
	search := c.Query("search")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search parameter is required"})
		return
	}

	players, err := h.service.SearchPlayerByNameService(c.Request.Context(), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) CreatePlayerHandle(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player := &models.Player{
		FullName:  req.FullName,
		BirthYear: req.BirthYear,
		AvatarURL: req.AvatarURL,
	}

	created, err := h.service.CreatePlayerService(c.Request.Context(), player)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}
