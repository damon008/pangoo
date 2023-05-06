# pangoo
全新高性能底层依赖模块

# 版权声明

开源不代表免费，本项目遵循 [GPL-3.0](https://github.com/damon008/pangoo/blob/main/LICENSE) 开源协议发布，并提供技术交流学习，但**绝不允许修改后和衍生的代码做为闭源的商业软件发布和销售！** 如果需要将本产品在本地进行任何附带商业化性质行为使用，**请联系项目负责人进行商业授权**，以遵守 GPL 协议保证您的正常使用。

目前在国内 GPL 协议**具备合同特征，是一种民事法律行为** ，属于我国《合同法》调整的范围。本项目原始团队保留一切诉讼权利。

[相关案例：违反 GPL 协议赔偿 50 万，国内首例!](https://mp.weixin.qq.com/s/YQ6sNjbDS-P7BViLZIsaoA)

我们团队拥有对本开源协议的最终解释权。


## 相关技术支持

| 拓展                                                                                                 | 描述                                                                                                |
|----------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------|
| [Autotls](https://github.com/hertz-contrib/autotls)                                                | 为 Hertz 支持 Let's Encrypt 。                                                                        |
| [Http2](https://github.com/hertz-contrib/http2)                                                    | 提供对 HTTP2 的支持。                                                                                    |
| [Websocket](https://github.com/hertz-contrib/websocket)                                            | 使 Hertz 支持 Websocket 协议。                                                                          |
| [Etag](https://github.com/hertz-contrib/etag)                                                      | 提供 ETag HTTP 响应标头。                                                                                |
| [Limiter](https://github.com/hertz-contrib/limiter)                                                | 提供了基于 bbr 算法的限流器。                                                                                 |
| [Monitor-prometheus](https://github.com/hertz-contrib/monitor-prometheus)                          | 提供基于 Prometheus 服务监控功能。                                                                           |
| [Obs-opentelemetry](https://github.com/hertz-contrib/obs-opentelemetry)                            | Hertz 的 Opentelemetry 扩展，支持 Metric、Logger、Tracing 并且达到开箱即用。                                       |
| [Opensergo](https://github.com/hertz-contrib/opensergo)                                            | Opensergo 扩展。                                                                                     |
| [Pprof](https://github.com/hertz-contrib/pprof)                                                    | Hertz 集成 Pprof 的扩展。                                                                               |
| [Registry](https://github.com/hertz-contrib/registry)                                              | 提供服务注册与发现功能。到现在为止，支持的服务发现拓展有 nacos， consul， etcd， eureka， polaris， servicecomb， zookeeper， redis。 |
| [Sentry](https://github.com/hertz-contrib/hertzsentry)                                             | Sentry 拓展提供了一些统一的接口来帮助用户进行实时的错误监控。                                                                |
| [Tracer](https://github.com/hertz-contrib/tracer)                                                  | 基于 Opentracing 的链路追踪。                                                                             |
| [Basicauth](https://github.com/cloudwego/hertz/tree/develop/pkg/app/middlewares/server/basic_auth) | Basicauth 中间件能够提供 HTTP 基本身份验证。                                                                    |
| [Jwt](https://github.com/hertz-contrib/jwt)                                                        | Jwt 拓展。                                                                                           |
| [Keyauth](https://github.com/hertz-contrib/keyauth)                                                | 提供基于 token 的身份验证。                                                                                 |
| [Requestid](https://github.com/hertz-contrib/requestid)                                            | 在 response 中添加 request id。                                                                        |
| [Sessions](https://github.com/hertz-contrib/sessions)                                              | 具有多状态存储支持的 Session 中间件。                                                                           |
| [Casbin](https://github.com/hertz-contrib/casbin)                                                  | 通过 Casbin 支持各种访问控制模型。                                                                             |
| [Cors](https://github.com/hertz-contrib/cors)                                                      | 提供跨域资源共享支持。                                                                                       |
| [Csrf](https://github.com/hertz-contrib/csrf)                                                      | Csrf 中间件用于防止跨站点请求伪造攻击。                                                                            |
| [Secure](https://github.com/hertz-contrib/secure)                                                  | 具有多配置项的 Secure 中间件。                                                                               |
| [Gzip](https://github.com/hertz-contrib/gzip)                                                      | 含多个可选项的 Gzip 拓展。                                                                                  |
| [I18n](https://github.com/hertz-contrib/i18n)                                                      | 可帮助将 Hertz 程序翻译成多种语言。                                                                             |
| [Lark](https://github.com/hertz-contrib/lark-hertz)                                                | 在 Hertz 中处理 Lark/飞书的卡片消息和事件的回调。                                                                   |
| [Loadbalance](https://github.com/hertz-contrib/loadbalance)                                        | 提供适用于 Hertz 的负载均衡算法。                                                                              |
| [Logger](https://github.com/hertz-contrib/logger)                                                  | Hertz 的日志拓展，提供了对 zap、logrus、zerologs 日志框架的支持。                                                     |
| [Recovery](https://github.com/cloudwego/hertz/tree/develop/pkg/app/middlewares/server/recovery)    | Hertz 的异常恢复中间件。                                                                                   |
| [Reverseproxy](https://github.com/hertz-contrib/reverseproxy)                                      | 实现反向代理。                                                                                           |
| [Swagger](https://github.com/hertz-contrib/swagger)                                                | 使用 Swagger 2.0 自动生成 RESTful API 文档。                                                               |
| [Cache](https://github.com/hertz-contrib/cache)                                                    | 用于缓存 HTTP 接口内容的 Hertz 中间件，支持多种客户端。                                                                |


### 压力测试结论

1. 正常的请求下，Kitex-mux性能最佳，Kitex次之，建议优先使用Thrift编解码

2. 基于GRPC协议下，Kitex性能最佳，比原生gRPC强，建议优先使用Thrift编解码

3. 基于Streaming下，Kitex性能最佳，基于Protobuf编码，有两种：Kitex Protobuf 和 gRPC。

4. 基于Thrift、Protobuf下，Kitex-mux性能最佳，Kitex次之，建议优先使用Thrift编解码

Kitex 在使用 Thrift 作为 Payload 的情况下，性能优于官方 gRPC，吞吐接近 gRPC 的两倍；
此外，在 Kitex 使用定制的 Protobuf 协议时，性能也优于 gRPC。

### 项目结构

- hello-server 此为kitex server端demo
- producer-service 此为kitex cli端demo，同时又是hertz server端demo
- consumer-service 此为hertz cli端demo，同时又是hertz server端demo
- idl 此为项目代码生成的idl文件所在，thrift、pb文件
- kitex_gen 此为根据idl的文件生成的kitex所需要使用到的中间件，便于kitex cli和server端互通
- pkg 此为所有的依赖包、函数、工具等资源池

### 代码生成

```shell
go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/cloudwego/kitex/tool/server/kitex@latest
go install github.com/cloudwego/thriftgo@latest

#kitex server 发现不太友好：去掉,template=slim
kitex -module "primus" -thrift frugal_tag -service hello idl/hello.thrift
mkdir hello-server
cd hello-server
...

#code server for hertz
#不需要编解码代码加参数-t=template=slim，pb文件同理
cd producer-service
hz new -module producer-service --idl=../idl/hello.thrift -t=template=slim
hz update --idl=../idl/hello.thrift -t=template=slim
hz update --idl=../idl/auth.thrift -t=template=slim
#hertz client
cd consumer-service
hz client --mod=primus/consumer-service --idl=../idl/hello.thrift --model_dir=model -t=template=slim --client_dir=service/hello
hz client --mod=primus/consumer-service --idl=../idl/auth.thrift --model_dir=model -t=template=slim --client_dir=service/auth

```

### 构建

```
cd primus/producer-service
go build -o producer-service
```

## License Scene

本次开发中，设计了软件的授权场景：比如：我们有一个服务生产者：专门在对客户机器部署软件时，看是否已经被授权。
- 在生产者服务中，在main中启动一个定时器：逻辑是拿到授权License后进行授权认证看是否被授权：状态
- 在该服务中，生成一个授权接口，为提供给其他服务查询授权状态(同时，License生成也在该服务中) 
- 所有服务全局路由时，查询上一步的服务中的授权接口，看状态是否ok


## API requests

### 授权码生成
```shell
##授权码生成
curl -H "Content-Type: application/json" -X POST -d '{"appId": "***.com","issuedTime": 1669082100,"notBefore": 1669082100,"notAfter": 1669116589,"customerInfo": "***公司","authorization": "all,training,inference","machineCodes": ["JL32YL2"],"nodeNum":1}' "http://localhost:2999/v1/createLicence"
```

### 特别鸣谢

本框架基于字节开源的hertz、kitex框架进行扩展，感谢字节机构大力支持：https://github.com/cloudwego。