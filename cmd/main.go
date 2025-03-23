package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rustamnr/cover-letter-generator/internal/server"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("Не найден .env файл")
	}

	// Запуск сервера
	srv := server.NewServer()
	srv.Run("8080")
}
