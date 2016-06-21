package main

import (
	"errors"
	"log"
	"os"
)

const (
	// This is the FQDN from the certs generated
	CONN_HOST = "thebankoftoday.com"
	CONN_PORT = "3300"
	CONN_TYPE = "tcp"
	HTTP_PORT = "8443"
)

func main() {
	argClientServer := "http"
	// http server is default mode
	if len(os.Args) > 1 {
		argClientServer = os.Args[1]
	}

	err := parseArguments(argClientServer)
	if err != nil {
		log.Fatalf("Error starting, err: %v\n", err)
	}
	os.Exit(0)
}

func parseArguments(arg string) (err error) {
	switch arg {
	case "http":
		err := RunHttpServer()
		if err != nil {
			log.Fatalf("Could not start HTTP server. " + err.Error())
		}
		break
	case "client":
		// Run client for bank system
		runClient("tls")
		break
	case "clientNoTLS":
		// Run client for bank system
		runClient("no-tls")
		break
	case "server":
		// Run server for bank system
		for {
			runServer("tls")
		}
	case "serverNoTLS":
		// Run server for bank system
		for {
			runServer("no-tls")
		}
	default:
		return errors.New("No valid option chosen. Valid options: client, clientNoTLS, server, serverNoTLS")
	}

	return
}
