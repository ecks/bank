#! /bin/bash
source ../scripts/json.sh

DOMAIN=`readJson ../config.json FQDN` || exit 1;

if [ ! "$DOMAIN" ]; then
    echo "Error: Cannot find FQDN in config.json" >&2;
    exit 1;
fi;

# Use the below lines to generate the certs
openssl req -new -nodes -x509 -out $DOMAIN.pem -keyout $DOMAIN.key -days 365
openssl req -new -nodes -x509 -out client.pem -keyout client.key -days 365
