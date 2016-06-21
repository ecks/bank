package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ksred/bank/accounts"
	"github.com/ksred/bank/appauth"
	"github.com/ksred/bank/configuration"
	"github.com/ksred/bank/payments"
	"github.com/ksred/bank/push"
)

func RunHttpServer() (err error) {
	fmt.Println("HTTP Server called")

	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		return errors.New("server.runServer: " + err.Error())
	}

	// Set config in packages
	accounts.SetConfig(&Config)
	payments.SetConfig(&Config)
	appauth.SetConfig(&Config)
	push.SetConfig(&Config)

	router := NewRouter()

	//err = http.ListenAndServeTLS(":8443", "certs/server.pem", "certs/server.key", router)
	err = http.ListenAndServeTLS(":8443", "certs/thebankoftoday.com.crt", "certs/thebankoftoday.com.key", router)
	fmt.Println(err)
	return
}

func Response(responseSuccess string, responseError error, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	req := make(map[string]string)

	// Check for error
	if responseError != nil {
		req["error"] = responseError.Error()
		jsonResponse, err := json.Marshal(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{error: 'Could not parse response'}"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// Create response
	req["response"] = string(responseSuccess)
	jsonResponse, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{error: 'Could not parse response'}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
