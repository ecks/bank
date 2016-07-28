package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/bvnk/bank/accounts"
	"github.com/bvnk/bank/appauth"
	"github.com/bvnk/bank/configuration"
	"github.com/bvnk/bank/payments"
	"github.com/bvnk/bank/push"
)

var Config configuration.Configuration

func runServer(mode string) (message string, err error) {

	// Load app config
	Config, err := configuration.LoadConfig()
	if err != nil {
		return "", errors.New("server.runServer: " + err.Error())
	}

	// Set config in packages
	accounts.SetConfig(&Config)
	payments.SetConfig(&Config)
	appauth.SetConfig(&Config)
	push.SetConfig(&Config)

	switch mode {
	case "tls":
		cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
		if err != nil {
			return "", err
		}

		// Load config and generate seed
		config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
		config.Rand = rand.Reader

		// Listen for incoming connections.
		l, err := tls.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT, &config)
		if err != nil {
			return "", err
		}

		// Close the listener when the application closes.
		defer l.Close()
		fmt.Println("Listening on secure " + CONN_HOST + ":" + CONN_PORT)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				return "", err
			}
			// Handle connections in a new goroutine.
			go handleTCPRequest(conn)
		}
	case "no-tls":
		// Listen for incoming connections.
		l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
		if err != nil {
			return "", err
		}

		// Close the listener when the application closes.
		defer l.Close()
		fmt.Println("Listening on unsecure " + CONN_HOST + ":" + CONN_PORT)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				return "", err
			}
			// Handle connections in a new goroutine.
			go handleTCPRequest(conn)
		}
	}

	return
}

// Handles incoming requests.
func handleTCPRequest(conn net.Conn) (err error) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}
	s := string(buf[:])

	// Process
	result, err := processCommand(s)

	// Convert response to text
	textResponse := "1~" + result
	if err != nil {
		textResponse = "0~" + err.Error()
	}

	// Send a response back to person contacting us.
	conn.Write([]byte(textResponse + "\n"))
	// Close the connection when you're done with it.
	conn.Close()

	return
}

func processCommand(text string) (result string, err error) {
	// Commands are received split by tilde (~)
	// command~DATA
	cleanText := strings.Replace(text, "\n", "", -1)
	fmt.Printf("### %s ####\n", cleanText)
	command := strings.Split(cleanText, "~")

	// Check if we received a command
	if len(command) == 0 {
		fmt.Println("No command received")
		return
	}

	// Remove null termination from data
	command[len(command)-1] = string(bytes.Trim([]byte(command[len(command)-1]), "\x00"))

	// Check application auth. This is always the first value, if no token a 0 is sent
	isCreateAccount := (command[0] == "0" && command[1] == "acmt" && command[2] == "1")
	isLogIn := (command[0] == "0" && command[1] == "appauth" && command[2] == "2")
	isCreateUserPassword := (command[0] == "0" && command[1] == "appauth" && command[2] == "3")

	if !isCreateAccount || !isLogIn || !isCreateUserPassword {
		err := appauth.CheckToken(command[0])
		if err != nil {
			return "", errors.New("server.processCommand: " + err.Error())
		}
	}

	switch command[1] {
	case "appauth":
		// Check "help"
		if command[2] == "help" {
			return "Format of appauth: appauth~userName~password", nil
		}
		result, err = appauth.ProcessAppAuth(command)
		if err != nil {
			return "", errors.New("server.processCommand: " + err.Error())
		}
		break
	case "pain":
		// Check "help"
		if command[2] == "help" {
			return "Format of PAIN transaction:\npain\npainType~senderAccountNumber@SenderBankNumber\nreceiverAccountNumber@ReceiverBankNumber\ntransactionAmount\n\nBank numbers may be left void if bank is local", nil
		}
		result, err = payments.ProcessPAIN(command)
		if err != nil {
			return "", errors.New("server.processCommand: " + err.Error())
		}
	case "camt":
	case "acmt":
		// Check "help"
		if command[2] == "help" {
			return "", nil // @TODO Help section
		}
		result, err = accounts.ProcessAccount(command)
	case "remt":
	case "reda":
	case "pacs":
	case "auth":
		break
	default:
		return "No valid command received", nil
	}

	return
}
