GRANT REPLICATION SLAVE ON *.* TO 'replication'@'%' IDENTIFIED BY 'replication-password';
FLUSH PRIVILEGES;
