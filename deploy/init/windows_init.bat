@echo off
REM 创建数据库
mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS linkme DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;"

REM 检查用户是否存在
for /f "tokens=*" %%i in ('mysql -uroot -proot -sse "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = 'canal');"') do set USER_EXISTS=%%i

if "%USER_EXISTS%"=="0" (
    mysql -uroot -proot -e "CREATE USER 'canal'@'%%' IDENTIFIED BY 'canal';"
)

REM 授予权限
mysql -uroot -proot -e "GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%%' WITH GRANT OPTION;"
mysql -uroot -proot -e "FLUSH PRIVILEGES;"

REM 获取当前目录
set currentDir=%~dp0

REM 导入当前目录下的SQL文件
mysql -uroot -proot linkme < "%currentDir%linkme.sql"

REM 提示完成
echo Database import complete.
pause