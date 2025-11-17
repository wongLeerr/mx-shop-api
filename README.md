# mx-shop-api
调用底层grpc服务暴露为上层http服务。

# user-web 目录
负责暴露底层user的grpc服务为上层http服务。

# go日志库 zap
分为logger和sugarLogger，sugarLogger提供简单易用的日志打印api，logger打印日志api用起来稍复杂但是性能极致。
日志是分级别的，例如分开发环境、生产环境。
debug、info、warn、error、fetal。
zap.L是zap.Logger的简易调用方式，zap.S是zap.SugaredLogger的简易调用方式，前者性能更好但需明确说明数据类型，后者调用更方便。