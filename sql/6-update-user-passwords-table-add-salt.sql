ALTER TABLE accounts_auth 
ADD `salt` char(64) NOT NULL
AFTER `password`;
