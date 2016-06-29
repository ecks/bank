package appauth

import (
	"encoding/hex"
	"errors"
	"time"

	"gopkg.in/redis.v3"

	"github.com/ksred/bank/configuration"
	"github.com/pzduniak/argon2"
	"github.com/satori/go.uuid"
)

const (
	TOKEN_TTL           = time.Hour // One hour
	MIN_PASSWORD_LENGTH = 8
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func ProcessAppAuth(data []string) (result string, err error) {
	//@TODO: Change from []string to something more solid, struct/interface/key-pair
	if len(data) < 3 {
		return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
	}
	switch data[2] {
	// Auth an existing account
	case "1":
		// TOKEN~appauth~1
		if len(data) < 3 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		err := CheckToken(data[0])
		if err != nil {
			return "", err
		}
		return result, nil
	// Log in
	case "2":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = CreateToken(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	// Create an account
	case "3":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = CreateUserPassword(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	// Remove an account
	case "4":
		if len(data) < 5 {
			return "", errors.New("appauth.ProcessAppAuth: Not all required fields present")
		}
		result, err = RemoveUserPassword(data[3], data[4])
		if err != nil {
			return "", err
		}
		return result, nil
	}
	return "", errors.New("appauth.ProcessAppAuth: No valid option chosen")
}

func CreateUserPassword(user string, password string) (result string, err error) {
	//TEST 0~appauth~3~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword

	// Check for existing account
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count > 0 {
		return "", errors.New("appauth.CreateUserPassword: Account already exists")
	}

	// Check password length
	if len(password) < MIN_PASSWORD_LENGTH {
		return "", errors.New("appauth.CreateUserPassword: Password must be at least " + string(MIN_PASSWORD_LENGTH) + " characters")
	}

	// Generate hash
	output, err := argon2.Key([]byte(password), []byte(Config.PasswordSalt), 3, 4, 4096, 32, argon2.Argon2i)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Could not generate secure hash. " + err.Error())
	}

	hash := hex.EncodeToString(output)

	// Prepare statement for inserting data
	insertStatement := "INSERT INTO accounts_auth (`accountNumber`, `password`, `timestamp`) "
	insertStatement += "VALUES(?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Error with insert. " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Convert variables
	t := time.Now()
	sqlTime := int32(t.Unix())

	_, err = stmtIns.Exec(user, hash, sqlTime)

	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Could not save account. " + err.Error())
	}

	result = "Successfully created account"
	return
}

func RemoveUserPassword(user string, hashedPassword string) (result string, err error) {
	// Check for existing account
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Error with select query. " + err.Error())
	}
	defer rows.Close()

	// @TODO Must be easy way to get row count returned
	count := 0
	for rows.Next() {
		count++
	}

	if count == 0 {
		return "", errors.New("appauth.RemoveUserPassword: Account auth does not exists")
	}

	// Prepare statement for inserting data
	delStatement := "DELETE FROM accounts_auth WHERE `accountNumber` = ? AND `password` = ? "
	stmtDel, err := Config.Db.Prepare(delStatement)
	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Error with delete. " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(user, hashedPassword)

	if err != nil {
		return "", errors.New("appauth.RemoveUserPassword: Could not delete account. " + err.Error())
	}

	result = "Successfully deleted account"
	return
}

func CreateToken(user string, password string) (token string, err error) {
	rows, err := Config.Db.Query("SELECT `password` FROM `accounts_auth` WHERE `accountNumber` = ?", user)
	if err != nil {
		return "", errors.New("appauth.CreateToken: Error with select query. " + err.Error())
	}
	defer rows.Close()

	count := 0
	hashedPassword := ""
	for rows.Next() {
		if err := rows.Scan(&hashedPassword); err != nil {
			return "", errors.New("appauth.CreateToken: Could not retreive account details")
		}
		count++
	}

	// Generate hash
	output, err := argon2.Key([]byte(password), []byte(Config.PasswordSalt), 3, 4, 4096, 32, argon2.Argon2i)
	if err != nil {
		return "", errors.New("appauth.CreateUserPassword: Could not generate secure hash. " + err.Error())
	}

	hash := hex.EncodeToString(output)

	if hash != hashedPassword {
		return "", errors.New("appauth.CreateToken: Authentication credentials invalid")
	}

	newUuid := uuid.NewV4()
	token = newUuid.String()

	// @TODO Remove all tokens for this user
	err = Config.Redis.Set(token, user, TOKEN_TTL).Err()
	if err != nil {
		return "", errors.New("appauth.CreateToken: Could not set token. " + err.Error())
	}

	return
}

func RemoveToken(token string) (result string, err error) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	_, err = Config.Redis.Del(token).Result()

	if err == redis.Nil {
		return "", errors.New("appauth.RemoveToken: Token not found. " + err.Error())
	} else if err != nil {
		return "", errors.New("appauth.RemoveToken: Could not remove token. " + err.Error())
	} else {
		result = "Token removed"
	}

	return
}

func CheckToken(token string) (err error) {
	//TEST 0~appauth~480e67e3-e2c9-48ee-966c-8d251474b669
	user, err := Config.Redis.Get(token).Result()

	if err == redis.Nil {
		return errors.New("appauth.CheckToken: Token not found. " + err.Error())
	} else if err != nil {
		return errors.New("appauth.CheckToken: Could not get token. " + err.Error())
	} else {
		// Extend token
		err := Config.Redis.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			return errors.New("appauth.CheckToken: Could not extend token. " + err.Error())
		}
	}

	return
}

func GetUserFromToken(token string) (user string, err error) {
	//TEST 0~appauth~~181ac0ae-45cb-461d-b740-15ce33e4612f~testPassword
	user, err = Config.Redis.Get(token).Result()
	if err != nil {
		return "", errors.New("appauth.GetUserFromToken: Could not get token. " + err.Error())
	}

	// If valid then extend
	if user != "" {
		err := Config.Redis.Set(user, token, TOKEN_TTL).Err()
		if err != nil {
			return "", errors.New("appauth.GetUserFromToken: Could not extend token. " + err.Error())
		}
	}

	return
}
