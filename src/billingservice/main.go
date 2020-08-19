package main

import (
	"encoding/json"
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
	clientID       string
	clientPassword string
}{
	clientID:       "oauth-nailed-app-1-token-checker",
	clientPassword: "7e9247c4-bcbe-4783-89f0-880d83ac147f",
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
		fmt.Println(nerr.Error())
		s := &model.BillingServicesErrorResponse{Error: nerr.Error()}
		rs.Header().Add("Content-Type", "application/json")
		rs.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(rs)
		encoder.Encode(s)
		return
	}
	fmt.Println("Access token provided:", accessToken)
	//
	if !isAccessTokenValid(accessToken) {
		fmt.Println("Access token is not valid or not active")
		rs.WriteHeader(http.StatusForbidden)
		return
	}
	fmt.Println("Access token is valid and active")
	//
	s := model.BillingServicesResponse{Services: []string{"electric", "phone", "internet", "water"}}
	rs.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(rs)
	encoder.Encode(s)
}

func getAccessToken(rq *http.Request) (string, error) {
	// On Authorization Request Header Field, see https://tools.ietf.org/html/rfc6750#section-2.1
	accessToken := rq.Header.Get("Authorization")
	if accessToken != "" {
		authorizationHeaderParts := strings.Split(accessToken, " ")
		if len(authorizationHeaderParts) != 2 {
			return accessToken, fmt.Errorf("Authorization header format is invalid")
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
	return accessToken, fmt.Errorf("Access token is not provided")
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
