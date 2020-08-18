package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"learn.oauth.billingservice/model"
)

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
	s := model.BillingServicesResponse{Services: []string{"electric", "phone", "internet", "water"}}
	rs.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(rs)
	encoder.Encode(s)
}
