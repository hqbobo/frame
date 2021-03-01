# service

所有微服务相关函数都封装在这里

##  主要函数
``func NewService(name string, config *conf.GlobalConfig) *Service``

本函数主要用来创建微服务需要配合 ``/qygame/common/conf``使用

***内部结合了***

- 注册中心
- 消息队列
- Cache缓存
- Mysql数据库
- 链路监控

所有内部相关模块将会自动启动只要有相关配置

---
``func NewCliService(name string, config *conf.GlobalConfig) *Service``

本函数主要是用来创建一个微服务客户端调用微服务使用
```go
    //创建微服务客户端
    service := service.NewCliService(define.SVC_USER, config)
    //绑定user微服务
    svc := proto.NewUserService(define.SVC_USER, service.Client())
```
---
`func SetClientBalance(ctx context.Context, key int64) context.Context`

为客户端调用设置负载均衡key值 选取策略为 key % (节点数量) 配合`NewCliService`使用

---
``Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error``

发送消息到消息队列,要使用本函数必须优先初始化``NewService``并且在配置内包含相关配置
```go
service.Publish("test111", &broker.Message{
    Header: map[string]string{
        "AAA": "BBBBB",
        "CCCCC": "DDDDDD",
    },
    Body: []byte("消息内容"+time.Now().String()),
}); err != nil {
    log.Fatal("%s",err.Error())
}
```

---
``Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error)``

注册消息队列监听函数,要使用本函数必须优先初始化``NewService``并且在配置内包含相关配置
```go
service.Subscribe("test111", func(p broker.Event) error {
    log.Warn("[sub] received message:%s header %s", string(p.Message().Body), p.Message().Header)
    return nil
})
```

