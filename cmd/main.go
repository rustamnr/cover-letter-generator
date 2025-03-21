package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	_ "github.com/joho/godotenv/autoload"
)

const (
	hhAuthURL  = "https://hh.ru/oauth/authorize"
	hhTokenURL = "https://hh.ru/oauth/token"
	hhAPIURL   = "https://api.hh.ru"
)

var (
	clientID     string
	clientSecret string
	redirectURI  string
)

// HHClient - клиент для работы с hh.ru API
type HHClient struct {
	httpClient *resty.Client
}

// NewHHClient создает новый клиент
func NewHHClient() *HHClient {
	return &HHClient{
		httpClient: resty.New(),
	}
}

// 🔹 1. Авторизация: редирект на hh.ru
func (h *HHClient) AuthHandler(c *gin.Context) {
	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s",
		hhAuthURL, clientID, url.QueryEscape(redirectURI))
	c.Redirect(http.StatusFound, authURL)
}

// 🔹 2. Получение access_token
func (h *HHClient) CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	// Запрос на обмен `code` -> `access_token`
	resp, err := h.httpClient.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code":          code,
			"redirect_uri":  redirectURI,
		}).
		Post(hhTokenURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса к hh.ru"})
		return
	}

	if resp.StatusCode() != http.StatusOK {
		c.JSON(resp.StatusCode(), gin.H{"error": resp.String()})
		return
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &tokenData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ответа"})
		return
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access_token not found"})
		return
	}

	// Получаем user_id
	userID, err := h.getUserID(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения user_id"})
		return
	}

	// Сохраняем access_token в сессию
	session := sessions.Default(c)
	session.Set("access_token", accessToken)
	session.Set("user_id", userID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message":      "Успешная авторизация!",
		"user_id":      userID,
		"access_token": accessToken,
	})
}

// 🔹 3. Получение user_id (чтобы понимать, кто залогинен)
func (h *HHClient) getUserID(accessToken string) (string, error) {
	resp, err := h.httpClient.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Accept", "application/json").
		Get(hhAPIURL + "/me")

	if err != nil || resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("ошибка получения user_id")
	}

	var userData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &userData); err != nil {
		return "", err
	}

	userID, ok := userData["id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id не найден")
	}

	return userID, nil
}

// 🔹 4. Получение резюме залогиненного пользователя
func (h *HHClient) GetUserResumesHandler(c *gin.Context) {
	session := sessions.Default(c)
	accessToken := session.Get("access_token")
	userID := session.Get("user_id")

	if accessToken == nil || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	resp, err := h.httpClient.R().
		SetHeader("Authorization", "Bearer "+accessToken.(string)).
		SetHeader("Accept", "application/json").
		Get(hhAPIURL + "/resumes/mine")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func main() {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Не найдены переменные окружения")
	}

	hhClient := NewHHClient()
	router := gin.Default()

	// Настройка сессий с использованием cookie
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("my_session", store))

	// 🔹 Авторизация
	router.GET("/auth", hhClient.AuthHandler)
	router.GET("/auth/callback", hhClient.CallbackHandler)

	// 🔹 Доступ к резюме
	router.GET("/resumes/mine", hhClient.GetUserResumesHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}
