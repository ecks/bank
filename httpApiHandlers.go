package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bvnk/bank/accounts"
	"github.com/bvnk/bank/appauth"
	"github.com/bvnk/bank/transactions"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API Index")
}

func getTokenFromHeader(w http.ResponseWriter, r *http.Request) (token string, err error) {
	// Get token from header
	token = r.Header.Get("X-Auth-Token")
	if token == "" {
		return "", errors.New("httpApiHandlers: Could not retrieve token from headers")
	}

	// Check token
	err = appauth.CheckToken(token)
	if err != nil {
		return "", errors.New("httpApiHandlers: Token invalid")
	}

	return
}

// Extend token
func AuthIndex(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	//Extend token
	response, err := appauth.ProcessAppAuth([]string{token, "appauth", "1"})
	fmt.Println(response)
	fmt.Println(err)
	Response(response, err, w, r)
	return
}

// Get token
func AuthLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get token")
	user := r.FormValue("User")
	password := r.FormValue("Password")

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", user, password})
	Response(response, err, w, r)
	return
}

// Create auth account
func AuthCreate(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("UserIdentificationNumber")
	password := r.FormValue("Password")

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "3", userID, password})
	Response(response, err, w, r)
	return
}

// Remove auth account
func AuthRemove(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	user := r.FormValue("User")
	password := r.FormValue("Password")

	response, err := appauth.ProcessAppAuth([]string{token, "appauth", "4", user, password})
	Response(response, err, w, r)
	return
}

func AccountIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Account Index")
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1001"})
	Response(response, err, w, r)
	return
}

func AccountCreate(w http.ResponseWriter, r *http.Request) {
	// Get values from POST
	accountHolderGivenName := r.FormValue("AccountHolderGivenName")
	accountHolderFamilyName := r.FormValue("AccountHolderFamilyName")
	accountHolderDateOfBirth := r.FormValue("AccountHolderDateOfBirth")
	accountHolderIdentificationNumber := r.FormValue("AccountHolderIdentificationNumber")
	accountHolderContactNumber1 := r.FormValue("AccountHolderContactNumber1")
	accountHolderContactNumber2 := r.FormValue("AccountHolderContactNumber2")
	accountHolderEmailAddress := r.FormValue("AccountHolderEmailAddress")
	accountHolderAddressLine1 := r.FormValue("AccountHolderAddressLine1")
	accountHolderAddressLine2 := r.FormValue("AccountHolderAddressLine2")
	accountHolderAddressLine3 := r.FormValue("AccountHolderAddressLine3")
	accountHolderPostalCode := r.FormValue("AccountHolderPostalCode")
	accountType := r.FormValue("AccountType")

	req := []string{
		"0",
		"acmt",
		"1",
		accountHolderGivenName,
		accountHolderFamilyName,
		accountHolderDateOfBirth,
		accountHolderIdentificationNumber,
		accountHolderContactNumber1,
		accountHolderContactNumber2,
		accountHolderEmailAddress,
		accountHolderAddressLine1,
		accountHolderAddressLine2,
		accountHolderAddressLine3,
		accountHolderPostalCode,
		accountType,
	}

	response, err := accounts.ProcessAccount(req)
	Response(response, err, w, r)
	return
}

func AccountGet(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1002", accountId})
	Response(response, err, w, r)
	return
}

func AccountRetrieve(w http.ResponseWriter, r *http.Request) {
	// Set these in the header as they are sensitive
	ID := r.Header.Get("X-IDNumber")
	givenName := r.Header.Get("X-GivenName")
	familyName := r.Header.Get("X-FamilyName")
	email := r.Header.Get("X-EmailAddress")

	response, err := accounts.ProcessAccount([]string{"", "acmt", "1006", ID, givenName, familyName, email})
	Response(response, err, w, r)
	return
}

func AccountGetAll(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1000"})
	Response(response, err, w, r)
	return
}

func AccountTokenPost(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	pushToken := r.FormValue("PushToken")
	platform := r.FormValue("Platform")

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1003", pushToken, platform})
	Response(response, err, w, r)
	return
}

func AccountTokenDelete(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	pushToken := r.FormValue("PushToken")
	platform := r.FormValue("Platform")

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1004", pushToken, platform})
	Response(response, err, w, r)
	return
}

func AccountSearch(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	searchTerm := r.FormValue("Search")

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1005", searchTerm})
	Response(response, err, w, r)
	return
}

func TransactionCreditInitiation(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	senderDetails := r.FormValue("SenderDetails")
	recipientDetails := r.FormValue("RecipientDetails")
	amount := r.FormValue("Amount")
	lat := r.FormValue("Lat")
	lon := r.FormValue("Lon")
	desc := r.FormValue("Desc")

	response, err := transactions.ProcessPAIN([]string{token, "pain", "1", senderDetails, recipientDetails, amount, lat, lon, desc})
	Response(response, err, w, r)
	return
}

func TransactionDepositInitiation(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}

	accountDetails := r.FormValue("AccountDetails")
	amount := r.FormValue("Amount")
	lat := r.FormValue("Lat")
	lon := r.FormValue("Lon")
	desc := r.FormValue("Desc")

	response, err := transactions.ProcessPAIN([]string{token, "pain", "1000", accountDetails, amount, lat, lon, desc})
	Response(response, err, w, r)
	return
}

func TransactionList(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(w, r)
	if err != nil {
		Response("", err, w, r)
		return
	}
	// Get account number from header
	accountNumber := r.Header.Get("X-Auth-AccountNumber")
	if accountNumber == "" {
		Response("", errors.New("httpApiHandlers.TransactionList: Could not retrieve accountNumber from headers"), w, r)
		return
	}

	vars := mux.Vars(r)
	perPage := vars["perPage"]
	page := vars["page"]
	timestamp := vars["timestamp"]

	response, err := transactions.ProcessPAIN([]string{token, "pain", "1001", accountNumber, page, perPage, timestamp})
	Response(response, err, w, r)
	return
}
