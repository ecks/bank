package transactions

import (
	"errors"
	"time"

	"github.com/bvnk/bank/configuration"
	"github.com/shopspring/decimal"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func savePainTransaction(transaction PAINTrans) (err error) {
	// Prepare statement for inserting data
	// Construct geoText. These values are already cleared
	geoText := transaction.Geo.ToWKT()
	insertStatement := "INSERT INTO transactions (`transaction`, `type`, `senderAccountNumber`, `senderBankNumber`, `receiverAccountNumber`, `receiverBankNumber`, `transactionAmount`, `feeAmount`, `desc`, `timestamp`, `status`, `geo`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, GeomFromText(?))"

	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	t := time.Now()
	sqlTime := int32(t.Unix())
	transaction.Timestamp = sqlTime

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount.Mul(transaction.Fee)

	_, err = stmtIns.Exec("pain", transaction.PainType, transaction.Sender.AccountNumber, transaction.Sender.BankNumber, transaction.Receiver.AccountNumber, transaction.Receiver.BankNumber,
		transaction.Amount, feeAmount, transaction.Desc, transaction.Timestamp, transaction.Status, geoText)

	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}

	return
}

// This is for testing. Transactions should never be removed
func removePainTransaction(transaction PAINTrans) (err error) {
	// Prepare statement for inserting data
	delStatement := "DELETE FROM transactions WHERE `transaction` = ? AND `type` = ? AND `senderAccountNumber` = ? AND `senderBankNumber` = ? AND `receiverAccountNumber` = ? AND `receiverBankNumber` = ? AND `transactionAmount` = ? AND `feeAmount` = ? "
	stmtDel, err := Config.Db.Prepare(delStatement)
	if err != nil {
		return errors.New("payments.removePainTransaction: " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount.Mul(transaction.Fee)

	_, err = stmtDel.Exec("pain", transaction.PainType, transaction.Sender.AccountNumber, transaction.Sender.BankNumber, transaction.Receiver.AccountNumber, transaction.Receiver.BankNumber,
		transaction.Amount, feeAmount)

	if err != nil {
		return errors.New("payments.removePainTransaction: " + err.Error())
	}

	return
}

//func updateAccounts(sender AccountHolder, receiver AccountHolder, transactionAmount float64, transactionFee float64) {
func updateAccounts(transaction PAINTrans) (err error) {
	t := time.Now()
	sqlTime := int32(t.Unix())

	// The feePerc is a percentage, convert to amount
	feeAmount := transaction.Amount.Mul(transaction.Fee)

	switch transaction.PainType {
	// Payment
	case 1:
		err = processCreditInitiation(transaction, sqlTime, feeAmount)
		if err != nil {
			return errors.New("payments.updateAccounts: " + err.Error())
		}
		break
	// Deposit
	case 1000:
		err = processDepositInitiation(transaction, sqlTime, feeAmount)
		if err != nil {
			return errors.New("payments.updateAccounts: " + err.Error())
		}
		break
	}

	err = updateBankHoldingAccount(feeAmount, sqlTime)
	if err != nil {
		return errors.New("payments.updateAccounts: " + err.Error())
	}

	return

}

func updateBankHoldingAccount(feeAmount decimal.Decimal, sqlTime int32) (err error) {
	// Add fees to bank holding account
	// Only one row in this account for now - only holds single holding bank's balance
	updateBank := "UPDATE `bank_account` SET `balance` = (`balance` + ?), `timestamp` = ?"
	stmtUpdBank, err := Config.Db.Prepare(updateBank)
	if err != nil {
		return errors.New("payments.updateBankHoldingAccount: " + err.Error())
	}
	defer stmtUpdBank.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtUpdBank.Exec(feeAmount, sqlTime)

	if err != nil {
		return errors.New("payments.updateBankHoldingAccount: " + err.Error())
	}
	return
}

// @TODO Look at using accounts.getAccountDetails here
func checkBalance(account AccountHolder) (balance decimal.Decimal, err error) {
	rows, err := Config.Db.Query("SELECT `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", account.AccountNumber)
	if err != nil {
		return decimal.NewFromFloat(0.), errors.New("payments.checkBalance: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&balance); err != nil {
			return decimal.NewFromFloat(0.), errors.New("payments.checkBalance: Could not retrieve account details. " + err.Error())
		}
		count++
	}

	if count > 1 {
		return decimal.NewFromFloat(0.), errors.New("payments.checkBalance: More than one account found with uuid")
	}

	return
}

func processCreditInitiation(transaction PAINTrans, sqlTime int32, feeAmount decimal.Decimal) (err error) {
	// Only update if account local
	if transaction.Sender.BankNumber == "" {
		updateSenderStatement := "UPDATE accounts SET `accountBalance` = (`accountBalance` - ?), `availableBalance` = (`availableBalance` - ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdSender, err := Config.Db.Prepare(updateSenderStatement)
		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}
		defer stmtUpdSender.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdSender.Exec(transaction.Amount.Add(feeAmount), transaction.Amount.Add(feeAmount), sqlTime, transaction.Sender.AccountNumber)

		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}

	} else {
		// Drop onto ledger
	}

	// Update receiver account
	// Only update if account local
	if transaction.Receiver.BankNumber == "" {
		updateStatementReceiver := "UPDATE accounts SET `accountBalance` = (`accountBalance` + ?), `availableBalance` = (`availableBalance` + ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdReceiver, err := Config.Db.Prepare(updateStatementReceiver)
		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}
		defer stmtUpdReceiver.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdReceiver.Exec(transaction.Amount, transaction.Amount, sqlTime, transaction.Receiver.AccountNumber)

		if err != nil {
			return errors.New("payments.processCreditInitiation: " + err.Error())
		}
	} else {
		// Drop onto ledger
	}
	return
}

func processDepositInitiation(transaction PAINTrans, sqlTime int32, feeAmount decimal.Decimal) (err error) {
	// We don't update sender as it is deposit
	// Update receiver account
	// The total received amount is the deposited amount minus the fee
	depositTransactionAmount := transaction.Amount.Sub(feeAmount)
	// Only update if account local
	if transaction.Receiver.BankNumber == "" {
		updateStatementReceiver := "UPDATE accounts SET `accountBalance` = (`accountBalance` + ?), `availableBalance` = (`availableBalance` + ?), `timestamp` = ? WHERE `accountNumber` = ? "
		stmtUpdReceiver, err := Config.Db.Prepare(updateStatementReceiver)
		if err != nil {
			return errors.New("payments.processDepositInitiation: " + err.Error())
		}
		defer stmtUpdReceiver.Close() // Close the statement when we leave main() / the program terminates

		_, err = stmtUpdReceiver.Exec(depositTransactionAmount, depositTransactionAmount, sqlTime, transaction.Receiver.AccountNumber)

		if err != nil {
			return errors.New("payments.processDepositInitiation: " + err.Error())
		}
	} else {
		// Drop onto ledger
	}
	return
}

func getTransactionList(accountNumber string, offset int, perPage int) (allTransactions []PAINTrans, err error) {
	rows, err := Config.Db.Query("SELECT `type`, `senderAccountNumber`, `senderBankNumber`, `receiverAccountNumber`, `receiverBankNumber`, `transactionAmount`, `feeAmount`, `desc`, `timestamp`, `status`, `geo` FROM `transactions` WHERE `senderAccountNumber` = ? OR `receiverAccountNumber` = ?  ORDER BY `id` DESC LIMIT ?, ?", accountNumber, accountNumber, offset, perPage)
	if err != nil {
		return []PAINTrans{}, errors.New("transactions.ListTransactions: " + err.Error())
	}
	defer rows.Close()

	allTransactions = []PAINTrans{}
	for rows.Next() {
		transaction := PAINTrans{}
		if err := rows.Scan(&transaction.PainType, &transaction.Sender.AccountNumber, &transaction.Sender.BankNumber, &transaction.Receiver.AccountNumber, &transaction.Receiver.BankNumber, &transaction.Amount, &transaction.Fee, &transaction.Desc, &transaction.Timestamp, &transaction.Status, &transaction.Geo); err != nil {
			//@TODO Throw error
			return []PAINTrans{}, errors.New("transactions.ListTransactions: " + err.Error())
		}
		allTransactions = append(allTransactions, transaction)
	}

	return
}
