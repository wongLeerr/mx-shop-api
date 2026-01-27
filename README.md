# mx-shop-api
调用底层grpc服务暴露为上层http服务。

# user-web 目录
负责暴露底层user的grpc服务为上层http服务。

# go日志库 zap
分为logger和sugarLogger，sugarLogger提供简单易用的日志打印api，logger打印日志api用起来稍复杂但是性能极致。
日志是分级别的，例如分开发环境、生产环境。
debug、info、warn、error、fetal。
zap.L是zap.Logger的简易调用方式，zap.S是zap.SugaredLogger的简易调用方式，前者性能更好但需明确说明数据类型，后者调用更方便。

# 使用 protoc 生成 go 代码
生成普通proto结构体代码: `protoc --go_out=. user.proto`
生成gRPC service接口代码:  `protoc --go-grpc_out=. user.proto`

# DTO
DTO（Data Transfer Object）

# go的配置文件管理
viper
why viper? 支持默认值、监听配置文件变动、很多简单易用的能力。

# redis
基于内存的 Key-Value 数据库

启动 redis：`brew services start redis`
测试 redis 是否运行成功：`redis-cli ping`
启动redis服务端：`redis-server`
启动redis客户端：`redis-cli`

| 配置项          | 值               |
| ------------ | --------------- |
| **Host**     | `127.0.0.1`     |
| **Port**     | `6379`          |
| **Password** | 空（如果你没设置密码的话）   |

# 服务注册 服务发现
启动 consul：`consul agent -dev`
访问可视化界面：`http://localhost:8500/ui`
使用dig解析服务name对应ip和port：`dig @127.0.0.1 -p 8600 web.service.consul` 其中：dig @127.0.0.1 -p 8600 表示连接本地的consul服务，解析web服务。xxx.service.consul中xxx表示服务名称，后面为固定写法。

# 常用的负载均衡算法
1. 轮询法（Round Robin） （平均将请求分配给各个服务器）
2. 随机法，同一
3. 源地址哈希法（大致意思是通过某种算法，使得同一个客户端IP访问的始终是同一台服务器）
4. 加权轮询（考虑机器性能等情况）
5. 加权随机（考虑机器性能等情况）
6. 最小连接数（考虑服务器的连接数，将请求分配给连接数较小的服务器）

# 分布式配置中心选型
apollo: 携程开源，大而全
nacos: 阿里开源，小而全

### nacos
Nacos 是一个 Java 服务，本质是一个 Web 应用
本地启动nacos：进入到dev目录 进入nacos目录执行：`sh bin/startup.sh -m standalone`
访问：`http://127.0.0.1:8848/nacos`

nacos中的一些概念：
- 命名空间 - 可以理解为一个项目就可以创建一个命令空间，例如user-web一个命名空间，user-srv一个命名空间。
- 组 - 可以用来做环境隔离，例如dev组，prod组。
- 配置集（data id）可以理解为具体的配置文件。

### ngrok
内网穿透工具，可将本地指定端口服务映射为公网IP。
用法: `ngrok http xxxx` xxxx 为本地服务端口。

