ALTER TABLE accounts_user_auth 
ADD `authUser` char(36) NOT NULL
AFTER `accountHolderIdentificationNumber`;

/* Here we reset the table as the authUser will be blank */
TRUNCATE TABLE accounts_user_auth;
