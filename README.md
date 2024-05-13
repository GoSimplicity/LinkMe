# LinkMe项目启动
## 首先需要clone项目源码到本地目录下
```bash
git@github.com:wangzijian2002/LinkMe.git
```
## 然后拉取依赖包
```bash
go mod tidy
```
## 手动创建数据库
```bash
mysql -uroot -proot -e "create database linkme;"
```
## 使用wire进行依赖注入
```bash
wire # 注意需要在wire.go文件所在目录下使用
```
## 构建并启动项目
```bash
go build -o linkme . && ./linkme
```
## 使用air启动项目(可选)
```bash
go install github.com/cosmtrek/air@1.49.0 # 注意go版本不得低于air指定版本，本项目使用golang版本为1.22
air
```