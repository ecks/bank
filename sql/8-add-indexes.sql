CREATE UNIQUE INDEX account_num
ON accounts (accountNumber);

CREATE INDEX account_auth_num
ON accounts_auth (accountNumber);

CREATE UNIQUE INDEX account_meta_num
ON accounts_meta (accountNumber);

CREATE INDEX account_push_account_num
ON accounts_push_tokens (accountNumber);

CREATE INDEX bank_transactions_sender_num
ON bank_transactions (senderBankNumber);
CREATE INDEX bank_transactions_receiver_num
ON bank_transactions (receiverBankNumber);

CREATE INDEX transactions_sender_num
ON transactions (senderAccountNumber);
CREATE INDEX transactions_receiver_num
ON transactions (receiverAccountNumber);

/* Down
ALTER TABLE accounts
DROP INDEX account_num;

ALTER TABLE accounts_auth
DROP INDEX account_auth_num;

ALTER TABLE accounts_meta
DROP INDEX account_meta_num;

ALTER TABLE accounts_push_tokens
DROP INDEX account_push_account_num;

ALTER TABLE bank_transactions
DROP INDEX bank_transactions_sender_num;
ALTER TABLE bank_transactions 
DROP INDEX bank_transactions_receiver_num;

ALTER TABLE transactions
DROP INDEX transactions_sender_num;
ALTER TABLE transactions
DROP INDEX transactions_receiver_num;
*/
