package limiter

import (
	"sync/atomic"
)

//计数器
type counter interface {
	Store(val int)
	Load() int
	Add(delta int)
} 


type atomicCounter struct {
	count int32
}

func (mc *atomicCounter) Store(val int) {
	atomic.StoreInt32(&mc.count, 1)
}

func (mc *atomicCounter) Load() int {
	return int(atomic.LoadInt32(&mc.count))
}

func (mc *atomicCounter) Add(delta int) {
	atomic.AddInt32(&mc.count, int32(delta))
}
