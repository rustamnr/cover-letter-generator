package server

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/rustamnr/cover-letter-generator/internal/handlers"
	"github.com/rustamnr/cover-letter-generator/internal/middleware"
	"github.com/rustamnr/cover-letter-generator/internal/services"

	"github.com/gin-gonic/gin"
)

// registerRoutes настраивает маршруты API
func registerRoutes(router *gin.Engine) {
	// hhService := services.NewResumeService("https://api.hh.ru")
	// resumeHandler := handlers.NewResumeHandler(hhService)

	hhService := services.NewHHService(
		os.Getenv("HH_CLIENT_ID"),
		os.Getenv("HH_CLIENT_SECRET"),
		os.Getenv("HH_REDIRECT_URI"),
	)
	chatGPTService := services.NewChatGPTService(
		os.Getenv("CHATGPT_API_URL"),
		os.Getenv("CHATGPT_API_KEY"),
	)
	deepSeekService := services.NewDeepSeekService(
		os.Getenv("DEEPSEEK_API_URL"),
		os.Getenv("DEEPSEEK_API_KEY"),
	)

	hhHandler := handlers.NewHHHandler(hhService)

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

	chatGPTHandler := handlers.NewChatGPTHandler(chatGPTService)
	deepSeekHandler := handlers.NewDeepSeekHandler(deepSeekService)

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/resumes", hhHandler.GetUserResumes)
		api.POST("/resumes/select", hhHandler.SelectResume)
		api.GET("/resumes/current", hhHandler.GetCurrentResume)
		api.GET("/vacancy", hhHandler.GetVacancy)
		api.GET("/negotiations", hhHandler.GetUserApplications)
		api.GET("/negotiation/", hhHandler.GetUserFirstFoundedApplication)
		api.POST("/message", hhHandler.SendNewMessage)

		router.POST("/generate/chatgpt", chatGPTHandler.HandleChatGPT)
		api.POST("/deepseek", deepSeekHandler.HandleDeepSeek)
		api.POST("/generate/deepseek")
	}
}
