package main

import (
	"github.com/joho/godotenv"
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/server"
)


func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("failed to load .env file: %v", err)
	}

	// Запуск сервера
	srv := server.NewServer()
	srv.Run("8080")
}
