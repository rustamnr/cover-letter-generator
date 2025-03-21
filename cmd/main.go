package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const (
	hhAuthURL        = "https://hh.ru/oauth/authorize"
	hhTokenURL       = "https://hh.ru/oauth/token"
	hhResumesMineURL = "https://api.hh.ru/resumes/mine"
	redirectURI      = "http://localhost:8080/auth/callback"
)

var (
	clientID     string
	clientSecret string
)

type HHClient struct {
	httpClient  *resty.Client
	accessToken string
}

func NewHHClient() *HHClient {
	return &HHClient{
		httpClient: resty.New(),
	}
}

// 🔹 Авторизация пользователя через hh.ru
func (h *HHClient) AuthHandler(c *gin.Context) {
	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s",
		hhAuthURL, clientID, url.QueryEscape(redirectURI))
	c.Redirect(http.StatusFound, authURL)
}

// 🔹 Обработка редиректа от hh.ru и обмен кода на access_token
func (h *HHClient) AuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	token, err := h.getAccessToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		log.Println("Error getting access token:", err)
		return
	}

	h.accessToken = token
	c.JSON(http.StatusOK, gin.H{"message": "User authenticated", "access_token": token})
}

// 🔹 Обмен кода на access_token
func (h *HHClient) getAccessToken(code string) (string, error) {
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
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %s", resp.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response")
	}

	return accessToken, nil
}

// 🔹 Получение резюме пользователя
// 🔹 Обработчик запроса к API hh.ru
func (h *HHClient) GetUserResumesHandler(c *gin.Context) {
	if h.accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	resumes, err := h.getUserResumes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user resumes"})
		return
	}

	c.JSON(http.StatusOK, resumes)
}

// 🔹 Внутренний метод для работы с API hh.ru
func (h *HHClient) getUserResumes() (map[string]interface{}, error) {
	resp, err := h.httpClient.R().
		SetHeader("Authorization", "Bearer "+h.accessToken).
		SetHeader("Accept", "application/json").
		Get(hhResumesMineURL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("ошибка запроса: %s", resp.String())
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 🔹 Запуск сервера
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// 🔹 Читаем переменные
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	hhClient := NewHHClient()
	router := gin.Default()

	router.GET("/", hhClient.AuthHandler)
	router.GET("/auth/callback", hhClient.AuthCallbackHandler)
	router.GET("/resumes/main", hhClient.GetUserResumesHandler)

	log.Println("Server running on http://localhost:" + port)
	log.Fatal(router.Run(":" + port))
}
