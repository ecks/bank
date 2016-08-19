package push

import (
	"database/sql"
	"errors"

	"github.com/bvnk/bank/configuration"
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

func getPushTokens(accountNumber string) (pushDevices []PushDevice, err error) {
	rows, err := Config.Db.Query("SELECT DISTINCT `token`, `platform` FROM `accounts_push_tokens` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return []PushDevice{}, errors.New("push.getPushTokens: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		pushDevice := PushDevice{}
		err := rows.Scan(&pushDevice.Token, &pushDevice.Platform)
		if err != nil {
			break
		}
		pushDevices = append(pushDevices, pushDevice)
	}

	return
}
