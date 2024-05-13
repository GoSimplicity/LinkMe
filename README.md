# LinkMe项目启动
## 首先需要clone项目源码到本地目录下
```bash
git@github.com:wangzijian2002/LinkMe.git
```
## 然后拉取依赖包
```bash
go mod tidy
```
## 构建并启动项目
```bash
go build -o linkme . && ./linkme
```
## 使用air启动项目(可选)
```air
go install github.com/cosmtrek/air@1.49.0 # 注意go版本不得低于air指定版本
air
```