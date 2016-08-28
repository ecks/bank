package transactions

/*
PAIN transactions are as follows

Payments initiation:
1 - CustomerCreditTransferInitiationV06
2 - CustomerPaymentStatusReportV06
7 - CustomerPaymentReversalV05
8 - CustomerDirectDebitInitiationV05

Payments mandates:
9 - MandateInitiationRequestV04
10 - MandateAmendmentRequestV04
11 - MandateCancellationRequestV04
12 - MandateAcceptanceReportV04

#### Custom payments
1000 - CustomerDepositInitiation (@FIXME Will need to implement this properly, for now we use it to demonstrate functionality)
1001 - ListTransactions

*/

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bvnk/bank/accounts"
	"github.com/bvnk/bank/appauth"
	"github.com/bvnk/bank/push"
	"github.com/paulmach/go.geo"
	"github.com/shopspring/decimal"
)

const TRANSACTION_FEE = 0.0001 // 0.01%

// @TODO Have this struct not repeat in payments and accounts
type AccountHolder struct {
	AccountNumber string
	BankNumber    string
}

type PAINTrans struct {
	ID        int32
	PainType  int64
	Sender    AccountHolder
	Receiver  AccountHolder
	Amount    decimal.Decimal
	Fee       decimal.Decimal
	Geo       geo.Point
	Desc      string
	Status    string
	Timestamp int32
}

func ProcessPAIN(data []string) (result interface{}, err error) {
	//There must be at least 3 elements
	if len(data) < 3 {
		return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
	}

	// Get type
	painType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return "", errors.New("payments.ProcessPAIN: Could not get type of PAIN transaction. " + err.Error())
	}

	switch painType {
	case 1:
		//There must be at least 9 elements
		if len(data) < 9 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present.")
		}

		result, err = painCreditTransferInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	case 1000:
		//There must be at least 8 elements
		//token~pain~type~amount~lat~lon~desc
		if len(data) < 8 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present.")
		}
		result, err = customerDepositInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	case 1001:
		//token~pain~type~page~perpage
		if len(data) < 6 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present.")
		}
		result, err = listTransactions(data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	}

	return
}

func painCreditTransferInitiation(painType int64, data []string) (result string, err error) {

	// Validate input
	sender, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	receiver, err := parseAccountHolder(data[4])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	// Check if sender valid
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	err = accounts.CheckUserAccountValidFromToken(tokenUser, sender.AccountNumber)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: Sender not valid")
	}

	lat, err := strconv.ParseFloat(data[6], 64)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: Could not parse coordinates into float")
	}
	lon, err := strconv.ParseFloat(data[7], 64)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: Could not parse coordinates into float")
	}
	desc := data[8]

	geo := *geo.NewPoint(lat, lon)
	transaction := PAINTrans{0, painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), geo, desc, "approved", 0}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable, err := checkBalance(transaction.Sender)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	// Comparing decimals results in -1 if <
	if balanceAvailable.Cmp(transaction.Amount) == -1 {
		return "", errors.New("payments.painCreditTransferInitiation: Insufficient funds available")
	}

	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	go push.SendNotification(sender.AccountNumber, "💸 Payment sent!", 1, "default")
	go push.SendNotification(receiver.AccountNumber, "💸 Payment received!", 1, "default")

	return
}

func processPAINTransaction(transaction PAINTrans) (result string, err error) {
	// Test: pain~1~1b2ca241-0373-4610-abad-da7b06c50a7b@~181ac0ae-45cb-461d-b740-15ce33e4612f@~20

	// Save in transaction table
	err = savePainTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.processPAINTransaction: " + err.Error())
	}

	// Amend sender and receiver accounts
	// Amend bank's account with fee addition
	err = updateAccounts(transaction)
	if err != nil {
		return "", errors.New("payments.processPAINTransaction: " + err.Error())
	}

	return
}

func parseAccountHolder(account string) (accountHolder AccountHolder, err error) {
	accountStr := strings.Split(account, "@")

	if len(accountStr) < 2 {
		return AccountHolder{}, errors.New("payments.parseAccountHolder: Not all details present")
	}

	accountHolder = AccountHolder{accountStr[0], accountStr[1]}
	return
}

func customerDepositInitiation(painType int64, data []string) (result string, err error) {
	// Validate input
	// Sender is bank
	sender, err := parseAccountHolder("0@0")
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	receiver, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	trAmt := strings.TrimRight(data[4], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	// Check if sender valid
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: " + err.Error())
	}
	err = accounts.CheckUserAccountValidFromToken(tokenUser, receiver.AccountNumber)
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: Sender not valid")
	}

	lat, err := strconv.ParseFloat(data[5], 64)
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: Could not parse coordinates into float")
	}
	lon, err := strconv.ParseFloat(data[6], 64)
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: Could not parse coordinates into float")
	}
	desc := data[7]

	// Issue deposit
	// @TODO This flow show be fixed. Maybe have banks approve deposits before initiation, or
	// immediate approval below a certain amount subject to rate limiting
	geo := *geo.NewPoint(lat, lon)
	transaction := PAINTrans{0, painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), geo, desc, "approved", 0}
	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	go push.SendNotification(receiver.AccountNumber, "💸 Deposit received!", 1, "default")

	return
}

func listTransactions(data []string) (result []PAINTrans, err error) {
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return []PAINTrans{}, errors.New("payments.ListTransactions: " + err.Error())
	}
	accountNumber := data[3]
	err = accounts.CheckUserAccountValidFromToken(tokenUser, accountNumber)
	if err != nil {
		return []PAINTrans{}, errors.New("payments.ListTransactions: " + err.Error())
	}

	page, err := strconv.Atoi(data[4])
	if err != nil {
		return []PAINTrans{}, errors.New("payments.ListTransactions: " + err.Error())
	}
	perPage, err := strconv.Atoi(data[5])
	if err != nil {
		return []PAINTrans{}, errors.New("payments.ListTransactions: " + err.Error())
	}
	// We limit perPage to 100
	if perPage > 100 {
		return []PAINTrans{}, errors.New("payments.ListTransactions: Cannot retrieve more than 100 results per request")
	}

	// Check if timestamp present
	if data[6] != "" {
		timestamp, err := strconv.Atoi(data[6])
		if err != nil {
			return []PAINTrans{}, errors.New("payments.ListTransactions: Could not convert timestamp to int. " + err.Error())
		}

		result, err = getTransactionListAfterTimestamp(accountNumber, (page * perPage), perPage, timestamp)
		if err != nil {
			return []PAINTrans{}, errors.New("payments.ListTransactions: " + err.Error())
		}

		return result, nil
	}

	result, err = getTransactionList(accountNumber, (page * perPage), perPage)
	if err != nil {
		return []PAINTrans{}, errors.New("payments.ListTransactions: " + err.Error())
	}

	return
}
