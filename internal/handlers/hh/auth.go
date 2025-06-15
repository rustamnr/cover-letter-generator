package handlers_hh

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
)

// AuthHandler redirects user to the authorization page
func (h *HHHandler) AuthHandler(c *gin.Context) {
	authURL := fmt.Sprintf("https://hh.ru/oauth/authorize?response_type=code&client_id=%s", h.hhClient.ClientID)
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler handles the OAuth callback
func (h *HHHandler) CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not found"})
		return
	}

	accessToken, err := h.hhClient.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.hhClient.GetUserID(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set(constants.AccessToken, accessToken)
	session.Set(constants.UserId, userID)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hhClient.SetAccessToken(accessToken)

	c.JSON(http.StatusOK, gin.H{"message": "authorized", "user_id": userID, "access_token": accessToken})
}

func GetAuthURL(clientID, redirectURI string) string {
	return constants.HHURL + constants.Authorize + "?response_type=code&client_id=" +
		clientID + "&redirect_uri=" + redirectURI
}
