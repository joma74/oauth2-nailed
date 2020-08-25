package model

// AccessTokenPayload Payload of an access token
type AccessTokenPayload struct {
	Exp            int         `json:"exp"`
	Iat            int         `json:"iat"`
	AuthTime       int         `json:"auth_time"`
	Jti            string      `json:"jti"`
	Iss            string      `json:"iss"`
	Aud            interface{} `json:"aud"`
	Sub            string      `json:"sub"`
	Typ            string      `json:"typ"`
	Azp            string      `json:"azp"`
	SessionState   string      `json:"session_state"`
	Acr            string      `json:"acr"`
	AllowedOrigins []string    `json:"allowed-origins"`
	RealmAccess    struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess struct {
		Account struct {
			Roles []string `json:"roles"`
		} `json:"account"`
	} `json:"resource_access"`
	Scope             string `json:"scope"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

// AudAsSlice return audience as slice of string
func (a *AccessTokenPayload) AudAsSlice() []string {
	result := []string{}
	switch a.Aud.(type) {
	case string:
		result = []string{a.Aud.(string)}
	case []interface{}:
		audiences, ok := a.Aud.([]interface{})
		if !ok {
			return result
		}
		for _, audience := range audiences {
			if sAud, ok := audience.(string); ok {
				result = append(result, sAud)
			}
		}
	}
	return result
}
