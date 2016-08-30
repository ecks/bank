/*
Accounts users accounts table links a single accounts user to many accounts
*/
CREATE TABLE IF NOT EXISTS accounts_users_accounts (
`id` int NOT NULL AUTO_INCREMENT,
`accountHolderIdentificationNumber` text NOT NULL, 
`accountNumber` char(36) NOT NULL, 
`bankNumber` char(36) NOT NULL, 
`timestamp` int NOT NULL, 
PRIMARY KEY (`id`)
);
