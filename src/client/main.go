package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"

	"learn.oauth.client/model"
)

var oauthServer = struct {
	authURL   string
	tokenURL  string
	logoutURL string
}{
	authURL:   "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/auth",
	tokenURL:  "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/token",
	logoutURL: "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/logout",
}

var oauthClient = struct {
	clientID       string
	clientPassword string
	afterAuthURL   string
	afterLogoutURL string
}{
	clientID:       "oauth-nailed-app-1",
	clientPassword: "0c061d83-f4f6-4678-94aa-5dc8d9584eea",
	afterAuthURL:   "http://localhost:9110/authCodeRedirect",
	afterLogoutURL: "http://localhost:9110/",
}

var t = template.Must(template.ParseFiles("template/index.html"))

var authCodeVars = struct {
	Code         string
	SessionState string
	AccessToken  string
	RefreshToken string
	TokenScope   string
}{Code: "???", SessionState: "???", AccessToken: "???", RefreshToken: "???", TokenScope: "???"}

func main() {
	fmt.Println("Server starting")
	http.HandleFunc("/", withMethodLogging(home))
	http.HandleFunc("/login", withMethodLogging(login))
	http.HandleFunc("/accessToken", withMethodLogging(accessToken))
	http.HandleFunc("/authCodeRedirect", withMethodLogging(authCodeRedirect))
	http.HandleFunc("/logout", withMethodLogging(logout))
	http.ListenAndServe(":9110", nil)
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
	t.Execute(rs, authCodeVars)
}

func login(rs http.ResponseWriter, rq *http.Request) {
	nrq, nerr := http.NewRequest("GET", oauthServer.authURL, nil)
	if nerr != nil {
		fmt.Print(nerr)
		return
	}
	qs := url.Values{}
	qs.Add("client_id", oauthClient.clientID)
	qs.Add("response_type", "code")
	qs.Add("state", "123")
	qs.Add("redirect_uri", oauthClient.afterAuthURL)
	nrq.URL.RawQuery = qs.Encode()
	http.Redirect(rs, rq, nrq.URL.String(), http.StatusFound)
}

func accessToken(rs http.ResponseWriter, rq *http.Request) {
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCodeVars.Code)
	form.Add("redirect_uri", oauthClient.afterAuthURL)
	form.Add("client_id", oauthClient.clientID)
	nrq, nerr := http.NewRequest("POST", oauthServer.tokenURL, strings.NewReader(form.Encode()))
	nrq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	nrq.SetBasicAuth(oauthClient.clientID, oauthClient.clientPassword)
	if nerr != nil {
		fmt.Println("Could not create new request", nerr)
		return
	}
	c := http.Client{}
	nrs, nerr := c.Do(nrq)
	if nerr != nil {
		fmt.Println("Could not get access token", nerr)
		return
	}
	byteBody, nerr := ioutil.ReadAll(nrs.Body)
	defer nrs.Body.Close()
	if nerr != nil {
		fmt.Println("Could not read body", nerr)
		return
	}
	accessTokenResponse := &model.AccessTokenResponse{}
	nerr = json.Unmarshal(byteBody, accessTokenResponse)
	if nerr != nil {
		fmt.Println("Could not unmarshal response to model.AccessTokenResponse", nerr)
		return
	}
	authCodeVars.AccessToken = accessTokenResponse.AccessToken
	authCodeVars.RefreshToken = accessTokenResponse.RefreshToken
	authCodeVars.TokenScope = accessTokenResponse.Scope
	//
	var out bytes.Buffer
	nerr = json.Indent(&out, byteBody, "", "   ")
	if nerr != nil {
		fmt.Println("Could not pretty print response", nerr)
		return
	}
	fmt.Printf("Access token response: %v\n", out.String())
	//
	http.Redirect(rs, rq, "/", http.StatusFound)
}

/**
* Location: http://localhost:9110/authCodeRedirect?state=123&session_state=6c634b86-8a30-...beaf&code=a16dcfbc-d53b-...-a66dbcfac9c1
 */
func authCodeRedirect(rs http.ResponseWriter, rq *http.Request) {
	fmt.Printf("After authentication the delivered data from Keycloak are:\n%v\n", rq.URL.Query())
	authCodeVars.Code = rq.URL.Query().Get("code")
	authCodeVars.SessionState = rq.URL.Query().Get("session_state")
	http.Redirect(rs, rq, "/", http.StatusFound)
}

func logout(rs http.ResponseWriter, rq *http.Request) {
	nrq, nerr := http.NewRequest("GET", oauthServer.logoutURL, nil)
	if nerr != nil {
		log.Print(nerr)
		return
	}
	qs := url.Values{}
	qs.Add("redirect_uri", oauthClient.afterLogoutURL)
	nrq.URL.RawQuery = qs.Encode()
	authCodeVars.Code = "???"
	authCodeVars.SessionState = "???"
	authCodeVars.AccessToken = "???"
	authCodeVars.RefreshToken = "???"
	authCodeVars.TokenScope = "???"
	http.Redirect(rs, rq, nrq.URL.String(), http.StatusFound)
}
