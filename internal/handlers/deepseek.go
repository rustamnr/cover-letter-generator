package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/services"
)

// DeepSeekHandler содержит сервис DeepSeek
type DeepSeekHandler struct {
	deepSeekService *services.DeepSeekService
}

// NewDeepSeekHandler создает новый обработчик
func NewDeepSeekHandler(deepSeekService *services.DeepSeekService) *DeepSeekHandler {
	return &DeepSeekHandler{deepSeekService: deepSeekService}
}

// HandleDeepSeek обрабатывает запросы к DeepSeek API
func (h *DeepSeekHandler) HandleDeepSeek(c *gin.Context) {
	var request struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	response, err := h.deepSeekService.SendRequest(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DeepSeekHandler) HandleDeepSeekCoverLetterGenerate(c *gin.Context) {
	var request struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	response, err := h.deepSeekService.SendRequest(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
