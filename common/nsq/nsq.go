package nsq

import (
	"sync"
	"time"

	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/rand"
	"github.com/nsqio/go-nsq"
)

//Logger 日志
type Logger struct {
}

//Output 输出
func (l *Logger) Output(calldepth int, s string) error {
	log.Traceln(s)
	return nil
}

//Nsq 生成者
type Nsq struct {
	pos      int
	count    int
	address  []string
	producer *nsq.Producer
	consumer []*nsq.Consumer
	lock     sync.Mutex
}

//NewNsq 新建一个Nsq
func NewNsq(address []string) *Nsq {
	n := new(Nsq)
	n.pos = 0
	n.count = len(address)
	n.address = address
	index := rand.IntRand(0, n.count)
	addr := n.address[index]
	if err := n.initProducer(addr); err != nil {
		log.Fatalf("连接%v失败: %s", n.address, err.Error())
	}
	return n
}

func (n *Nsq) initProducer(addr string) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, config)
	if err != nil {
		return err
	}
	producer.SetLogger(new(Logger), nsq.LogLevelInfo)
	err = producer.Ping()
	if err != nil {
		producer.Stop()
		return err
	}
	n.producer = producer
	return nil
}

func (n *Nsq) setConsumer(consumer *nsq.Consumer) {
	n.lock.Lock()
	n.consumer = append(n.consumer, consumer)
	n.lock.Unlock()
}

func (n *Nsq) getProducer() *nsq.Producer {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.producer
}

//Publish 发送消息
func (n *Nsq) Publish(topic string, msg []byte) error {
	if err := n.getProducer().Ping(); err != nil {
		n.getProducer().Stop()
		for _, addr := range n.address {
			//此处全部断线重连一下
			if err := n.initProducer(addr); err == nil {
				break
			}
		}
	}
	return n.getProducer().Publish(topic, msg)
}

//DeferredPublish 延迟消息
func (n *Nsq) DeferredPublish(topic string, delay time.Duration, msg []byte) error {
	if err := n.getProducer().Ping(); err != nil {
		n.getProducer().Stop()
		for _, addr := range n.address {
			//此处全部断线重连一下
			if err := n.initProducer(addr); err == nil {
				break
			}
		}
	}
	return n.getProducer().DeferredPublish(topic, delay, msg)
}

//Consumer 创建消费者
func (n *Nsq) Consumer(topic, channel string, f func(msg []byte)) error {
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second
	config.MaxInFlight = len(n.address)
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Errorf("创建消费者失败: %s", err.Error())
		return err
	}
	consumer.SetLogger(new(Logger), nsq.LogLevelInfo)
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		f(message.Body)
		message.Finish()
		return nil
	}))
	//建立多个nsqd连接
	if err := consumer.ConnectToNSQDs(n.address); err != nil {
		log.Errorf("连接%v失败: %s", n.address, err.Error())
		return err
	}
	n.setConsumer(consumer)
	<-consumer.StopChan
	stats := consumer.Stats()
	log.Tracef("message received %d, finished %d, requeued:%s, connections:%s", stats.MessagesReceived, stats.MessagesFinished, stats.MessagesRequeued, stats.Connections)
	return nil
}

//Close 关闭
func (n *Nsq) Close() {
	for _, consumer := range n.consumer {
		consumer.Stop()
	}
	n.producer.Stop()
}
