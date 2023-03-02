package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var certFile = "/certs/tls.crt"
var keyFile = "/certs/tls.key"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/check-namespace", CheckNamespace)
	errs := make(chan error, 1)

	go func() {
		errs <- http.ListenAndServe(":80", mux)
	}()

	go func() {
		errs <- http.ListenAndServeTLS(":443", certFile, keyFile, mux)
	}()

	log.Println("Listening...")
	<-errs
}

func CheckNamespace(w http.ResponseWriter, r *http.Request) {
	status, data := parseNamespace(r)
	if status != http.StatusOK {
		http.Error(w, data, status)
		return
	}

	if data == "" || data == "default" {
		w.WriteHeader(200)
		w.Write([]byte("{\"allowed\": false}"))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("{\"allowed\": true}"))
	}
}

func parseNamespace(r *http.Request) (int, string) {
	switch r.Method {
	case http.MethodGet:
		return http.StatusOK, r.URL.Query().Get("namespace")

	case http.MethodPost:
		var requestData map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			return http.StatusBadRequest, err.Error()
		}

		return http.StatusOK, requestData["namespace"]

	default:
		return http.StatusMethodNotAllowed, "Method not allowed"
	}
}
