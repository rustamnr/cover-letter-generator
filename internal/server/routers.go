package server

import (
	"net/http"

	clients_deepseek "github.com/rustamnr/cover-letter-generator/internal/clients/deepseek"
	clients_hh "github.com/rustamnr/cover-letter-generator/internal/clients/hh"
	handlers_app "github.com/rustamnr/cover-letter-generator/internal/handlers/app"
	handlers_hh "github.com/rustamnr/cover-letter-generator/internal/handlers/hh"
	"github.com/rustamnr/cover-letter-generator/internal/middleware"
	"github.com/rustamnr/cover-letter-generator/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// registerRoutes настраивает маршруты API
func registerRoutes(router *gin.Engine) {
	// Инициализация клиентов
	hhClient := clients_hh.NewHHClient()
	deepSeekClient := clients_deepseek.NewDeepSeekClient()

	vacancyProvider := services.NewHHProvider(hhClient)
	textGenerator := services.NewDeepSeekService(deepSeekClient)
	// Инициализация сервисов
	vacancyQueue := services.NewSliceVacancyQueue()
	applicationService := services.NewApplicationService(vacancyProvider, vacancyQueue, textGenerator)

	// Инициализация хендлеров
	hhHandler := handlers_hh.NewHHHandler(hhClient)
	applicationHandler := handlers_app.NewApplicationHandler(applicationService)

	// Настройка сессий
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",       // Доступность для всех путей
		Domain:   "",        // Пусто = текущий домен
		MaxAge:   86400 * 7, // Время жизни в секундах (7 дней)
		Secure:   false,     // true для HTTPS только
		HttpOnly: true,      // Запрет доступа из JavaScript
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("session", store))

	// HH.ru API
	router.GET("/auth", hhHandler.AuthHandler)
	router.GET("/auth/callback", hhHandler.CallbackHandler)

	// API группы
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/resumes", hhHandler.GetUserResumes)
		api.POST("/resumes/current", hhHandler.SetCurrnetResume)
		api.GET("/resumes/current", hhHandler.GetCurrentResume)

		api.GET("/vacancies/similar", hhHandler.GetSimilarVacancies)
		api.GET("/vacancies/similar/first", hhHandler.GetFirstSimilarVacancy)
		api.GET("/vacancies/:vacancy_id", hhHandler.GetVacancyByID)
		api.POST("/vacancies/apply/:vacancy_id", applicationHandler.ApplyToVacancy)

		api.POST("/cover-letter", applicationHandler.GenerateCoverLetter)
	}
}
