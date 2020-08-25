package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"learn.oauth.billingservice/model"
)

var oauthServer = struct {
	introspectionURL string
}{
	introspectionURL: "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/token/introspect",
}

var oauthClient = struct {
	clientID                   string
	clientPassword             string
	scopeNameBillingservice    string
	audienceNameBillingservice string
}{
	clientID:                   "oauth-nailed-app-1-token-checker",
	clientPassword:             "7e9247c4-bcbe-4783-89f0-880d83ac147f",
	scopeNameBillingservice:    "billingService",
	audienceNameBillingservice: "billingService",
}

func main() {
	fmt.Println("Server starting")
	http.HandleFunc("/", withMethodLogging(home))
	http.HandleFunc("/billing/v1/services", withMethodLogging(services))
	http.ListenAndServe(":9111", nil)
}

func withMethodLogging(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rs http.ResponseWriter, rq *http.Request) {
		methodName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		fmt.Printf("\033[1;36m%s\033[0m\n", "--> "+methodName)
		handler(rs, rq)
		fmt.Printf("\033[1;36m%s\033[0m\n", "<-- "+methodName)
	}
}

func home(rs http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(rs, "hello")
}

func services(rs http.ResponseWriter, rq *http.Request) {
	accessToken, nerr := getAccessToken(rq)
	if nerr != nil {
		sendErrorResponseMessage(nerr, http.StatusBadRequest, rs)
		return
	}
	fmt.Println("Access token provided:", accessToken)
	//
	if !isAccessTokenValid(accessToken) {
		nerr := errors.New("Access token is not valid or not active")
		sendErrorResponseMessage(nerr, http.StatusForbidden, rs)
		return
	}
	fmt.Println("Access token is valid and active")
	//
	payload, nerr := getPayload(accessToken)
	if nerr != nil {
		nerr := fmt.Errorf("Could not decode paylod from accessToken: %v", nerr)
		sendErrorResponseMessage(nerr, http.StatusBadRequest, rs)
		return
	}

	accessTokenPayload := &model.AccessTokenPayload{}
	nerr = json.Unmarshal(payload, accessTokenPayload)
	if nerr != nil {
		fmt.Println("Could not unmarshal accessToken payload to model.AccessTokenPayload", nerr)
		sendErrorResponseMessage(nerr, http.StatusBadRequest, rs)
		return
	}
	//
	scopes := strings.Split(accessTokenPayload.Scope, " ")
	fmt.Println("Scopes from accessToken payload are")
	for _, scope := range scopes {
		fmt.Println(" - ", scope)
	}
	if !strings.Contains(accessTokenPayload.Scope, oauthClient.scopeNameBillingservice) {
		nerr := fmt.Errorf("Denied access because a access token's scope of '" + oauthClient.scopeNameBillingservice + "' is required but was not provided.")
		sendErrorResponseMessage(nerr, http.StatusForbidden, rs)
		return
	}
	//
	isValidAudience := false
	fmt.Println("Audiences from accessToken payload are")
	for _, audience := range accessTokenPayload.AudAsSlice() {
		fmt.Println(" - ", audience)
		if audience == oauthClient.audienceNameBillingservice {
			isValidAudience = true
		}
	}
	if !isValidAudience {
		nerr := fmt.Errorf("Denied access because a access token's audience of '" + oauthClient.audienceNameBillingservice + "' is required but was not provided.")
		sendErrorResponseMessage(nerr, http.StatusForbidden, rs)
		return
	}
	//
	s := model.BillingServicesResponse{Services: []string{"electric", "phone", "internet", "water"}}
	rs.Header().Add("Content-Type", "application/json")
	rs.Header().Add("Access-Control-Allow-Origin", "*")
	encoder := json.NewEncoder(rs)
	encoder.Encode(s)
}

func getAccessToken(rq *http.Request) (string, error) {
	// On Authorization Request Header Field, see https://tools.ietf.org/html/rfc6750#section-2.1
	accessToken := rq.Header.Get("Authorization")
	if accessToken != "" {
		authorizationHeaderParts := strings.Split(accessToken, " ")
		if len(authorizationHeaderParts) != 2 {
			return accessToken, errors.New("Authorization header format is invalid")
		}
		return authorizationHeaderParts[1], nil
	}
	// On Form-Encoded Body Parameter, see https://tools.ietf.org/html/rfc6750#section-2.2
	accessToken = rq.FormValue("access_token")
	if accessToken != "" {
		return accessToken, nil
	}
	// On URI Query Parameter, see https://tools.ietf.org/html/rfc6750#section-2.3
	accessToken = rq.URL.Query().Get("access_token")
	if accessToken != "" {
		return accessToken, nil
	}
	return accessToken, errors.New("Access token is not provided")
}

func isAccessTokenValid(accessToken string) bool {
	form := url.Values{}
	form.Add("token", accessToken)
	form.Add("token_type_hint", "requesting_party_token")
	nrq, nerr := http.NewRequest("POST", oauthServer.introspectionURL, strings.NewReader(form.Encode()))
	nrq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	nrq.SetBasicAuth(oauthClient.clientID, oauthClient.clientPassword)
	if nerr != nil {
		fmt.Print("Could not prepare request", nerr)
		return false
	}
	c := http.Client{}
	nrs, nerr := c.Do(nrq)
	if nerr != nil {
		fmt.Println("Could not send introspection request", nerr)
		return false
	}
	if nrs.StatusCode != 200 {
		fmt.Println("Introspection response status is expected to be 200 but was ", nrs.StatusCode)
		return false
	}
	byteBody, nerr := ioutil.ReadAll(nrs.Body)
	defer nrs.Body.Close()
	if nerr != nil {
		fmt.Println("Could not read introspection response body", nerr)
		return false
	}

	introspectionRequestingPartyTokenResponse := &model.IntrospectionRequestingPartyTokenResponse{}
	nerr = json.Unmarshal(byteBody, introspectionRequestingPartyTokenResponse)
	if nerr != nil {
		fmt.Println("Could not unmarshal response to model.IntrospectionRequestingPartyTokenResponse", nerr)
		return false
	}

	return introspectionRequestingPartyTokenResponse.Active
}

func getPayload(accessToken string) ([]byte, error) {
	tokenParts := strings.Split(accessToken, ".")
	payload, nerr := base64.RawURLEncoding.DecodeString(tokenParts[1])
	return payload, nerr
}

func sendErrorResponseMessage(nerr error, statusCode int, rs http.ResponseWriter) {
	fmt.Println(nerr.Error())
	s := &model.BillingServicesErrorResponse{Error: nerr.Error()}
	rs.Header().Add("Content-Type", "application/json")
	rs.Header().Add("Access-Control-Allow-Origin", "*")
	rs.WriteHeader(statusCode)
	encoder := json.NewEncoder(rs)
	encoder.Encode(s)
}
