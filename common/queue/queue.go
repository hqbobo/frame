package queue

import (
	"fmt"
	//"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/log"
	"sync"

	"github.com/streadway/amqp"
)

type topic struct {
	name string
	ch   *amqp.Channel
	que  amqp.Queue
	conn *amqp.Connection
}

type topics struct {
	list []*topic
	lock sync.Mutex
}

func (tops *topics) GetChannel(name string, conn *amqp.Connection) (*amqp.Channel, error) {
	tops.lock.Lock()
	defer tops.lock.Unlock()
	for _, v := range tops.list {
		if v.name == name {
			return v.ch, nil
		}
	}
	top, err := tops.addTopic(name, conn)
	return top.ch, err
}

const (
	queuePrefix = "queuePreix"
)

func (tops *topics) addTopic(name string, conn *amqp.Connection) (ntopic *topic, err error) {
	ntopic = new(topic)
	ntopic.name = name
	ntopic.ch, err = conn.Channel()
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	// 声明一个主要使用的 exchange
	err = ntopic.ch.ExchangeDeclare(
		name,     // name
		"fanout", // type
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	ntopic.que, err = ntopic.ch.QueueDeclare(
		queuePrefix+name, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	err = ntopic.ch.QueueBind(
		ntopic.que.Name, // queue name, 这里指的是 test_logs
		"",              // routing key
		name,            // exchange
		false,
		nil)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	/**
	 * 注意,这里是重点!!!!!
	 * 声明一个延时队列, ß我们的延时消息就是要发送到这里
	 */
	ntopic.ch.QueueDeclare(
		queuePrefix+name+"_delay", // name
		false,                     // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		amqp.Table{
			// 当消息过期时把消息发送到 logs 这个 exchange
			"x-dead-letter-exchange": name,
		}, // arguments
	)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	tops.list = append(tops.list, ntopic)
	return ntopic, nil
}

// NewQueue 申请一个新的队列
func NewQueue(addr string) (q *Queue, err error) {
	//var err error
	q = new(Queue)
	q.conn, err = amqp.Dial(addr)
	return q, err
}

// Queue 消息队列
// 持久化队列保证消息到达
type Queue struct {
	conn *amqp.Connection
	tops topics
}

// Publish 发送消息
func (q *Queue) Publish(topic string, data []byte) error {
	table := make(map[string]interface{}, 0)
	table["delay"] = false
	table["retrans"] = 0
	return q.publish(topic, data, 0, table)
}

// PublishDelay 发送消息 time 毫秒
func (q *Queue) PublishDelay(topic string, data []byte, time int) error {
	table := make(map[string]interface{}, 0)
	table["delay"] = false
	table["retrans"] = 0
	return q.publish(topic, data, time, table)
}

// Send 发送消息 time 毫秒
func (q *Queue) publish(topic string, data []byte, time int, table amqp.Table) error {
	table["delaytime"] = time
	chn, err := q.tops.GetChannel(topic, q.conn)
	if err != nil {
		return err
	}
	return chn.Publish(
		"",                         // exchange
		queuePrefix+topic+"_delay", // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			Headers:     table,
			ContentType: "text/plain",
			Body:        data,
			Expiration:  fmt.Sprint(time), // 设置五秒的过期时间
		})
}

// Consume 接收消息
// f func(data []byte, retrans bool) bool
// retrans 是否为重传消息
// 当返回失败的时候 会触发延迟重传
func (q *Queue) Consume(topic string, f func(data []byte, retrans int) bool) error {
	chn, err := q.tops.GetChannel(topic, q.conn)
	if err != nil {
		return err
	}
	msgs, err := chn.Consume(
		queuePrefix+topic, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return err
	}
	go func(t string) {
		for d := range msgs {
			//log.Debug(d.Headers)
			//判断执行返回
			retrans := int(d.Headers["retrans"].(int32))
			if !f(d.Body, retrans) {
				d.Headers["retrans"] = retrans + 1
				delay := int(d.Headers["delaytime"].(int32))
				if delay == 0 {
					delay = 1000
				} else {
					delay += 1000
				}
				q.publish(t, d.Body, delay, d.Headers)
			}
		}
	}(topic)
	return nil
}
