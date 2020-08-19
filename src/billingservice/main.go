package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"learn.oauth.billingservice/model"
)

var oauthServer = struct {
	introspectionURL string
}{
	introspectionURL: "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/token/introspect",
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

	}
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
		return accessToken, nil
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
	return true
}
