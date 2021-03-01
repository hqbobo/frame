# db
本目录主要是对数据库操作进行封装

## 介绍
- `common` - 数据库操作相关
    - `impl` - 数据库接口实现
    - `manage` - 数据库连接申请
    
- `model` - 定义数据库结构

## 使用流程
- `添加model` 在model目录下添加数据库表对应的结构
- `添加impl` 在`common/imlp` 目录下实现数据库操作
- `添加manage` 在`manage.go` 内部添加对应的New函数



