ALTER DATABASE bank CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

ALTER TABLE transactions CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER TABLE transactions CHANGE `desc` `desc` VARCHAR(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

/*
/etc/mysql/my.cnf needs updating:

[client]
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4

[mysqld]
character-set-client-handshake = FALSE
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci

THEN

mysqlcheck -u root -p --auto-repair --optimize --all-databases
*/
