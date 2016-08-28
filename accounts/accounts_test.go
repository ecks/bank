package accounts

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestProcessAccountTooFewFields(t *testing.T) {
	tst := []string{"", ""}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "Data string does not have enough fields", nil)
	}
}

func TestProcessAccountACMTTypeNotSet(t *testing.T) {
	tst := []string{"", "", ""}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "Could not get type of ACMT transaction", nil)
	}
}

func TestProcessAccountACMTTypeIncorrect(t *testing.T) {
	tst := []string{"", "", "-1000"}
	_, err := ProcessAccount(tst)

	if err == nil {
		t.Errorf("ProcessAccount does not pass. Looking for %v, got %v", "ACMT transaction code invalid", nil)
	}
}

//@TODO Implement valid ACMT tests

func TestOpenCloseAccount(t *testing.T) {
	tst := []string{"", "", ""}
	_, err := openAccount(tst)

	if err == nil {
		t.Errorf("OpenAccount does not pass. Looking for %v, got %v", "Not all fields present", nil)
	}

	_, err = closeAccount(tst)

	if err == nil {
		t.Errorf("CloseAccount does not pass. Looking for %v, got %v", "Not all fields present", nil)
	}
}

func TestSetAccountDetails(t *testing.T) {
	tst := []string{
		"", // 0
		"", // acmt
		"", // 1
		"John",
		"Doe",
		"1900-01-01",
		"19000101-1000-100",
		"555-123-1234",
		"",
		"test@user.com",
		"Address 1",
		"Address 2",
		"Address 3",
		"22202",
		"cheque",
	}

	accountDetails, err := setAccountDetails(tst)

	if err != nil {
		t.Errorf("SetAccountDetails does not pass. ERROR. Looking for %v, got %v", nil, err)
	}

	if reflect.TypeOf(accountDetails).String() != "accounts.AccountDetails" {
		t.Errorf("SetAccountDetails does not pass. TYPE. Looking for %v, got %v", "accounts.AccountDetails", reflect.TypeOf(accountDetails).String())
	}

	if !accountDetails.Overdraft.Equals(decimal.NewFromFloat(OPENING_OVERDRAFT)) {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", decimal.NewFromFloat(OPENING_OVERDRAFT), accountDetails.Overdraft)
	}

	if !accountDetails.AccountBalance.Equals(decimal.NewFromFloat(OPENING_BALANCE)) {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", decimal.NewFromFloat(OPENING_BALANCE), accountDetails.AccountBalance)
	}

	if !accountDetails.AvailableBalance.Equals(decimal.NewFromFloat(OPENING_BALANCE + OPENING_OVERDRAFT)) {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", decimal.NewFromFloat(OPENING_BALANCE+OPENING_OVERDRAFT), accountDetails.AvailableBalance)
	}

	if accountDetails.AccountHolderName != "Doe,John" {
		t.Errorf("SetAccountDetails does not pass. DETAILS. Looking for %v, got %v", "Doe,John", accountDetails.AccountHolderName)
	}
}

func TestSetAccountHolderDetailsFailure(t *testing.T) {
	tst := []string{"", "", "", "John", "Doe"}
	_, err := setAccountHolderDetails(tst)
	if err == nil {
		t.Errorf("etAccountHolderDetailsFailure does not pass. Should fail. Looking for %v, got %v", "Not all field values present", nil)
	}
}

func TestSetAccountHolderDetails(t *testing.T) {
	tst := []string{"", "", "", "John", "Doe", "01011900", "010119001234123", "111", "222", "user@domain.com", "address 1", "address 2", "address 3", "2000", "cheque"}
	accountHolderDetails, err := setAccountHolderDetails(tst)

	if err != nil {
		t.Errorf("SetAccountHolderDetails does not pass.  Looking for %v, got %v", nil, err)
	}

	if reflect.TypeOf(accountHolderDetails).String() != "accounts.AccountHolderDetails" {
		t.Errorf("SetAccountHolderDetails does not pass. TYPE. Looking for %v, got %v", "accounts.AccountHolderDetails", reflect.TypeOf(accountHolderDetails).String())
	}

	if accountHolderDetails.GivenName != "John" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "John", accountHolderDetails.GivenName)
	}

	if accountHolderDetails.FamilyName != "Doe" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "Doe", accountHolderDetails.FamilyName)
	}

	if accountHolderDetails.DateOfBirth != "01011900" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "01011900", accountHolderDetails.DateOfBirth)
	}

	if accountHolderDetails.IdentificationNumber != "010119001234123" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "010119001234123", accountHolderDetails.IdentificationNumber)
	}

	if accountHolderDetails.ContactNumber1 != "111" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "111", accountHolderDetails.ContactNumber1)
	}

	if accountHolderDetails.ContactNumber2 != "222" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "222", accountHolderDetails.ContactNumber2)
	}

	if accountHolderDetails.EmailAddress != "user@domain.com" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "user@domain.com", accountHolderDetails.EmailAddress)
	}

	if accountHolderDetails.AddressLine1 != "address 1" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "address 1", accountHolderDetails.AddressLine1)
	}

	if accountHolderDetails.AddressLine2 != "address 2" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "address 2", accountHolderDetails.AddressLine2)
	}

	if accountHolderDetails.AddressLine3 != "address 3" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "address 3", accountHolderDetails.AddressLine3)
	}

	if accountHolderDetails.PostalCode != "2000" {
		t.Errorf("SetAccountHolderDetails does not pass. DETAILS. Looking for %v, got %v", "2000", accountHolderDetails.PostalCode)
	}
}

func BenchmarkSetAccountHolderDetails(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tst := []string{"", "", "", "John", "Doe", "01011900", "010119001234123", "111", "222", "user@domain.com", "address 1", "address 2", "address 3", "2000", "cheque"}
		_, _ = setAccountHolderDetails(tst)
	}
}

func TestSetAccountHolderDetailsSuccessAccountType(t *testing.T) {
	accountTypes := []string{"savings", "cheque", "merchant", "money-market", "cd", "ira", "rcp", "credit", "mortgage", "loan"}
	for _, v := range accountTypes {
		tst := []string{"", "", "", "John", "Doe", "01011900", "010119001234123", "111", "222", "user@domain.com", "address 1", "address 2", "address 3", "2000", v}
		_, err := setAccountHolderDetails(tst)

		if err != nil {
			t.Errorf("SetAccountHolderDetails does not pass.  Looking for %v, got %v", nil, err)
		}
	}
}

func TestSetAccountHolderDetailsFailureAccountType(t *testing.T) {
	tst := []string{"", "", "", "John", "Doe", "01011900", "010119001234123", "111", "222", "user@domain.com", "address 1", "address 2", "address 3", "2000", "not-valid-account-type"}
	_, err := setAccountHolderDetails(tst)

	if err != nil {
		t.Errorf("SetAccountHolderDetails does not pass.  Looking for %v, got %v", nil, err)
	}

}

/* @TODO
None of the above tests run against the functionality
Add functional tests like the one below, currently throwing nil pointer exception
func TestAddAccountPushTokenFunctional(t *testing.T) {
	//token~acmt~1003~token~platform
	tst := []string{"", "", "", "test-push-token", "other"}
	err := addAccountPushToken(tst)

	if err != nil {
		t.Errorf("AddAccountPushToken does not pass. Add token. Looking for %v, got %v", nil, err)
	}

	err = removeAccountPushToken(tst)
	if err != nil {
		t.Errorf("AddAccountPushToken does not pass. Remove token. Looking for %v, got %v", nil, err)
	}

}
*/
