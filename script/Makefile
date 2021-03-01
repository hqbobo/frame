#基础定义
include mk.def

#跨平台定义
ifdef ARCH
include mk.$(ARCH)
else
include mk.linux
endif

ifdef SVC
export APP=$(SVC)
endif

ifdef NOCOLOR
include include/print.def
else
include include/colorprint.def
endif


.PHONY:app console deps clean run services api run runapi docker proto

all:deps api services test

runapi:
	$(Q) $(PYellow) "-------------开始启动网关服务--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(APIROOT)" $(PClear);
	$(Q) $(CD) $(APIROOT);make run;
	$(Q) $(PYellow) "-------------完成启动网关服务--------------"$(PClear);

api:
	$(Q) $(PYellow) "-------------开始编译网关--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(APIROOT)" $(PClear);
	$(Q) $(CD) $(APIROOT);make bins;
	$(Q) $(PYellow) "-------------完成编译网关--------------"$(PClear);

services:
	$(Q) $(PYellow) "-------------开始编译微服务--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(BINROOT)" $(PClear);
	$(Q) $(CD) $(BINROOT);make bins;
	$(Q) $(PGreen) "-------------完成编译微服务--------------"$(PClear);

proto:
	$(Q) $(PYellow) "-------------开始准备proto--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(TOOL) 编译proto工具" $(PClear);
	$(Q) $(CD) $(TOOL);go build -o protoc-gen-ext;
	$(Q) $(PGreen) "进入目录:$(BINROOT)" $(PClear);
	$(CD) $(BINROOT);make proto;
	$(Q) $(PGreen) "-------------完成准备proto--------------"$(PClear);

deps:
	$(Q) $(PYellow) "-------------微服务依赖处理--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(BINROOT)" $(PClear);
	$(Q) $(CD) $(BINROOT);make deps;
	$(Q) $(PYellow) "------------完成微服务依赖处理--------------"$(PClear);

runconsole:
	$(CWD)/$(BIN)/$(CONSOLE) &

run:
	$(Q) $(PYellow) "-------------开始启动微服务--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(BINROOT)" $(PClear);
	$(Q) $(CD) $(BINROOT);make run;
	$(Q) $(PYellow) "-------------完成启动微服务--------------"$(PClear);

test:
	$(Q) $(PYellow) "-------------微服务GO test 检查--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(BINROOT)" $(PClear);
	$(Q) $(CD) $(BINROOT);make test;
	$(Q) $(PYellow) "------------完成微服务GO test 检查--------------"$(PClear);

clean:
	$(Q) $(CD) $(BINROOT);make clean;
	$(Q) $(CD) $(APIROOT);make clean;
	$(Q) $(CD) $(TOOL);$(RM) protoc-gen-ext
	$(Q) $(RM) $(BIN);

docker:services
	$(Q) $(PYellow) "-------------创建镜像--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(BINROOT)" $(PClear);	
	$(Q) $(CD) $(BINROOT);make docker;
	$(Q) $(PYellow) "------------创建镜像完成------------"$(PClear);

dockerapi:
	$(Q) $(PYellow) "-------------创建API镜像--------------" $(PClear);
	$(Q) $(PGreen) "进入目录:$(APIROOT)" $(PClear);	
	$(Q) $(RM) $(APIROOT)/proto; mkdir $(APIROOT)/proto;find ../ |grep pb | xargs -n1 -I {} $(CP) -f {} $(APIROOT)/proto
	$(Q) $(CD) $(APIROOT);make docker;
	$(Q) $(RM) $(APIROOT)/proto
	$(Q) $(PYellow) "------------创建镜像完成------------"$(PClear);
