# Services

本目录主要是各个微服务主要代码

## 开发基本原则

- 请完善服务的README.md文件介绍服务相关信息
- 请创建服务后将测试统一放到`tests`目录下
- 避免引用外面库,如有必要请通过``/common``做封装
- 请服务内部对必要的耗时处理函数做好链路追踪
- 请准备创建好自己服务的Dockerfile通过`make docker` 进行镜像编译 
- 跨服务调用:
    - `svc := proto.NewTestService(define.SvcTest, svc.Client())` 主要入参是Client()
    - 如果本身是微服务请使用`NewService()`生成的服务来产生Client()
    - 如果本身不是微服务请使用`service.NewCliService(define.SvcTest, conf.CFG()).Client()`来产生client
    
    
- 消息队列:
    - 使用消息队列请在`/define/broker_msg.go`下定义队列topic并且使用`prefixMsq`前缀
    - 使用消息队列请在`/define/broker_msg.go`下定义结构
    - 定义任意消息队列结构的时候请包含`encode` 结构体做统一的编解码封装
    
    
- Event:
    - Event会广播到所有注册节点, 相对于消息队列比较轻量级
    - 使用Event请在`/define/broker_msg.go`下定义队列topic并且使用`prefixEvent`前缀
    - 定义任意Event结构的时候请包含`encode` 结构体做统一的编解码封装


- 统一错误`AppError`:
    - 所有错误ID统一定义到`/define/errors/def.go`下面,同时定义统一错误描述
    - 微服务互相调用通过 `ParseAppError(error)` 解析错误到 `AppError`

## 服务介绍

### test
 测试微服务这是一个最基础的微服务例子,请优先参考本例子再进行开发
 - 包含内容:
    - [cli](./test/cli.go) - 客户端案例
    - [db](./test/db.go) - 数据库例子
    - [event](./test/event.go) - event例子
    - [lock](./test/lock.go) - 分布式锁
    - [mq](./test/mq.go) - 消息队列例子
    - [transport](./test/transport.go) - transport例子
---



