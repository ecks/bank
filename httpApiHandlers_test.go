package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bvnk/bank/accounts"
	"github.com/bvnk/bank/appauth"
	"github.com/bvnk/bank/configuration"
	"github.com/bvnk/bank/payments"
	"github.com/bvnk/bank/push"
)

func TestLoadConfiguration(t *testing.T) {
	// Load app config
	_, err := configuration.LoadConfig()
	if err != nil {
		t.Errorf("loadDatabase does not pass. Configuration does not load, looking for %v, got %v", nil, err)
	}
}

func loadAllConfig(t *testing.T) {
	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		t.Errorf("loadDatabase does not pass. Configuration does not load, looking for %v, got %v", nil, err)
	}

	// Set config in packages
	accounts.SetConfig(&Config)
	payments.SetConfig(&Config)
	appauth.SetConfig(&Config)
	push.SetConfig(&Config)
}

func TestIndex(t *testing.T) {
	loadAllConfig(t)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Index)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//    t.Errorf("handler returned unexpected body: got %v want %v",
	//        rr.Body.String(), expected)
	//}
}

// Extend token
func TestAuthIndexFailNoToken(t *testing.T) {
	loadAllConfig(t)

	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAuthIndexFailIncorrectToken(t *testing.T) {
	loadAllConfig(t)

	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Auth-Token", "incorrect-token")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

/*
func TestAccountCreate(t *testing.T) {
	loadAllConfig(t)

	user := map[string]string{
		"AccountHolderGivenName":            "Test",
		"AccountHolderFamilyName":           "Account",
		"AccountHolderDateOfBirth":          "01011990",
		"AccountHolderIdentificationNumber": "1234567890",
		"AccountHolderContactNumber1":       "180012345",
		"AccountHolderContactNumber2":       "",
		"AccountHolderEmailAddress":         "test@email.com",
		"AccountHolderAddressLine1":         "Address Line 1",
		"AccountHolderAddressLine2":         "Address Line 2",
		"AccountHolderAddressLine3":         "Address Line 3",
		"AccountHolderPostalCode":           "AB1 CD2",
	}

	req, err := http.NewRequest("POST", "/account", user)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
*/
