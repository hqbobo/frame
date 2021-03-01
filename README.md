# 框架描述
本框架基于go-micro做二次封装,将各种中间件开发进行了统一封装

## 目录描述

- [api](./api) - API相关服务目录
- [common](./common) - 框架通用封装
- [db](./db) - 数据库相关
- [define](./define) - 通用定义
- [script](./script) - 脚本相关
- [services](./services) - 微服务目录


## 中间件选型
    
- ***NOSQL***
    - redis
    
- ***消息队列*** 
    - kakfa
    - nats
    - [rabbitmq] - 管理后台
- ***数据库*** 
    - mysql
    
- ***注册中心***
    - consul

- ***链路追踪***
    - [Jaeger] - 链路追踪页面 

## 周边环境地址
