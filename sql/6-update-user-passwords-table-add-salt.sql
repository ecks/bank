ALTER TABLE accounts_auth 
ADD `salt` char(128) NOT NULL
AFTER `password`;
