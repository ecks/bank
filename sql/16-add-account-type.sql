ALTER TABLE accounts 
ADD `type` enum('savings', 'cheque', 'merchant', 'money-market', 'cd', 'ira', 'rcp', 'credit', 'mortgage', 'loan') NOT NULL DEFAULT 'cheque'
AFTER `availableBalance`;

