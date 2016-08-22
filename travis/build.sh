#! /bin/bash

# Copy config
cp $HOME/gopath/src/github.com/ksred/bank/travis/config.json $HOME/gopath/src/github.com/ksred/bank/

# Create database
mysql -u root -e "create database bank;"
mysql -u root bank < $HOME/gopath/src/github.com/ksred/bank/sql/0-create-trans-table.sql
mysql -u root bank < $HOME/gopath/src/github.com/ksred/bank/sql/1-create-accounts-table.sql
mysql -u root bank < $HOME/gopath/src/github.com/ksred/bank/sql/2-create-account-meta-table.sql
mysql -u root bank < $HOME/gopath/src/github.com/ksred/bank/sql/3-create-bank-holding-account.sql
mysql -u root bank < $HOME/gopath/src/github.com/ksred/bank/sql/4-create-account-auth-table.sql
mysql -u root bank < $HOME/gopath/src/github.com/ksred/bank/sql/5-create-account-push-token-table.sql