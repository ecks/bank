ALTER TABLE transactions
ADD `status` enum('approved', 'rejected', 'pending') NOT NULL DEFAULT 'approved',
ADD `lat` POINT DEFAULT NULL,
ADD `long` POINT DEFAULT NULL,
ADD `desc` varchar(512) DEFAULT NULL
AFTER `feeAmount`;

--- http://dev.mysql.com/doc/refman/5.6/en/using-spatial-data.html
