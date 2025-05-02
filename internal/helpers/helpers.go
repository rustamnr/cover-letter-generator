package helpers

import "github.com/rustamnr/cover-letter-generator/internal/constants"

func GetAuthURL(clientID, redirectURI string) string {
	return constants.HHURL + constants.Authorize + "?response_type=code&client_id=" +
		clientID + "&redirect_uri=" + redirectURI
}
