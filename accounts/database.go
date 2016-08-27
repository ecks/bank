package accounts

/*
@TODO Fix DB repetition
*/

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bvnk/bank/configuration"
	"github.com/satori/go.uuid"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func loadDatabase() (db *sql.DB, err error) {
	// Test connection with ping
	err = Config.Db.Ping()
	if err != nil {
		return
	}

	return
}

func createAccount(accountDetails *AccountDetails, accountHolderDetails *AccountHolderDetails) (err error) {
	// Convert variables
	t := time.Now()
	sqlTime := int32(t.Unix())

	err = doCreateAccount(sqlTime, accountDetails)
	if err != nil {
		return errors.New("accounts.createAccount: " + err.Error())
	}

	err = doCreateAccountMeta(sqlTime, accountHolderDetails, accountDetails)
	if err != nil {
		return errors.New("accounts.createAccount: " + err.Error())
	}

	return
}

func deleteAccount(accountDetails *AccountDetails, accountHolderDetails *AccountHolderDetails) (err error) {
	err = doDeleteAccount(accountDetails)
	if err != nil {
		return errors.New("accounts.deleteAccount: " + err.Error())
	}

	err = doDeleteAccountMeta(accountHolderDetails, accountDetails)
	if err != nil {
		return errors.New("accounts.deleteAccount: " + err.Error())
	}

	return
}

func doCreateAccount(sqlTime int32, accountDetails *AccountDetails) (err error) {
	// Create account
	insertStatement := "INSERT INTO accounts (`accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}

	// Prepare statement for inserting data
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Generate account number
	newUuid := uuid.NewV4()
	accountDetails.AccountNumber = newUuid.String()

	_, err = stmtIns.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName, accountDetails.AccountBalance, accountDetails.Overdraft, accountDetails.AvailableBalance, sqlTime)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}
	return
}

func doDeleteAccount(accountDetails *AccountDetails) (err error) {
	deleteStatement := "DELETE FROM accounts WHERE `accountNumber` = ? AND `bankNumber` = ? AND `accountHolderName` = ? "
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccount: " + err.Error())
	}

	// Prepare statement for inserting data
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName)
	if err != nil {
		return errors.New("accounts.doDeleteAccount: " + err.Error())
	}
	// Can use db.RowsAffected() to check
	return
}

func doCreateAccountMeta(sqlTime int32, accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	// Create account meta
	insertStatement := "INSERT INTO accounts_users (`accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccountMeta: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtIns.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode, sqlTime)

	if err != nil {
		return errors.New("accounts.doCreateAccountMeta: " + err.Error())
	}

	return
}

func doDeleteAccountMeta(accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	// Create account meta
	deleteStatement := "DELETE FROM accounts_users WHERE `accountNumber` = ? AND `bankNumber` = ? AND `accountHolderGivenName` = ? AND `accountHolderFamilyName` = ? AND `accountHolderDateOfBirth` = ? AND `accountHolderIdentificationNumber` = ? AND `accountHolderContactNumber1` = ? AND `accountHolderContactNumber2` = ? AND `accountHolderEmailAddress` = ? AND `accountHolderAddressLine1` = ? AND `accountHolderAddressLine2` = ? AND `accountHolderAddressLine3` = ? AND `accountHolderPostalCode` = ? "
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccountMeta: " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtDel.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode)

	if err != nil {
		return errors.New("accounts.doDeleteAccountMeta: " + err.Error())
	}

	return
}

func getAccountDetails(id string) (accountDetails AccountDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", id)
	if err != nil {
		return AccountDetails{}, errors.New("accounts.getAccountDetails: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.AccountHolderName, &accountDetails.AccountBalance, &accountDetails.Overdraft, &accountDetails.AvailableBalance)
		if err != nil {
			break
		}
		count++
	}

	if count == 0 {
		return AccountDetails{}, errors.New("accounts.getAccountDetails: Account not found")
	}

	if count > 1 {
		//@TODO: Allow user to have multiple accounts
		return AccountDetails{}, errors.New("accounts.getAccountDetails: More than one account found")
	}

	return
}

func getAccountMeta(id string) (accountDetails AccountHolderDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode` FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ?", id)
	if err != nil {
		return AccountHolderDetails{}, errors.New("accounts.getAccountMeta: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.DateOfBirth, &accountDetails.IdentificationNumber, &accountDetails.ContactNumber1, &accountDetails.ContactNumber2, &accountDetails.EmailAddress, &accountDetails.AddressLine1, &accountDetails.AddressLine2,
			&accountDetails.AddressLine3, &accountDetails.PostalCode); err != nil {
			//@TODO Throw error
			break
		}
		count++
	}

	if count == 0 {
		return AccountHolderDetails{}, errors.New("accounts.getAccountMeta: Account not found")
	}

	if count > 1 {
		//@TODO: Allow user to have multiple accounts
		return AccountHolderDetails{}, errors.New("accounts.getAccountMeta: More than one account found")
	}

	return
}

func getAllAccountDetails() (allAccounts []AccountDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName` FROM `accounts`")
	if err != nil {
		return []AccountDetails{}, errors.New("accounts.getAllAccountDetails: Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	allAccounts = make([]AccountDetails, 0)

	for rows.Next() {
		accountDetailsSingle := AccountDetails{}
		if err := rows.Scan(&accountDetailsSingle.AccountNumber, &accountDetailsSingle.BankNumber, &accountDetailsSingle.AccountHolderName); err != nil {
			break
		}

		allAccounts = append(allAccounts, accountDetailsSingle)
		count++
	}

	return
}

func getSingleAccountDetail(accountNumber string) (account AccountDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return AccountDetails{}, errors.New("accounts.getSingleAccountDetail: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&account.AccountNumber, &account.BankNumber, &account.AccountHolderName, &account.AccountBalance, &account.Overdraft, &account.AvailableBalance); err != nil {
			break
		}

		count++
	}

	return
}

func getSingleAccountNumberByID(userID string) (accountID string, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ?", userID)
	if err != nil {
		return "", errors.New("accounts.getSingleAccountNumberByID: " + err.Error())
	}
	defer rows.Close()

	count := 0
	// @TODO Right now this will return the latest account only, if there are two accounts
	for rows.Next() {
		if err := rows.Scan(&accountID); err != nil {
			break
		}
		count++
	}

	if count == 0 {
		return "", errors.New("accounts.getSingleAccountNumberByID: Account not found")
	}

	return
}

func doAddAccountPushToken(accountNumber string, pushToken string, platform string) (err error) {
	t := time.Now()
	sqlTime := int32(t.Unix())

	// Check if push token already exists for user
	pushTokenExists, err := fetchAccountPushToken(accountNumber, pushToken, platform)
	if err != nil {
		return errors.New("accounts.doAddAccountPushToken: " + err.Error())
	}

	if pushTokenExists == true {
		// Delete current push token
		err = doDeleteAccountPushToken(accountNumber, pushToken, platform)
		if err != nil {
			return errors.New("accounts.doAddAccountPushToken: Could not delete existing push token: " + err.Error())
		}
	}

	insertStatement := "INSERT INTO accounts_push_tokens (`accountNumber`, `token`, `platform`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doAddAccountPushToken: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(accountNumber, pushToken, platform, sqlTime)

	if err != nil {
		return errors.New("accounts.doAddAccountPushToken: " + err.Error())
	}

	return
}

func fetchAccountPushToken(accountNumber string, pushToken string, platform string) (pushTokenExists bool, err error) {
	rows, err := Config.Db.Query("SELECT * FROM `accounts_push_tokens` WHERE `accountNumber` = ? AND `token` = ? AND `platform` = ?", accountNumber, pushToken, platform)
	if err != nil {
		return false, errors.New("accounts.fetchAccountPushToken: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&pushToken); err != nil {
			pushTokenExists = true
			break
		}
	}

	return
}

func doDeleteAccountPushToken(accountNumber string, pushToken string, platform string) (err error) {
	deleteStatement := "DELETE FROM `accounts_push_tokens` WHERE `accountNumber` = ? AND `token` = ? AND `platform` = ?"
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccountPushToken: " + err.Error())
	}

	defer stmtDel.Close()

	_, err = stmtDel.Exec(accountNumber, pushToken, platform)
	if err != nil {
		return errors.New("accounts.doDeleteAccountPushToken: " + err.Error())
	}
	// Can use db.RowsAffected() to check
	return
}

func getAccountFromSearchData(searchString string) (allAccountDetails []AccountHolderDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderEmailAddress` FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ? OR `accountHolderGivenName` = ? OR `accountHolderFamilyName` = ? OR  `accountHolderEmailAddress` = ? LIMIT 10", searchString, searchString, searchString, searchString)
	if err != nil {
		return []AccountHolderDetails{}, errors.New("accounts.getAccountMeta: " + err.Error())
	}
	defer rows.Close()

	allAccountDetails = []AccountHolderDetails{}
	count := 0
	for rows.Next() {
		accountDetails := AccountHolderDetails{}
		if err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.EmailAddress); err != nil {
			//@TODO Throw error
			break
		}
		allAccountDetails = append(allAccountDetails, accountDetails)
		count++
	}

	return
}

func getAccountByHolderDetails(ID string, givenName string, familyName string, email string) (accountID string, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ? AND `accountHolderGivenName` = ? AND `accountHolderFamilyName` = ? AND `accountHolderEmailAddress` = ?", ID, givenName, familyName, email)
	if err != nil {
		return "", errors.New("accounts.getAccountByHolderDetails: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountID); err != nil {
			//@TODO Throw error
			break
		}
		count++
	}

	return
}
