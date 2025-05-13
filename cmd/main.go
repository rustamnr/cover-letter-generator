package main

import (
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/server"
	"github.com/rustamnr/cover-letter-generator/internal/storage"
	"github.com/rustamnr/cover-letter-generator/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("failed to load .env file: %v", err)
	}

	go func() {
		storage.InitRedis("localhost:6379")
	}()

	// Запуск сервера
	go func() {
		srv := server.NewServer()
		srv.Run("8080")
	}()

	// Запуск Telegram бота
	bot, err := telegram.NewBot()
	if err != nil {
		logger.Fatalf("failed to create bot: %v", err)
	}
	bot.Start()
}
