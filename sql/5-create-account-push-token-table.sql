CREATE TABLE IF NOT EXISTS accounts_push_tokens (
`id` int NOT NULL AUTO_INCREMENT,
`accountNumber` char(36) NOT NULL, 
`token` varchar(255) NOT NULL,
`platform` enum('ios', 'android', 'blackberry', 'windows', 'other') NOT NULL,
`timestamp` int NOT NULL,
PRIMARY KEY (`id`)
);

