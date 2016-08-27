/*
We copy the existing data over to the new table
*/
INSERT INTO accounts_users_accounts (`accountHolderIdentificationNumber`, `accountNumber`, `bankNumber`, `timestamp`)
SELECT `accountHolderIdentificationNumber`, `accountNumber`, `bankNumber`, `timestamp`
FROM accounts_users;
