package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/rustamnr/cover-letter-generator/internal/clients"
	"github.com/rustamnr/cover-letter-generator/internal/handlers"
	"github.com/rustamnr/cover-letter-generator/internal/middleware"
	"github.com/rustamnr/cover-letter-generator/internal/services"

	"github.com/gin-gonic/gin"
)

// registerRoutes настраивает маршруты API
func registerRoutes(router *gin.Engine) {
	// Инициализация клиентов
	hhClient := clients.NewHHClient()
	deepSeekClient := clients.NewDeepSeekClient()

	vacancyProvider := services.NewHHProvider(hhClient)
	textGenerator := services.NewDeepSeekService(deepSeekClient)
	// Инициализация сервисов
	applicationService := services.NewApplicationService(vacancyProvider, textGenerator)

	// Инициализация хендлеров
	hhHandler := handlers.NewHHHandler(hhClient)
	applicationHandler := handlers.NewApplicationHandler(applicationService)

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
		api.POST("/resumes/select", hhHandler.SetCurrnetResume)
		api.GET("/resumes/current", hhHandler.GetCurrentResume)

		api.GET("/vacancies/similar", hhHandler.GetSimilarVacancies)
		api.GET("/vacancies/similar/first", hhHandler.GetFirstSimilarVacancy)
		api.GET("/vacancies/:vacancy_id", hhHandler.GetVacancyByID)
		api.GET("/vacancy", hhHandler.GetVacancyByID)

		api.POST("/cover-letter", applicationHandler.GenerateCoverLetter)
	}
}
