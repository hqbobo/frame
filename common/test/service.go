package test

import (
	"os"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/memory"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func newDefaultRegistry() registry.Registry {
	return registry.DefaultRegistry
}

func newTestMemRegistry() registry.Registry {
	return memory.NewRegistry()
}

func newTestConsulRegistry(addr string) registry.Registry {
	return consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{addr}
	})
}

// LocalService 创建本地测试服务
func LocalService(name string) micro.Service {
	r := newDefaultRegistry()
	os.Setenv("MICRO_CLIENT_POOL_SIZE", "10")
	svc := micro.NewService(
		micro.Name(name),
		micro.Registry(r),
	)
	return svc
}

// ConsulService 创建Consul测试服务
func ConsulService(name, addr string) micro.Service {
	r := newTestConsulRegistry(addr)
	// create service
	os.Setenv("MICRO_CLIENT_POOL_SIZE", "10")
	svc := micro.NewService(
		micro.Name(name),
		micro.Registry(r),
	)
	return svc
}
