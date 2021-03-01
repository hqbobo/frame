#定义过滤的目录
FILTER = $(BINROOT) $(BINROOT)/allconfig $(BINROOT)/emqttd
#遍历所有目录

ifdef SVC
ALLDIRS = $(shell find $(BINROOT)  -maxdepth 1  -type d |grep  "\/$(SVC)" |grep  -w "$(SVC)")
else
ALLDIRS = $(shell find $(BINROOT)  -maxdepth 1  -type d)
endif
#过滤目录
DIRS = $(filter-out $(FILTER),$(ALLDIRS))

.PHONY:bins deps test run docker proto

bins:proto
	$(Q) for i in $(DIRS); \
		do \
		$(PYellow) "进入目录编译:$$i"$(PClear);\
		$(CD) $$i;go build;file $$i;\
		if [ $$? != 0 ]; then exit 1; fi;\
		$(PGreen) "编译完成离开目录:$$i"$(PClear);\
		done\

deps:
	$(Q) for i in $(DIRS); \
		do \
		$(PYellow) "进行golint检查:$$i========"$(PClear);$(PRed);\
		$(CD) $$i;golint ./... | grep -v vendor > tmp;\
		if [ -s ./tmp ]; then golint ./...; $(RM) tmp;$(PYellow) $(PClear);$(PRW) "golint 检查失败" $(PClear); exit 1; fi; $(RM) tmp;\
		$(PGreen) "golint检查完成:$$i========"$(PClear);\
		done\

proto:
	$(Q) for i in $(DIRS); \
		do \
		$(PYellow) "进行proto相关编译:$$i========"$(PClear);\
        $(CD) $$i/proto;protoc --go_out=. --ext_out=. *.proto;$(EMTPYCMD);\
        $(PGreen) "完成proto相关编译完成相关代码在:$$i/proto目录下========"$(PClear);\
		$(PYellow) "进行proto编译:$$i========"$(PClear);$(PRed);\
		protoc --go_out=. --micro_out=. *.proto;$(EMTPYCMD);$(V2CMD);\
		$(PGreen) "proto编译完成:$$i========"$(PClear);\
		done\

test:
	$(Q) for i in $(DIRS); \
		do \
		$(PYellow) "进入目录:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"services/\");print arr[2]}"`;\
		$(PYellow) "========[$$name]启动服务========"$(PClear);export RUNMODE="test";\
		$(CD) $$i;./$$name &\
		$(PYellow) "========[$$name]启动测试========"$(PClear); sleep 1;\
		cd tests/; go test -v;\
		if [ $$? != 0 ]; then $(PRed) "========[$$name]测试失败========"$(PClear); exit 1; fi;\
		$(PGreen) "========[$$name]测试完成========"$(PClear);\
		if [ $(OS) == "Darwin" ]; then killall $$name ; else  kill -9 `pidof $$name`; fi;\
		$(PGreen) "离开目录:$$i"$(PClear);\
		done\

run:bins
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "进入目录:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"services/\");print arr[2]}"`;\
		$(CD) $$i;./$$name ;\
		if [ $$? != 0 ]; then exit 1; fi;\
		$(PGreen) "离开目录:$$i"$(PClear);\
		done\

clean:
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "清理:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"services/\");print arr[2]}"`;\
		$(CD) $$i;rm -rf $$name ;$(RM) ./proto/cli;$(RM) ./proto/server;$(RM) ./proto/*.go;\
		if [ $$? != 0 ]; then exit 1; fi;\
		$(PGreen) "离开:$$i"$(PClear);\
		done\

docker:
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "编译镜像:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"services/\");print arr[2]}"`;\
		echo name $$name;\
		$(CD) $$i;\
		$(CP) $(SCRITP)/timezone Shanghai;\
		$(CP) $(CONFROOT)/$(RUNMODE).json config.json;\
		echo docker build -t $(HARBOR)/$$name:$(VER) .;\
		docker build -t $(HARBOR)/$$name:$(VER) .;\
		docker push $(HARBOR)/$$name:$(VER);\
		docker rmi $(HARBOR)/$$name:$(VER) -f;\
		rm -rf config.json Shanghai $$name;\
		$(PGreen) "完成编译镜像:$$name"$(PClear);\
		done\