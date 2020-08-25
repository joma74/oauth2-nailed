package model

// IntrospectionRequestingPartyTokenResponse Response of token introspection in the form was requested as token_type_hint with value requesting_party_token
type IntrospectionRequestingPartyTokenResponse struct {
	Exp    int         `json:"exp"`
	Nbf    int         `json:"nbf"`
	Iat    int         `json:"iat"`
	Jti    string      `json:"jti"`
	Aud    interface{} `json:"aud"`
	Typ    string      `json:"typ"`
	Acr    string      `json:"acr"`
	Active bool        `json:"active"`
}
