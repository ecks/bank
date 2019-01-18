#!/bin/sh
source ./json.sh
PASS=`readJson ../config.json MySQLPass` || exit 1;
USER=`readJson ../config.json MySQLUser` || exit 1;
DB=`readJson ../config.json MySQLDB` || exit 1;

if [ ! "$USER" ]; then
    echo "Error: Cannot find USER in config.json" >&2;
    exit 1;
fi;

if [ ! "$DB" ]; then
    echo "Error: Cannot find Database in config.json" >&2;
    exit 1;
fi;

if [ ! "$PASS" ]; then
	CREDS="-u $USER"
else
	CREDS="-u $USER -p$PASS"
fi;
	mysql $CREDS -e "DROP DATABASE IF EXISTS $DB"
	mysql $CREDS -e "CREATE DATABASE $DB"
	mysql $CREDS -e "USE $DB"
	find ../sql -name '*.sql' | awk '{ print "source",$0 }' | sort -V | mysql $CREDS $DB --batch

#clear redis DB
redis-cli flushdb
