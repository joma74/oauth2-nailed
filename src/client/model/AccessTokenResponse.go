package model

// AccessTokenResponse Access token from token call
type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	NotBeforePolicy  int64  `json:"not-before-policy"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	SessionState     string `json:"session_state"`
	TokenType        string `json:"token_type"`
}
