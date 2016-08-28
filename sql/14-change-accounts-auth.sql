RENAME TABLE `accounts_auth` TO `accounts_user_auth`;

/*
@TODO Create migration where the current accountNumbers are replaced with the users identification number
*/

DROP INDEX account_auth_num ON `accounts_user_auth`;
ALTER TABLE `accounts_user_auth` CHANGE `accountNumber` `accountHolderIdentificationNumber` VARCHAR(200);
CREATE INDEX account_auth_num ON accounts_user_auth(accountHolderIdentificationNumber);


