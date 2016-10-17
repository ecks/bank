package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bvnk/bank/accounts"
	"github.com/bvnk/bank/appauth"
	"github.com/bvnk/bank/configuration"
	"github.com/bvnk/bank/push"
	"github.com/bvnk/bank/transactions"
)

func RunHttpServer() (err error) {
	bLog(0, "HTTP server called", trace())

	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		return errors.New("server.runServer: " + err.Error())
	}

	// Set config in packages
	accounts.SetConfig(&Config)
	transactions.SetConfig(&Config)
	appauth.SetConfig(&Config)
	push.SetConfig(&Config)

	router := NewRouter()

	err = http.ListenAndServeTLS(":"+Config.HttpPort, configuration.ImportPath+"certs/"+Config.FQDN+".pem", configuration.ImportPath+"certs/"+Config.FQDN+".key", router)
	//err = http.ListenAndServeTLS(":8443", "certs/thebankoftoday.com.crt", "certs/thebankoftoday.com.key", router)
	fmt.Println(err)
	bLog(4, err.Error(), trace())
	return
}

func Response(responseSuccess interface{}, responseError error, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	req := make(map[string]interface{})

	// Check for error
	if responseError != nil {
		req["error"] = responseError.Error()
		bLog(3, "Response error: "+responseError.Error(), trace())
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

	req["response"] = responseSuccess
	jsonResponse, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{error: 'Could not parse response'}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	bLog(0, "Response success: "+string(jsonResponse), trace())
}
