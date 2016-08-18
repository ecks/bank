![Bank logo](https://bvnk.co/bvnk-dark-hor.png)
# Banking infrastructure built using Go
_Core banking in Golang_

[![Join the chat at https://gitter.im/ksred/bank](https://badges.gitter.im/ksred/bank.svg)](https://gitter.im/ksred/bank?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Overview

__Bank__ is an answer to the question “What would banking infrastructure look like if it were built today?”

The financial industry is well known for using outdated technology, some still relying on systems from the 1970s written using languages like COBOL. This project aims to bring that infrastructure through to current technological levels.

A secondary aim of this project is to lift the curtain on the financial industry. This is one of the most closed off industries, and finding out how it works has proven to be very difficult. As more is discovered, the possibilities open up. Bringing light to this industry will surely help other fintech developers and upstarts.

This project’s [development is documented](http://blog.ksred.me/tags/banking-infrastructure/).

### Standards (WIP)
- [ISO20222](http://www.iso20022.org/)
- [ISO 8583](https://en.wikipedia.org/wiki/ISO_8583)

### Functionality (WIP)
- Payments
- Securities
- Trade Services
- FX
- Cards

## Installation

__Coming soon:__ There will soon be a public facing instance of this server to build and test against.

This project is built using `Go 1.6`, but may work against older versions.

### Software Dependencies:
- MySQL
- Redis

In order to develop against this project, clone it using `git clone https://github.com/ksred/bank.git`. Make sure you have all dependencies the project uses:

- `go get github.com/go-sql-driver/mysql`
- `go get github.com/satori/go.uuid`
- `go get gopkg.in/redis.v3`
- `go get github.com/gorilla/mux`
- `go get github.com/pzduniak/argon2`
- `go get github.com/paulmach/go.geo`
- `go get github.com/kardianos/osext`

You can then `go build`.

## Install certs
Either create your own self-signed certs, or purchase them, and place them in the `certs/` directory. There is a script `generateCerts.sh` to assist with this. These certs must keep the same naming convention:

- `server.key`
- `server.pem`
- `client.key`
- `client.pem`

## Running the HTTP API

Once built, you can run the HTTP API:

- Secure: `./bank http`
- Insecure: There is no insecure mode implemented

The API will then be available at the [FQDN and port specified](https://github.com/ksred/bank/blob/master/main.go#L11).

## Running the CLI server

You can run the CLI server:

- Secure: `./bank server`
- Insecure: `./bank serverNoTLS`

## Running the CLI client

The program comes with a CLI client. A TLS client can only connect to a TLS server, so the modes must match. The server accepts commands from the CLI separated by new lines. _A detailed explanation will be placed in the wiki._

You can run the client:

- Secure: `./bank client`
- Insecure: `./bank clientNoTLS`

## Example flow

The following illustrates an example flow for using the server through the CLI client.

__Create account__
```
0~acmt~1~Kyle~Redelinghuys~19700101~197001011234098~55512340987~~email@domain.com~Physical Address 1~~~1000
```

This call responds with the `user_id`:
```
1~52d27bde-9418-4a5d-8528-3fb32e1a5d69
```

__Create account login__
```
0~appauth~3~52d27bde-9418-4a5d-8528-3fb32e1a5d69~TestPassword
```

__Log into account__
```
0~appauth~2~52d27bde-9418-4a5d-8528-3fb32e1a5d69~TestPassword
```

This responds with a token:
```
1~cb485f9d-0a24-4385-a358-61ea0d44fdea
```

This token is valid for 60 minutes by default.

__Make a payment, here the payment amount is 20__
```
cb485f9d-0a24-4385-a358-61ea0d44fdea~pain~1~52d27bde-9418-4a5d-8528-3fb32e1a5d69@~137232cc-142e-474c-aaaa-43393f9b7c4c@~20
```

A full wiki and documentation will be added.

## License
[Q Public License 1.0 (QPL-1.0)](https://en.wikipedia.org/wiki/Q_Public_License)

Copyright Kyle Redelinghuys
