package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var oauthServer = struct {
	authURL   string
	logoutURL string
}{
	authURL:   "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/auth",
	logoutURL: "http://localhost:9112/auth/realms/myrealm/protocol/openid-connect/logout",
}

var oauthClient = struct {
	afterAuthURL   string
	afterLogoutURL string
}{
	afterAuthURL:   "http://localhost:9110/authCodeRedirect",
	afterLogoutURL: "http://localhost:9110/",
}

var t = template.Must(template.ParseFiles("template/index.html"))

var authCodeVars = struct {
	Code         string
	SessionState string
}{Code: "???", SessionState: "???"}

func main() {
	fmt.Println("Server starting")
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authCodeRedirect", authCodeRedirect)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":9110", nil)
}

func home(rs http.ResponseWriter, rq *http.Request) {
	t.Execute(rs, authCodeVars)
}

func login(rs http.ResponseWriter, rq *http.Request) {
	nrq, nerr := http.NewRequest("GET", oauthServer.authURL, nil)
	if nerr != nil {
		log.Print(nerr)
		return
	}
	qs := url.Values{}
	qs.Add("client_id", "oauth-nailed-app-1")
	qs.Add("response_type", "code")
	qs.Add("state", "123")
	qs.Add("redirect_uri", oauthClient.afterAuthURL)
	nrq.URL.RawQuery = qs.Encode()
	http.Redirect(rs, rq, nrq.URL.String(), http.StatusFound)
}

/**
* Location: http://localhost:9110/authCodeRedirect?state=123&session_state=6c634b86-8a30-...beaf&code=a16dcfbc-d53b-...-a66dbcfac9c1
 */
func authCodeRedirect(rs http.ResponseWriter, rq *http.Request) {
	fmt.Printf("Request query: %v\n", rq.URL.Query())
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
	http.Redirect(rs, rq, nrq.URL.String(), http.StatusFound)
}
