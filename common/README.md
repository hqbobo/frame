# common

本目录下提供所有微服务框架下所使用的功能

## service 

微服务的基础模块.

## dcache 

框架下使用的分布式缓存,目前使用redis实现,支持内存同步 (依赖service模块初始化)

## log

框架下使用的日志.

## conf

所有微服务的配置

## encode

`interface` 到 `byte` 的编解码

## transport

对传输层进行封装


## datakeeper 

数据缓存维护中间件

## test 

测试相关的


## limiter 

限流器 配置环境变量 LimiterRate = 3 //每秒3个

