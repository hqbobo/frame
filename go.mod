module github.com/hqbobo/frame

go 1.17

replace (
	github.com/asim/go-micro/v2 => github.com/micro/go-micro/v2 v2.3.0
	github.com/micro/go-micro/v2 => github.com/asim/go-micro/v2 v2.3.0
	k8s.io/api => k8s.io/api v0.0.0-20190708174958-539a33f6e817
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190808125512-07798873deee
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.3.0
	github.com/micro/go-plugins/broker/nsq/v2 v2.3.0
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.3.0
	github.com/micro/go-plugins/registry/consul/v2 v2.3.0
	github.com/micro/go-plugins/sync/leader/consul/v2 v2.3.0
	github.com/micro/go-plugins/sync/lock/consul/v2 v2.3.0
	github.com/micro/go-plugins/sync/lock/redis/v2 v2.3.0
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.3.0
	github.com/micro/protoc-gen-micro v1.0.0
	github.com/nsqio/go-nsq v1.0.8
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.1
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/stretchr/testify v1.6.1
	github.com/ugorji/go v1.1.13 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sys v0.0.0-20210112080510-489259a85091
	golang.org/x/text v0.3.4 // indirect
	gopkg.in/olivere/elastic.v5 v5.0.86
	gopkg.in/sohlich/elogrus.v2 v2.0.2
	xorm.io/core v0.7.3
	xorm.io/xorm v1.0.7
)
