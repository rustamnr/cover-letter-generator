package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Server структура для хранения объекта сервера
type Server struct {
	Router *gin.Engine
}

// NewServer создает новый сервер и настраивает маршруты
func NewServer() *Server {
	router := gin.Default()

	// Регистрируем маршруты
	registerRoutes(router)

	return &Server{Router: router}
}

// Run запускает сервер
func (s *Server) Run(port string) {
	log.Printf("Сервер запущен на :%s", port)
	if err := s.Router.Run(":" + port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
