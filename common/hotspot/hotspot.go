package hotspot

import (
	"math/rand"
	"time"
)

var hotspot *Hotspot

//GetEngine 获取对象实例
func GetEngine() *Hotspot {
	return hotspot
}

//Hotspot 热点缓存
type Hotspot struct {
	mem *Mem
}

//Get 获取缓存
func (h *Hotspot) Get(key string) (string, bool) {
	return h.mem.Get(key)
}

//Set 设置缓存
func (h *Hotspot) Set(key, value string) {
	h.mem.Set(key, string(value), int64(intRand(5, 10)))
}

func init() {
	hotspot = new(Hotspot)
	hotspot.mem = NewMem()
	rand.Seed(time.Now().UnixNano())
}

//intRand rand一个范围值
func intRand(min, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min) + min
}
