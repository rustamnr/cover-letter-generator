package constants

const (
	HHURL                  = "https://hh.ru"
	Authorize              = "/oauth/authorize"
	Me                     = "/me"
	ResumesMine            = "/resumes/mine"
	Negotiations           = "/negotiations"
	NegotiationsNidMessage = Negotiations + "/%s/messages"
)

func GetAuthURL(clientID, redirectURI string) string {
	return HHURL + Authorize + "?response_type=code&client_id=" + clientID + "&redirect_uri=" + redirectURI
}
