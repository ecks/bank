#! /bin/bash
source ../scripts/json.sh
FQDN=`readJson ../config.json FQDN` || exit 1;

openssl s_client -showcerts -debug -connect $FQDN:3300 -no_ssl2 -bugs