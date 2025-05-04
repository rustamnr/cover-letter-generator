package handlers

import (
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
