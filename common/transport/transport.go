package transport

import (
	"time"
)

// Transport is an interface which is used for communication between
// services. It uses connection based socket send/recv semantics and
// has various implementations; http, grpc, quic.
type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}

// Message 消息
type Message struct {
	Header map[string]string
	Body   interface{}
}

// Socket 消息
type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
	Local() string
	Remote() string
}

// Client 客户端
type Client interface {
	Socket
}

// Listener 服务器Listener
type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

// Option Option
type Option func(*Options)

// DialOption DialOption
type DialOption func(*DialOptions)

// ListenOption ListenOption
type ListenOption func(*ListenOptions)

var (
	// DefaultDialTimeout 默认连接超时
	DefaultDialTimeout = time.Second * 5
)

// NewTransport 默认传输层
func NewTransport() Transport {
	return newTCPTransport()
}
