#!/bin/bash

# 创建数据库 请自行修改数据库密码等信息
mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS linkme DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 在 MySQL 中创建一个用于 Canal 连接的用户，并赋予必要的权限
mysql -uroot -proot -e "CREATE USER 'canal'@'%' IDENTIFIED BY 'canal';"
mysql -uroot -proot -e "GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';"
mysql -uroot -proot -e "FLUSH PRIVILEGES;"

# 导入SQL文件 请自行修改SQL文件路径
mysql -uroot -proot linkme < linkme.sql

# 提示完成
echo "Database import complete."
