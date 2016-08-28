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

	err = doCreateAccount(sqlTime, accountDetails, accountHolderDetails)
	if err != nil {
		return errors.New("accounts.createAccount: " + err.Error())
	}

	err = doCreateAccountUser(sqlTime, accountHolderDetails, accountDetails)
	if err != nil {
		return errors.New("accounts.createAccount: " + err.Error())
	}

	err = doCreateAccountUserAccount(sqlTime, accountHolderDetails, accountDetails)
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

	err = doDeleteAccountUser(accountHolderDetails)
	if err != nil {
		return errors.New("accounts.deleteAccount: " + err.Error())
	}

	err = doDeleteAccountUserAccounts(accountHolderDetails)
	if err != nil {
		return errors.New("accounts.deleteAccount: " + err.Error())
	}

	return
}

func doCreateAccount(sqlTime int32, accountDetails *AccountDetails, accountHolderDetails *AccountHolderDetails) (err error) {
	// Create account
	insertStatement := "INSERT INTO accounts (`accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance`, `type`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}

	// Prepare statement for inserting data
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Generate account number
	newUuid := uuid.NewV4()
	accountDetails.AccountNumber = newUuid.String()

	_, err = stmtIns.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName, accountDetails.AccountBalance, accountDetails.Overdraft, accountDetails.AvailableBalance, accountDetails.Type, sqlTime)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}

	// We insert a record into account user accounts
	insertStatement = "INSERT INTO accounts_users_accounts (`accountHolderIdentificationNumber`, `accountNumber`, `bankNumber`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?)"
	stmtIns, err = Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}

	// Prepare statement for inserting data
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(accountHolderDetails.IdentificationNumber, accountDetails.AccountNumber, accountDetails.BankNumber, sqlTime)
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

func doCreateAccountUser(sqlTime int32, accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	// Check if the user already exists
	account, err := getAccountUser(accountHolderDetails.IdentificationNumber)
	if err != nil {
		return errors.New("accounts.doCreateAccountUser: " + err.Error())
	}

	// If account is not empty it exists, return without creating a new one
	if account != (AccountHolderDetails{}) {
		return
	}

	// Create account meta
	insertStatement := "INSERT INTO accounts_users (`accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccountUser: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode, sqlTime)

	if err != nil {
		return errors.New("accounts.doCreateAccountUser: " + err.Error())
	}

	return
}

func doCreateAccountUserAccount(sqlTime int32, accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	insertStatement := "INSERT INTO accounts_users_accounts (`accountHolderIdentificationNumber`, `accountNumber`, `bankNumber`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccountUserAccount: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtIns.Exec(accountHolderDetails.IdentificationNumber, accountDetails.AccountNumber, accountDetails.BankNumber, sqlTime)

	if err != nil {
		return errors.New("accounts.doCreateAccountUserAccount: " + err.Error())
	}

	return
}

func doDeleteAccountUser(accountHolderDetails *AccountHolderDetails) (err error) {
	// Create account meta
	deleteStatement := "DELETE FROM accounts_users WHERE `accountHolderGivenName` = ? AND `accountHolderFamilyName` = ? AND `accountHolderDateOfBirth` = ? AND `accountHolderIdentificationNumber` = ? AND `accountHolderContactNumber1` = ? AND `accountHolderContactNumber2` = ? AND `accountHolderEmailAddress` = ? AND `accountHolderAddressLine1` = ? AND `accountHolderAddressLine2` = ? AND `accountHolderAddressLine3` = ? AND `accountHolderPostalCode` = ? "
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccountMeta: " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode)

	if err != nil {
		return errors.New("accounts.doDeleteAccountUser: " + err.Error())
	}

	return
}

func doDeleteAccountUserAccounts(accountHolderDetails *AccountHolderDetails) (err error) {
	// Create account meta
	deleteStatement := "DELETE FROM accounts_users_accounts WHERE `accountHolderIdentificationNumber` = ?"
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccountMeta: " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(accountHolderDetails.IdentificationNumber)

	if err != nil {
		return errors.New("accounts.doDeleteAccountUserAccount: " + err.Error())
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
		// There cannot be more than one account with the same accountNumber
		return AccountDetails{}, errors.New("accounts.getAccountDetails: More than one account found")
	}

	return
}

func getAccountUser(id string) (accountDetails AccountHolderDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode` FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ?", id)
	if err != nil {
		return AccountHolderDetails{}, errors.New("accounts.getAccountUser: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.DateOfBirth, &accountDetails.IdentificationNumber, &accountDetails.ContactNumber1, &accountDetails.ContactNumber2, &accountDetails.EmailAddress, &accountDetails.AddressLine1, &accountDetails.AddressLine2,
			&accountDetails.AddressLine3, &accountDetails.PostalCode); err != nil {
			//@TODO Throw error
			break
		}
		count++
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

func getUserAccountsDetail(userID string) (accounts []AccountDetails, err error) {
	rows, err := Config.Db.Query(
		"SELECT a.accountNumber, a.bankNumber, a.accountHolderName, a.accountBalance, a.overdraft, a.availableBalance "+
			"FROM accounts a "+
			"LEFT JOIN accounts_users_accounts au "+
			"ON au.accountNumber = a.accountNumber "+
			"AND au.bankNumber = a.bankNumber "+
			"WHERE au.accountHolderIdentificationNumber = ?", userID)
	if err != nil {
		return nil, errors.New("accounts.getUserAccountsDetail: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var account AccountDetails
		if err := rows.Scan(&account.AccountNumber, &account.BankNumber, &account.AccountHolderName, &account.AccountBalance, &account.Overdraft, &account.AvailableBalance); err != nil {
			break
		}

		accounts = append(accounts, account)
		count++
	}

	return
}

func getAllAccountNumbersByID(userID string) (accountIDs []string, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_users_accounts` WHERE `accountHolderIdentificationNumber` = ?", userID)
	if err != nil {
		return nil, errors.New("accounts.getAllAccountNumbersByID: " + err.Error())
	}
	defer rows.Close()

	count := 0
	// Return an array
	for rows.Next() {
		var accountID string
		if err := rows.Scan(&accountID); err != nil {
			break
		}
		count++
		accountIDs = append(accountIDs, accountID)
	}

	if count == 0 {
		return nil, errors.New("accounts.getAllAccountNumbersByID: Account not found")
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
	rows, err := Config.Db.Query("SELECT `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderEmailAddress` FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ? OR `accountHolderGivenName` = ? OR `accountHolderFamilyName` = ? OR  `accountHolderEmailAddress` = ? LIMIT 10", searchString, searchString, searchString, searchString)
	if err != nil {
		return []AccountHolderDetails{}, errors.New("accounts.getAccountMeta: " + err.Error())
	}
	defer rows.Close()

	allAccountDetails = []AccountHolderDetails{}
	count := 0
	for rows.Next() {
		accountDetails := AccountHolderDetails{}
		if err := rows.Scan(&accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.EmailAddress); err != nil {
			//@TODO Throw error
			break
		}
		allAccountDetails = append(allAccountDetails, accountDetails)
		count++
	}

	return
}

func getAccountByHolderDetails(ID string, givenName string, familyName string, email string) (accountIDs []string, err error) {
	// First we make sure that the ID number matches up to all the other information
	rows, err := Config.Db.Query("SELECT * FROM `accounts_users` WHERE `accountHolderIdentificationNumber` = ? AND `accountHolderGivenName` = ? AND `accountHolderFamilyName` = ? AND `accountHolderEmailAddress` = ?", ID, givenName, familyName, email)
	if err != nil {
		return nil, errors.New("accounts.getAccountByHolderDetails: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	if count > 0 {
		rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_users_accounts` WHERE `accountHolderIdentificationNumber` = ?", ID)
		if err != nil {
			return nil, errors.New("accounts.getAccountByHolderDetails: " + err.Error())
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var accountID string
			if err := rows.Scan(&accountID); err != nil {
				//@TODO Throw error
				break
			}
			accountIDs = append(accountIDs, accountID)
			count++
		}
	}

	return
}
