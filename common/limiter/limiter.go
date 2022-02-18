package limiter

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hqbobo/frame/common/log"
)

type limiterData struct {
	addr    string
	count   counter
	timeout time.Time
}

//Limiter 熔断限流器
type Limiter struct {
	rate int
	maps *sync.Map
}

var gLimiter *Limiter
var defaultRate = 10 //默认每秒10次
var Unlimit = -1

func init() {
	gLimiter = new(Limiter)
	gLimiter.maps = new(sync.Map)
	rate, err := strconv.Atoi(os.Getenv("LimiterRate"))
	if err != nil {
		rate = defaultRate
	}
	log.Infof("限流模块启动:%d/s", rate)
	gLimiter.rate = rate
}

// Try 检查关键字
func (l *Limiter) Try(s string) bool {
	if l.rate == Unlimit {
		return true
	}
	if v, ok := l.maps.Load(s); ok {
		ld := v.(*limiterData)
		//判断是否过期
		if time.Now().After(ld.timeout) {
			//已经过期重制
			ld.count.Store(1)
			ld.timeout = time.Now().Add(time.Second)
		} else {
			ld.count.Add(1)
		}
		// log.Debugln("当前次数:", ld.count.Load())
		if ld.count.Load() > l.rate {
			return false
		}
		return true
	}
	ld := new(limiterData)
	ld.addr = s
	ld.timeout = time.Now().Add(time.Second)
	ld.count = new(atomicCounter)
	l.maps.Store(s, ld)
	return true
}

// Try 检查关键字限流次数
func Try(s string) bool { return gLimiter.Try(s) }
