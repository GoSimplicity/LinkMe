# LinkMe项目启动
## clone项目源码到本地目录下
```bash
git@github.com:wangzijian2002/LinkMe.git
```
## 目录结构
```markdown
.
├── config           # 项目配置文件目录
├── docs             # API文档目录
├── go.mod           # Go模块定义文件
├── go.sum           # Go模块依赖校验和
├── internal         # 项目内部包，含核心业务逻辑
├── ioc              # IoC容器配置，负责依赖注入设置
├── LICENSE          # 开源许可证书
├── LinkMe           # 可执行文件或快捷方式
├── main.go          # 项目主入口文件
├── pkg              # 自定义工具包与库
├── README.md        # 项目自述文件
├── tmp              # 临时文件目录
├── wire_gen.go      # Wire工具生成的代码
├── wire.go          # Wire配置，声明依赖注入关系
└── yaml             # Kubernetes部署配置文件目录
```
## 拉取依赖包
```bash
go mod tidy
```
## 手动创建数据库/使用提供的yaml文件创建数据库
```bash
mysql -uroot -proot -e "create database linkme;"
或
kubectl apply -f yaml/  # 需要有k8s环境
```
## 使用wire进行依赖注入
```bash
go install github.com/google/wire/cmd/wire@latest
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