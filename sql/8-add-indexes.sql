CREATE UNIQUE INDEX account_num
ON accounts (accountNumber);

CREATE UNIQUE INDEX account_auth_num
ON accounts_auth (accountNumber);

CREATE UNIQUE INDEX account_meta_num
ON accounts_meta (accountNumber);

CREATE UNIQUE INDEX account_push_account_num
ON accounts_push_tokens (accountNumber);

CREATE UNIQUE INDEX bank_transactions_sender_num
ON bank_transactions (senderBankNumber);
CREATE UNIQUE INDEX bank_transactions_receiver_num
ON bank_transactions (receiverBankNumber);

CREATE UNIQUE INDEX transactions_sender_num
ON transactions (senderAccountNumber);
CREATE UNIQUE INDEX transactions_receiver_num
ON transactions (receiverAccountNumber);
