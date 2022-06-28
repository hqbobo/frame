package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hqbobo/frame/common/conf"
	"github.com/hqbobo/frame/common/dcache"
	"github.com/hqbobo/frame/common/errors"
	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/nsq"
	"github.com/hqbobo/frame/common/queue"
	"github.com/hqbobo/frame/db"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"

	//初始化
	_ "github.com/micro/go-micro/v2/broker/nats"
	//tcpTrans "github.com/micro/go-plugins/transport/tcp"
	"github.com/micro/go-micro/v2/client"
	cgrpc "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	sgrpc "github.com/micro/go-micro/v2/server/grpc"

	"github.com/micro/go-micro/v2/sync/leader"
	"github.com/micro/go-micro/v2/sync/lock"

	//初始化
	_ "github.com/micro/go-plugins/broker/kafka/v2"
	//初始化
	_ "github.com/micro/go-plugins/broker/nsq/v2"
	_ "github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/micro/go-plugins/registry/consul/v2"
	leaderconsul "github.com/micro/go-plugins/sync/leader/consul/v2"
	dlock "github.com/micro/go-plugins/sync/lock/consul/v2"
	rlock "github.com/micro/go-plugins/sync/lock/redis/v2"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	goTracing "github.com/opentracing/opentracing-go"
)

//Dlock 分布式锁
type Dlock = lock.Lock

// Service 封装的Service
type Service struct {
	service micro.Service
	l       lock.Lock
	config  *conf.GlobalConfig
	lead    leader.Leader
	q       *queue.Queue
	nsq     *nsq.Nsq
}

// GetLock 获取
func (svc *Service) GetLock() Dlock {
	return svc.l
}

// Name 封装的Service Name
func (svc *Service) Name() string { return svc.service.Name() }

// Client 封装的Service Client
func (svc *Service) Client() client.Client { return svc.service.Client() }

// Server Server
func (svc *Service) Server() server.Server { return svc.service.Server() }

// Publish 用于发送带回执消息 保证消息送达一个接收端
func (svc *Service) Publish(topic string, data []byte) error {
	if svc.config.Broker.Type == "nsq" {
		return svc.nsq.Publish(topic, data)
	}
	return svc.q.Publish(topic, data)
}

// PublishDelay 用于发送延迟ms带回执消息 保证消息送达一个接收端
func (svc *Service) PublishDelay(topic string, data []byte, ms int) error {
	if svc.config.Broker.Type == "nsq" {
		return svc.nsq.DeferredPublish(topic, time.Duration(ms)*time.Second, data)
	}
	return svc.q.PublishDelay(topic, data, ms)
}

// Consume 接收消息能够消费普通消息和延迟消息
// f func(data []byte, retrans bool) bool
// retrans 是否为重传消息
// 当返回失败的时候 会触发延迟重传
func (svc *Service) Consume(topic, channel string, f func(data []byte, retrans int) bool) error {
	if svc.config.Broker.Type == "nsq" {
		return svc.nsq.Consumer(topic, channel, func(msg []byte) {
			f(msg, 0)
		})
	}
	return svc.q.Consume(topic, f)
}

// Run 开始运行阻塞模式
func (svc *Service) Run(leaderTasks ...func()) error {
	var err error
	if svc.lead != nil && len(leaderTasks) > 0 {
		var e leader.Elected
		var end, done chan bool
		end = make(chan bool, 0)
		done = make(chan bool, 0)
		go func() {
			e, err = svc.lead.Elect("id")
			// handle err
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			// operate while leader
			revoked := e.Revoked()
			for {
				select {
				case <-revoked:
					// re-elect
					svc.lead.Elect("id")
				case <-end:
					e.Resign()
					close(done)
					return
				default:
					// leader operation
					for _, task := range leaderTasks {
						task()
					}
				}
			}
		}()
		err = svc.service.Run()
		close(end)
		// resign leadership
		<-done
	} else {
		err = svc.service.Run()
	}
	return err
}

// Get 获取micro
func (svc *Service) Get() micro.Service { return svc.service }

type svrLogger struct {
}

func (svc svrLogger) Print(v ...interface{}) {
	log.Debugln("%s", fmt.Sprint(v...))
}

//BrokerEvent BrokerEvent定义
type BrokerEvent = broker.Event

func buildAddr(s []string) string {
	var out string
	for i := 0; i < len(s); i++ {
		var addr, port string
		r := strings.Split(s[i], ":")
		addr = r[0]
		port = r[1]
		if net.ParseIP(addr) == nil {
			addrs, err := net.LookupHost(addr)
			if err != nil {
				log.Errorln("%v", err)
			}
			log.Infof("%s -ip地址格式不正确解析dns ", addr, addrs)
			addr = addrs[0]
		}

		if i == 0 {
			out = addr + ":" + port
		} else {
			out += "," + addr + ":" + port
		}
	}
	log.Infoln(out)
	return out
}

// NewService 新建服务避免并发使用
func NewService(name string, config *conf.GlobalConfig) *Service {
	var err error
	var opts []micro.Option
	svc := new(Service)
	svc.config = config
	errors.SetID(name)
	switch config.Dlock.Type {
	case "consul":
		svc.l = dlock.NewLock(lock.Nodes(
			config.Dlock.Addr,
		))

	case "redis":
		svc.l = rlock.NewLock(lock.Nodes(
			config.Dlock.Addr,
		))
		//), lock.TTL(5 *time.Second), lock.Wait(5 *time.Second))
	default:
		log.Infoln("使用默认Dlock")
	}

	switch config.Registry.Type {
	case "consul":
		os.Setenv("MICRO_REGISTRY", "consul")
		os.Setenv("MICRO_REGISTRY_ADDRESS", buildAddr(config.Registry.Addrs))
		svc.lead = leaderconsul.NewLeader(
			leader.Group(name),
			leader.Nodes(config.Registry.Addrs[0]),
		)
	default:
		//reg = registry.DefaultRegistry
		log.Infoln("使用默认注册中心")
	}

	switch config.Broker.Type {
	case "kafka":
		os.Setenv("MICRO_BROKER", "kafka")
		os.Setenv("MICRO_BROKER_ADDRESS", buildAddr(config.Broker.Addrs))
		log.Infoln("消息队列使用kafka")
	case "nats":
		log.Infof("消息队列使用nats")
		os.Setenv("MICRO_BROKER", "nats")
		os.Setenv("MICRO_BROKER_ADDRESS", buildAddr(config.Broker.Addrs))
	case "rabbitmq":
		log.Infof("消息队列使用rabbitmq")
		os.Setenv("MICRO_BROKER", "rabbitmq")
		os.Setenv("MICRO_BROKER_ADDRESS", config.Broker.Addrs[0])
	case "nsq":
		log.Infof("消息队列使用nsq")
		os.Setenv("MICRO_BROKER", "nsq")
		os.Setenv("MICRO_BROKER_ADDRESS", config.Broker.Addrs[0])
	default:
		log.Infoln("使用默认注册中心")

	}
	cmd.Init()
	if err := broker.Init(); err != nil {
		log.Errorf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Errorf("Broker Connect error: %v", err)
	}
	os.Setenv("MICRO_CLIENT_POOL_SIZE", "200")

	//设置客户端链接数
	opts = append(opts, micro.Name(name))
	opts = append(opts, micro.Address(conf.CFG().Host))
	opts = append(opts, micro.RegisterTTL(time.Second*6))
	opts = append(opts, micro.RegisterInterval(time.Second*2))
	opts = append(opts, micro.Server(sgrpc.NewServer(server.Name(name),
		sgrpc.MaxMsgSize(1024*1024*100))))
	opts = append(opts, micro.Client(cgrpc.NewClient(
		client.RequestTimeout(time.Second*60),
		client.PoolSize(200),
		cgrpc.MaxSendMsgSize(1024*1024*100),
		cgrpc.MaxRecvMsgSize(1024*1024*100))))
	opts = append(opts, micro.WrapClient(newWrapper))
	if config.Tracelink != "" {
		//log.Debugf("启动链路追踪:%s", config.Tracelink)
		goTracing.SetGlobalTracer(newTracer(name, config.Tracelink))
		opts = append(opts, micro.WrapHandler(wrapperTrace.NewHandlerWrapper(goTracing.GlobalTracer())))
		opts = append(opts, micro.WrapClient(wrapperTrace.NewClientWrapper(goTracing.GlobalTracer())))
	}
	//opts = append(opts, micro.Transport(tcpTrans.NewTransport()))
	svc.service = micro.NewService(opts...)
	svc.service.Init()
	dcache.Init(conf.CFG().Cache.Addrs, conf.CFG().Cache.Pass)
	db.InitDal(config)
	if conf.CFG().Rabbitmq != "" {
		svc.q, err = queue.NewQueue(conf.CFG().Rabbitmq)
		if err != nil {
			log.Error("Queue启动失败", err)
		}
	}
	if config.Broker.Type == "nsq" {
		svc.nsq = nsq.NewNsq(config.Broker.Addrs)
	}
	return svc
}

// A Wrapper that creates a Datacenter Selector Option
type dcWrapper struct {
	client.Client
}

func (dc *dcWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	balancekey := md["Balancekey"]
	key, _ := strconv.Atoi(balancekey)
	filter := func(services []*registry.Service) []*registry.Service {
		for _, service := range services {
			total := len(service.Nodes)
			var nodes []*registry.Node
			for i, node := range service.Nodes {
				if key == 0 {
					nodes = append(nodes, node)
					continue
				}
				if i == key%total {
					nodes = append(nodes, node)
				}
			}
			service.Nodes = nodes
			// log.Debugf("callService:%s name:%s Key:%d nodes:%+v", req.Service, service.Name, key, service.Nodes)
		}
		return services
	}
	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(filter),
	))
	return dc.Client.Call(ctx, req, rsp, callOptions...)
}

func newWrapper(c client.Client) client.Client {
	return &dcWrapper{c}
}

// NewCliService 新建client
func NewCliService(name string, config *conf.GlobalConfig) *Service {
	var reg registry.Registry
	svc := new(Service)
	switch config.Registry.Type {
	case "consul":
		reg = consul.NewRegistry(func(op *registry.Options) {
			op.Addrs = config.Registry.Addrs
		})
	default:
		reg = registry.DefaultRegistry
		log.Infoln("没有配置注册中心")
	}
	if goTracing.IsGlobalTracerRegistered() {
		svc.service = micro.NewService(micro.Registry(reg),
			micro.Name(name),
			micro.WrapClient(wrapperTrace.NewClientWrapper(goTracing.GlobalTracer())),
			micro.WrapClient(newWrapper),
		)
	} else {
		svc.service = micro.NewService(micro.Registry(reg),
			micro.Name(name),
		)
	}
	svc.service.Init()
	return svc
}

//SetClientBalance 设置微服务调用负载均衡key
func SetClientBalance(ctx context.Context, key int64) context.Context {
	return metadata.NewContext(ctx, map[string]string{
		"Balancekey": fmt.Sprintf("%d", key),
	})
}

//BrokerMessage 消息队列定义
type BrokerMessage struct {
	Header map[string]string
	Body   []byte
}

func _brokerMsgParse(msg *BrokerMessage) *broker.Message {
	out := new(broker.Message)
	out.Header = msg.Header
	out.Body = msg.Body
	return out
}

//Deprecated Publish 消息队列发送用户全局广播
func Publish(topic string, msg *BrokerMessage, opts ...broker.PublishOption) error {
	return broker.Publish(topic, _brokerMsgParse(msg), opts...)
}

//BrokerHanlde 消息队列定义
type BrokerHanlde = broker.Handler

//Deprecated Subscribe 消息队列接收
func Subscribe(topic string, handler BrokerHanlde, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return broker.Subscribe(topic, handler, opts...)
}

// RegisterEvent is syntactic sugar for registering a subscriber
func RegisterEvent(topic string, s server.Server, h interface{}, opts ...server.SubscriberOption) error {
	return s.Subscribe(s.NewSubscriber(topic, h, opts...))
}

// Client 转换
type Client = client.Client

// Event 转换
type Event = micro.Event

// NewEvent Event事件
func NewEvent(topic string, cli client.Client) Event {
	return micro.NewEvent(topic, cli)
}
