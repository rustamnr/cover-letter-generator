package handlers_chatgpt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/services"
)

// ChatGPTHandler содержит сервис ChatGPT
type ChatGPTHandler struct {
	chatGPTService *services.ChatGPTService
}

// NewChatGPTHandler создает новый обработчик
func NewChatGPTHandler(chatGPTService *services.ChatGPTService) *ChatGPTHandler {
	return &ChatGPTHandler{chatGPTService: chatGPTService}
}

// HandleChatGPT обрабатывает запросы к ChatGPT API
func (h *ChatGPTHandler) HandleChatGPT(c *gin.Context) {
	var request struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	response, err := h.chatGPTService.SendRequest(request.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
