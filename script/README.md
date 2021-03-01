# 脚本
本项目基于Makefile 进行工程编译

## 编译参数
- `all` - 编译所有服务
- `test` - 进行go test必须先编译
- `dep` - 编译准备,目前做golint检查
- `run` - 启动服务 `serivces下`
- `runapi` - 启动api服务 `api目录下`
- `clean` - 清理项目
- `SVC` - 指定编译的服务 `SVC=user` 本选项可以可以配合上面选项使用
- `ARCH` - 指定输出构架 `ARCH=linux`
- `docker` - 输出服务docker镜像
- `dockerapi` - 输出api服务docker镜像
- `RUN` - 指定输出环境.
- `MOD` - 指定配置文件{develop local}.
- `proto` - 编译api服务之前需要先编译proto.

##使用例子

- `make proto` - 编译proto否则swag编译不过
- `make SVC=test` - 编译test服务,同时会进行`dep,test`操作
- `make test SVC=test` - 运行指定服务的go test
- `make clean SVC=test` - 清理test服务
- `make dockerapi SVC=gateway ARCH=linux` 编译api下gateway镜像


