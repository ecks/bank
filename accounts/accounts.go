package accounts

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bvnk/bank/appauth"
	"github.com/shopspring/decimal"
)

/*
Accounts package to deal with all account related queries.

@TODO Implement the ISO20022 standard
http://www.iso20022.org/full_catalogue.page - acmt

@TODO Consider moving checkBalances, updateBalance to here

Accounts (acmt) transactions are as follows:
1  - AccountOpeningInstructionV05
2  - AccountDetailsConfirmationV05
3  - AccountModificationInstructionV05
5  - RequestForAccountManagementStatusReportV03
6  - AccountManagementStatusReportV04
7  - AccountOpeningRequestV02
8  - AccountOpeningAmendmentRequestV02
9  - AccountOpeningAdditionalInformationRequestV02
10 - AccountRequestAcknowledgementV02
11 - AccountRequestRejectionV02
12 - AccountAdditionalInformationRequestV02
13 - AccountReportRequestV02
14 - AccountReportV02
15 - AccountExcludedMandateMaintenanceRequestV02
16 - AccountExcludedMandateMaintenanceAmendmentRequestV02
17 - AccountMandateMaintenanceRequestV02
18 - AccountMandateMaintenanceAmendmentRequestV02
19 - AccountClosingRequestV02
20 - AccountClosingAmendmentRequestV02
21 - AccountClosingAdditionalInformationRequestV02
22 - IdentificationModificationAdviceV02
23 - IdentificationVerificationRequestV02
24 - IdentificationVerificationReportV02

### Custom functionality
1000 - ListAllAccounts (Deprecated)
1001 - ListSingleAccount
1002 - CheckAccountByID
1003 - AddAccountPushToken
1004 - RemoveAccountPushToken
1005 - SearchForAccount
1006 - RetrieveAccount

## Merchant accounts
1100 - MerchantAccountCreate
1101 - MerchantAccountUpdate
1102 - MerchantAccountView
1103 - MerchantAccountDelete
1104 - MerchantAccountSearch

*/

/* acmt~1~
   AccountHolderGivenName~
   AccountHolderFamilyName~
   AccountHolderDateOfBirth~
   AccountHolderIdentificationNumber~
   AccountHolderContactNumber1~
   AccountHolderContactNumber2~
   AccountHolderEmailAddress~
   AccountHolderAddressLine1~
   AccountHolderAddressLine2~
   AccountHolderAddressLine3~
   AccountHolderPostalCode
*/
type AccountHolder struct {
	AccountNumber string
	BankNumber    string
}

type AccountHolderDetails struct {
	GivenName            string
	FamilyName           string
	DateOfBirth          string
	IdentificationNumber string
	ContactNumber1       string
	ContactNumber2       string
	EmailAddress         string
	AddressLine1         string
	AddressLine2         string
	AddressLine3         string
	PostalCode           string
}

type AccountHolderAccounts struct {
	IdentificationNumber string
	AccountNumber        string
	BankNumber           string
}

type AccountDetails struct {
	AccountNumber     string
	BankNumber        string
	AccountHolderName string
	AccountBalance    decimal.Decimal
	Overdraft         decimal.Decimal
	AvailableBalance  decimal.Decimal
	Type              string
	Timestamp         int
}

type MerchantDetails struct {
	ID                   string
	Name                 string
	Description          string
	ContactGivenName     string
	ContactFamilyName    string
	AddressLine1         string
	AddressLine2         string
	AddressLine3         string
	Country              string
	PostalCode           string
	BusinessSector       string
	Website              string
	ContactPhone         string
	ContactFax           string
	ContactEmail         string
	Logo                 string
	IdentificationNumber string
	Timestamp            int
}

// Set up some defaults
const (
	BANK_NUMBER       = "a0299975-b8e2-4358-8f1a-911ee12dbaac"
	OPENING_BALANCE   = 100.
	OPENING_OVERDRAFT = 0.
)

func ProcessAccount(data []string) (result interface{}, err error) {
	if len(data) < 3 {
		return "", errors.New("accounts.ProcessAccount: Not enough fields, minimum 3")
	}

	acmtType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return "", errors.New("accounts.ProcessAccount: Could not get ACMT type")
	}

	// Switch on the acmt type
	switch acmtType {
	case 1, 7:
		/*
		   @TODO
		   The differences between AccountOpeningInstructionV05 and AccountOpeningRequestV02 will be explored in detail, for now we treat the same - open an account
		*/
		result, err = openAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1001:
		result, err = fetchUserAccounts(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1002:
		if len(data) < 4 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = fetchAllAccountsByID(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
		break
	case 1003:
		if len(data) < 5 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		err = addAccountPushToken(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	case 1004:
		if len(data) < 5 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		err = removeAccountPushToken(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	case 1005:
		// acmt~1005~token~search-information
		if len(data) < 4 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = searchAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	case 1006:
		if len(data) < 7 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = retrieveAccount(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	// Merchant account create
	case 1100:
		if len(data) < 19 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = merchantAccountCreate(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	// Merchant account update
	case 1101:
		if len(data) < 19 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = merchantAccountUpdate(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	// Merchant account view
	case 1102:
		if len(data) < 4 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = merchantAccountView(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	// Merchant account delete
	case 1103:
		if len(data) < 4 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = merchantAccountDelete(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	// Merchant account search
	case 1104:
		if len(data) < 4 {
			err = errors.New("accounts.ProcessAccount: Not all fields present")
			return
		}
		result, err = merchantAccountSearch(data)
		if err != nil {
			return "", errors.New("accounts.ProcessAccount: " + err.Error())
		}
	default:
		err = errors.New("accounts.ProcessAccount: ACMT transaction code invalid")
		break
	}

	return
}

func openAccount(data []string) (result interface{}, err error) {
	// Validate string against required info/length
	if len(data) < 14 {
		err = errors.New("accounts.openAccount: Not all fields present")
		return
	}

	// Test: acmt~1~Kyle~Redelinghuys~19000101~190001011234098~1112223456~~email@domain.com~Physical Address 1~~~1000
	// @FIXME: Remove new line from data
//	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Create account
	accountHolderObject, err := setAccountDetails(data)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}
	accountHolderDetailsObject, err := setAccountHolderDetails(data)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}
	err = createAccount(&accountHolderObject, &accountHolderDetailsObject)
	if err != nil {
		return "", errors.New("accounts.openAccount: " + err.Error())
	}

	result = accountHolderObject.AccountNumber
	return
}

func closeAccount(data []string) (result interface{}, err error) {
	// Validate string against required info/length
	if len(data) < 15 {
		err = errors.New("accounts.closeAccount: Not all fields present")
		return
	}

	// Check if account already exists, check on ID number
	accountHolder, _ := getAccountUser(data[6])
	if accountHolder == (AccountHolderDetails{}) {
		return "", errors.New("accounts.closeAccount: Account does not exist.")
	}

	// @FIXME: Remove new line from data
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	// Delete account
	accountHolderObject, err := setAccountDetails(data)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}
	accountHolderDetailsObject, err := setAccountHolderDetails(data)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}
	err = deleteAccount(&accountHolderObject, &accountHolderDetailsObject)
	if err != nil {
		return "", errors.New("accounts.closeAccount: " + err.Error())
	}

	return
}

func setAccountDetails(data []string) (accountDetails AccountDetails, err error) {
	if len(data) < 14 {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Not all fields required present")
	}

	if data[4] == "" {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Given name cannot be empty")
	}
	accountDetails.BankNumber = BANK_NUMBER
	accountDetails.AccountHolderName = data[4] + "," + data[3] // Family Name, Given Name
	accountDetails.AccountBalance = decimal.NewFromFloat(OPENING_BALANCE)
	accountDetails.Overdraft = decimal.NewFromFloat(OPENING_OVERDRAFT)
	accountDetails.AvailableBalance = decimal.NewFromFloat(OPENING_BALANCE + OPENING_OVERDRAFT)
	// Get account type
	accountType := data[14]
	switch accountType {
	case "":
		accountType = "cheque" // Default to chequing account
		break
	case "savings", "cheque", "merchant", "money-market", "cd", "ira", "rcp", "credit", "mortgage", "loan":
		// Valid
		break
	default:
		return AccountDetails{}, errors.New("accounts.setAccountDetails: Account type not valid, must be one of savings, cheque, merchant, money-market, cd, ira, rcp, credit, mortgage, loan")
		break
	}
	accountDetails.Type = accountType

	return
}

func setAccountHolderDetails(data []string) (accountHolderDetails AccountHolderDetails, err error) {
	if len(data) < 12 {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Not all field values present")
	}
	//@TODO: Test date parsing in format ddmmyyyy
	if data[4] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Family name cannot be empty")
	}
	if data[3] == "" {
		return AccountHolderDetails{}, errors.New("accounts.setAccountHolderDetails: Given name cannot be empty")
	}

	// @TODO Integrity checks
	accountHolderDetails.GivenName = data[3]
	accountHolderDetails.FamilyName = data[4]
	accountHolderDetails.DateOfBirth = data[5]
	accountHolderDetails.IdentificationNumber = data[6]
	accountHolderDetails.ContactNumber1 = data[7]
	accountHolderDetails.ContactNumber2 = data[8]
	accountHolderDetails.EmailAddress = data[9]
	accountHolderDetails.AddressLine1 = data[10]
	accountHolderDetails.AddressLine2 = data[11]
	accountHolderDetails.AddressLine3 = data[12]
	accountHolderDetails.PostalCode = data[13]

	return
}

func fetchUserAccounts(data []string) (account interface{}, err error) {
	// Fetch user account. Must be user logged in
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.fetchSingleAccount: " + err.Error())
	}
	account, err = getUserAccountsDetail(tokenUser)
	if err != nil {
		return "", errors.New("accounts.fetchSingleAccount: " + err.Error())
	}

	return
}

func fetchAllAccountsByID(data []string) (userAccountNumber []string, err error) {
	// Format: token~acmt~1002~USERID
	userID := data[3]
	if userID == "" {
		return nil, errors.New("accounts.fetchSingleAccountByID: User ID not present")
	}

	userAccountNumber, err = getAllAccountNumbersByID(userID)
	if err != nil {
		return nil, errors.New("accounts.fetchSingleAccountByID: " + err.Error())
	}

	return
}

func addAccountPushToken(data []string) (err error) {
	//Format: token~acmt~1003~token~platform
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return errors.New("accounts.addAccountPushToken: " + err.Error())
	}

	// Check platform is correctly set
	// @FIXME This feels very heavy handed for "check if value in array"
	platform := data[4]
	platformPass := false
	allowedPlatforms := []string{"ios", "windows", "android", "blackberry", "other"}
	for _, v := range allowedPlatforms {
		if strings.Compare(v, platform) == 0 {
			platformPass = true
		}
	}

	if !platformPass {
		return errors.New("accounts.addAccountPushToken: Platform invalid")
	}

	err = doAddAccountPushToken(tokenUser, data[3], platform)
	if err != nil {
		return err
	}
	return nil
}

func removeAccountPushToken(data []string) (err error) {
	//Format: token~acmt~1003~token~platform
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return errors.New("accounts.addAccountPushToken: " + err.Error())
	}

	// Check platform is correctly set
	// @FIXME This feels very heavy handed for "check if value in array"
	platform := data[4]
	platformPass := false
	allowedPlatforms := []string{"ios", "windows", "android", "blackberry", "other"}
	for _, v := range allowedPlatforms {
		if strings.Compare(v, platform) == 0 {
			platformPass = true
		}
	}

	if !platformPass {
		return errors.New("accounts.addAccountPushToken: Platform invalid")
	}
	err = doDeleteAccountPushToken(tokenUser, data[3], platform)
	if err != nil {
		return err
	}
	return nil
}

// searchAccountData takes a search term and searches on id, first name and last name
// Results are limited to 10
func searchAccount(data []string) (accounts interface{}, err error) {
	//Format: acmt~1005~token~search-string
	_, err = appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.searchAccount: " + err.Error())
	}

	searchString := data[3]
	if len(searchString) < 4 {
		return "", errors.New("accounts.merchantAccountSearch: Please provide a search term of at least 4 characters.")
	}

	accounts, err = getAccountFromSearchData(searchString)
	if err != nil {
		return "", errors.New("accounts.searchAccount: Searching for account error. " + err.Error())
	}

	return
}

func retrieveAccount(data []string) (accountIDs []string, err error) {
	id := data[3]
	givenName := data[4]
	familyName := data[5]
	email := data[6]

	if id == "" || givenName == "" || familyName == "" || email == "" {
		return nil, errors.New("accounts.retrieveAccount: Not all fields present")
	}

	accountIDs, err = getAccountByHolderDetails(id, givenName, familyName, email)
	if err != nil {
		return nil, errors.New("accounts.retrieveAccount: Could not retrieve account. " + err.Error())
	}

	return
}

func CheckUserAccountValidFromToken(userID string, accountNumber string) (err error) {
	// Get list of accounts from userID
	userAccountNumbers, err := getAllAccountNumbersByID(userID)
	if err != nil {
		return errors.New("accounts.CheckUserAccountValidFromToken: " + err.Error())
	}

	senderValid := false
	for _, v := range userAccountNumbers {
		if strings.Compare(v, accountNumber) == 0 {
			senderValid = true
		}
	}

	if !senderValid {
		return errors.New("accounts.CheckUserAccountValidFromToken: Sender invalid")
	}
	return
}

func merchantAccountCreate(data []string) (result interface{}, err error) {
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: " + err.Error())
	}
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	accountHolder, err := getAccountUser(tokenUser)
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: Could not retrieve account user from token user.")
	}
	// Create account
	merchantObject, accountDetails, err := setMerchantDetails(data, accountHolder.IdentificationNumber, "create")
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: " + err.Error())
	}

	err = createMerchantAccount(&merchantObject, &accountDetails, &accountHolder)
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: " + err.Error())
	}

	result = merchantObject.ID
	return
}

func merchantAccountUpdate(data []string) (result interface{}, err error) {
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: " + err.Error())
	}
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	accountHolder, err := getAccountUser(tokenUser)
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: Could not retrieve account user from token user.")
	}

	// Set merchant details
	merchantObject, _, err := setMerchantDetails(data, accountHolder.IdentificationNumber, "update")
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: " + err.Error())
	}

	// Get list of accounts linked to merchant
	accountIDs, err := getAllMerchantAccountNumbersByMerchantID(merchantObject.ID)
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: Could not get accountIDs associated to Merchant ID: " + err.Error())
	}

	// Make sure the user has access to all merchant accounts
	for _, aID := range accountIDs {
		err = CheckUserAccountValidFromToken(accountHolder.IdentificationNumber, aID)
		if err != nil {
			return "", errors.New("accounts.merchantAccountUpdate: Account holder not valid for given Merchant ID. " + err.Error())
		}
	}

	currentMerchantObject, err := getMerchantFromMerchantID(merchantObject.ID)
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: Could not retrieve current merchant using given merchant ID. " + err.Error())
	}

	// Check which values are updated and set in new object
	compareUpdateMerchantWithCurrentMerchant(&merchantObject, &currentMerchantObject)

	err = updateMerchant(&merchantObject)
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: Could not update merchant. " + err.Error())
	}
	return
}

func merchantAccountDelete(data []string) (result interface{}, err error) {
	// Currently this deletes an account entirely
	// Deleting is needed but backups of all data should be done
	// The account could be marked as deleted when called through the API
	// @TODO: Find best solution here. One option: log all data in separate place
	// (perhaps in a file) and remove the record
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: " + err.Error())
	}
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	accountHolderDetails, err := getAccountUser(tokenUser)
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: Could not retrieve account user from token user.")
	}

	accountDetails, err := getAccountDetails(data[5])
	if err != nil {
		return "", errors.New("accounts.merchantAccountCreate: Could not retrieve account details from account number provided.")
	}

	merchantObject, err := getMerchantFromMerchantID(data[4])
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: Could not retrieve current merchant using given merchant ID. " + err.Error())
	}

	// Get list of accounts linked to merchant
	accountIDs, err := getAllMerchantAccountNumbersByMerchantID(merchantObject.ID)
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: Could not get accountIDs associated to Merchant ID: " + err.Error())
	}

	// Make sure the user has access to all merchant accounts
	for _, aID := range accountIDs {
		err = CheckUserAccountValidFromToken(accountHolderDetails.IdentificationNumber, aID)
		if err != nil {
			return "", errors.New("accounts.merchantAccountUpdate: Account holder not valid for given Merchant ID. " + err.Error())
		}
	}

	err = deleteMerchantAccount(&merchantObject, &accountDetails, &accountHolderDetails)
	if err != nil {
		return "", errors.New("accounts.merchantAccountUpdate: Could not update merchant. " + err.Error())
	}

	return
}

func merchantAccountView(data []string) (merchant MerchantDetails, err error) {
	_, err = appauth.GetUserFromToken(data[0])
	if err != nil {
		return MerchantDetails{}, errors.New("accounts.merchantAccountUpdate: " + err.Error())
	}
	data[len(data)-1] = strings.Replace(data[len(data)-1], "\n", "", -1)

	merchant, err = getMerchantFromMerchantID(data[3])
	if err != nil {
		return MerchantDetails{}, errors.New("accounts.merchantAccountView: Could not retrieve merchant from given Merchant ID.")
	}
	return
}

func merchantAccountSearch(data []string) (result interface{}, err error) {
	_, err = appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("accounts.merchantAccountSearch: " + err.Error())
	}

	searchString := data[3]
	if len(searchString) < 2 {
		return "", errors.New("accounts.merchantAccountSearch: Please provide a search term of at least 2 characters.")
	}

	result, err = getMerchantAccountFromSearchData(searchString)
	if err != nil {
		return "", errors.New("accounts.merchantAccountSearch: Searching for account error. " + err.Error())
	}

	return
}

func setMerchantDetails(data []string, identificationNumber string, setType string) (merchantDetails MerchantDetails, accountDetails AccountDetails, err error) {
	if len(data) < 19 {
		return MerchantDetails{}, AccountDetails{}, errors.New("accounts.setMerchantDetails: Not all field values present")
	}

	if data[3] == "" {
		return MerchantDetails{}, AccountDetails{}, errors.New("accounts.setMerchantDetails: Merchant name cannot be empty")
	}
	if data[4] == "" {
		return MerchantDetails{}, AccountDetails{}, errors.New("accounts.setMerchantDetails: Merchant description cannot be empty")
	}
	if data[12] == "" {
		return MerchantDetails{}, AccountDetails{}, errors.New("accounts.setMerchantDetails: Merchant business sector cannot be empty")
	}
	if data[16] == "" {
		return MerchantDetails{}, AccountDetails{}, errors.New("accounts.setMerchantDetails: Merchant contact email cannot be empty")
	}

	merchantDetails.Name = data[3]
	merchantDetails.Description = data[4]
	merchantDetails.ContactGivenName = data[5]
	merchantDetails.ContactFamilyName = data[6]
	merchantDetails.AddressLine1 = data[7]
	merchantDetails.AddressLine2 = data[8]
	merchantDetails.AddressLine3 = data[9]
	merchantDetails.Country = data[10]
	merchantDetails.PostalCode = data[11]
	merchantDetails.BusinessSector = data[12]
	merchantDetails.Website = data[13]
	merchantDetails.ContactPhone = data[14]
	merchantDetails.ContactFax = data[15]
	merchantDetails.ContactEmail = data[16]
	// @FIXME We leave logo out for now, not sure how to parse pictures
	merchantDetails.IdentificationNumber = identificationNumber

	accountDetails.BankNumber = BANK_NUMBER
	accountDetails.AccountHolderName = data[3] // Business Name
	accountDetails.AccountBalance = decimal.NewFromFloat(OPENING_BALANCE)
	accountDetails.Overdraft = decimal.NewFromFloat(OPENING_OVERDRAFT)
	accountDetails.AvailableBalance = decimal.NewFromFloat(OPENING_BALANCE + OPENING_OVERDRAFT)

	if setType == "create" {
		accountType := data[18]
		switch accountType {
		case "":
			accountType = "merchant" // Default to merchant account
			break
		case "merchant", "credit", "mortgage", "loan":
			// Valid
			break
		default:
			return MerchantDetails{}, AccountDetails{}, errors.New("accounts.setMerchantDetails: Account type not valid for a merchant, must be one of merchant, credit, mortgage, loan")
			break
		}
		accountDetails.Type = accountType
	} else if setType == "update" {
		if data[18] != "" {
			merchantDetails.ID = data[18]
		}
	}

	return
}

func compareUpdateMerchantWithCurrentMerchant(merchantUpdate *MerchantDetails, merchantCurrent *MerchantDetails) {
	if merchantUpdate.Name == "" {
		merchantUpdate.Name = merchantCurrent.Name
	}
	if merchantUpdate.Description == "" {
		merchantUpdate.Description = merchantCurrent.Description
	}
	if merchantUpdate.ContactGivenName == "" {
		merchantUpdate.ContactGivenName = merchantCurrent.ContactGivenName
	}
	if merchantUpdate.ContactFamilyName == "" {
		merchantUpdate.ContactFamilyName = merchantCurrent.ContactFamilyName
	}
	if merchantUpdate.AddressLine1 == "" {
		merchantUpdate.AddressLine1 = merchantCurrent.AddressLine1
	}
	if merchantUpdate.AddressLine2 == "" {
		merchantUpdate.AddressLine2 = merchantCurrent.AddressLine2
	}
	if merchantUpdate.AddressLine3 == "" {
		merchantUpdate.AddressLine3 = merchantCurrent.AddressLine3
	}
	if merchantUpdate.Country == "" {
		merchantUpdate.Country = merchantCurrent.Country
	}
	if merchantUpdate.PostalCode == "" {
		merchantUpdate.PostalCode = merchantCurrent.PostalCode
	}
	if merchantUpdate.BusinessSector == "" {
		merchantUpdate.BusinessSector = merchantCurrent.BusinessSector
	}
	if merchantUpdate.Website == "" {
		merchantUpdate.Website = merchantCurrent.Website
	}
	if merchantUpdate.ContactPhone == "" {
		merchantUpdate.ContactPhone = merchantCurrent.ContactPhone
	}
	if merchantUpdate.ContactFax == "" {
		merchantUpdate.ContactFax = merchantCurrent.ContactFax
	}
	if merchantUpdate.ContactEmail == "" {
		merchantUpdate.ContactEmail = merchantCurrent.ContactEmail
	}
}
