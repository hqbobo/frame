#定义过滤的目录
FILTER = $(APIROOT)
#遍历所有目录

ifdef SVC
ALLDIRS = $(shell find $(APIROOT)  -maxdepth 1  -type d |grep  "\/$(SVC)" |grep  -w "$(SVC)")
else
ALLDIRS = $(shell find $(APIROOT)  -maxdepth 1  -type d)
endif
#过滤目录
DIRS = $(filter-out $(FILTER),$(ALLDIRS))

.PHONY:bins deps test run docker

bins:
	$(Q) for i in $(DIRS); \
		do \
		$(PYellow) "进入目录编译:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"api/\");print arr[2]}"`;\
		$(PYellow) "进入目录编译:$$i"$(PClear);\
		$(CD) $$i;$(RM) proto;$(CP) ../proto proto;$(RM) response;$(CP) ../../define/response response;swag init;go build ;$(RM) proto;$(RM) response;\
		if [ $$? != 0 ]; then exit 1; fi;\
		$(PGreen) "编译完成离开目录:$$i"$(PClear);\
		done\

deps:
	$(Q) for i in $(DIRS); \
		do \
		$(PYellow) "进行golint检查:$$i"$(PClear);\
		$(CD) $$i;$(PRed);golint ./... >tmp;\
		if [ -s ./tmp ]; then golint ./...; $(RM) tmp;$(PYellow) $(PClear);$(PRW) "golint 检查失败" $(PClear); exit 1; fi; $(RM) tmp;\
		$(PGreen) "golint检查完成:$$i"$(PClear);\
		done\

run:bins
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "进入目录:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"api/\");print arr[2]}"`;\
		$(CD) $$i; ./$$name ;\
		if [ $$? != 0 ]; then exit 1; fi;\
		$(PGreen) "离开目录:$$i"$(PClear);\
		done\

test:
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "进入目录:$$i"$(PClear);\
		$(CD) $$i;go test;\
		$(PGreen) "离开目录:$$i"$(PClear);\
		done\

clean:
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "清理:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"api/\");print arr[2]}"`;\
		$(CD) $$i;$(RM) $$name;$(RM) Shanghai;\
		if [ $$? != 0 ]; then exit 1; fi;\
		$(PGreen) "离开:$$i"$(PClear);\
		done\

docker:bins
	$(Q) for i in $(DIRS); \
		do \
		$(PGreen) "编译镜像:$$i"$(PClear);\
		name=`echo | awk "{split(\"$$i\", arr, \"api/\");print arr[2]}"`;\
		$(CD) $$i;\
		$(CP) $(SCRITP)/timezone Shanghai;\
		echo docker build -t $(HARBOR)/$$name:$(VER) .;\
		docker build -t $(HARBOR)/$$name:$(VER) .;\
		docker push $(HARBOR)/$$name:$(VER);\
		docker rmi $(HARBOR)/$$name:$(VER) -f;\
		$(RM) Shanghai $$name;\
		$(PGreen) "完成编译镜像:$$name"$(PClear);\
		done\