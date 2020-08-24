package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadDecoder(t *testing.T) {
	assert := assert.New(t)
	//
	accessToken := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJaNFVObHQ0alZ5STNjUnNTZW9rMjdza2xvM0ZGMzVOaG5nWlA5NEk1R2xnIn0.eyJleHAiOjE1OTgyNTU5MjMsImlhdCI6MTU5ODI1NTYyMywiYXV0aF90aW1lIjoxNTk4MjUzODM0LCJqdGkiOiJmYzc0MTU2MC1mM2U4LTQzODMtYTk1Ni05NmVmMDY4MzE2ZmMiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjkxMTIvYXV0aC9yZWFsbXMvbXlyZWFsbSIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiI1NTdjYTc0OS02ODhjLTQyZDItYjhiMS05NTNiYzNlZDJmYTciLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJvYXV0aC1uYWlsZWQtYXBwLTEiLCJzZXNzaW9uX3N0YXRlIjoiMWEwZmYxNWYtNTYwYi00NmZhLThlMzUtMDI1YTc5OWIzNmE4IiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vbG9jYWxob3N0OjkxMTAiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJteSB1c2VyIiwicHJlZmVycmVkX3VzZXJuYW1lIjoibXl1c2VyIiwiZ2l2ZW5fbmFtZSI6Im15IiwiZmFtaWx5X25hbWUiOiJ1c2VyIn0.Kqmhfcj30frZlmlZj0OXj-TzgMSDotmTcknA_RADuS62U28D3uEnx_i2XncJDE8yMQJjie-95I2a3xvjQzcfVogei_fV99jX29ZRUxmNJkX6VbYP2BOmlF1cepWoMhIzrhCWIgjDOcTmHbbiMZGe4-uTpYyDetErduNR3jBSiZ2ajnP_3L9ZEMskLYs74hwghapqvYvMEQUI9iZdUpFEkkbx1yRV-bX96u8akdEY4qwOmGbZkjklh91sYrY44G4-ErQxcIxnoBVNVQW3I6LaW9XN5Hpw0BNcyX6itHzlK0Jqn7msP68ab_087diQIuyk1-TSh52YyVRc2_VTvly_Zw"
	tokenParts := strings.Split(accessToken, ".")
	claims, nerr := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if nerr != nil {
		nerr := fmt.Errorf("Could not decode paylod from accessToken: %v", nerr)
		t.Error(nerr)
	}
	expected := `{"exp":1598255923,"iat":1598255623,"auth_time":1598253834,"jti":"fc741560-f3e8-4383-a956-96ef068316fc","iss":"http://localhost:9112/auth/realms/myrealm","aud":"account","sub":"557ca749-688c-42d2-b8b1-953bc3ed2fa7","typ":"Bearer","azp":"oauth-nailed-app-1","session_state":"1a0ff15f-560b-46fa-8e35-025a799b36a8","acr":"0","allowed-origins":["http://localhost:9110"],"realm_access":{"roles":["offline_access","uma_authorization"]},"resource_access":{"account":{"roles":["manage-account","manage-account-links","view-profile"]}},"scope":"profile email","email_verified":false,"name":"my user","preferred_username":"myuser","given_name":"my","family_name":"user"}`
	actual := string(claims)
	assert.Equal(expected, actual)
}
