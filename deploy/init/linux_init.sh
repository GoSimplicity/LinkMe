#!/bin/bash
mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS linkme DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;"
USER_EXISTS=$(mysql -uroot -proot -sse "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = 'canal');")
if [ "$USER_EXISTS" != 1 ]; then
  mysql -uroot -proot -e "CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';"
fi
mysql -uroot -proot -e "GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' WITH GRANT OPTION;"
mysql -uroot -proot -e "FLUSH PRIVILEGES;"
mysql -uroot -proot linkme < /docker-entrypoint-initdb.d/linkme.sql
echo "Database import complete."