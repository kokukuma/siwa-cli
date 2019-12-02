package siwa

import "net/url"

var (
	appleCom     = "https://appleid.apple.com"
	authEndpoint = "https://appleid.apple.com/auth/authorize"
)

// Config summarize information of SIWA
type Config struct {
	TeamID      string
	ClientID    string
	RedirectURI string
	AppleKey    string
	Scope       string
}

// CreateAuthURL returns url to get auth code
func (c Config) CreateAuthURL(state string) string {
	values := url.Values{}
	values.Add("client_id", c.ClientID)
	values.Add("redirect_uri", c.RedirectURI)
	values.Add("response_type", "code id_token")
	values.Add("state", state)
	values.Add("scope", c.Scope)
	values.Add("response_mode", "form_post")

	return authEndpoint + "?" + values.Encode()
}
