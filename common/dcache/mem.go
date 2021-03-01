package dcache

import (
	"sync"
	"time"
)

type sobj struct {
	obj string
	ttl int64
}

type MemSession struct {
	cache map[string]sobj
	lock  *sync.RWMutex
}

func newMemSession() *MemSession {
	s := new(MemSession)
	s.cache = make(map[string]sobj, 0)
	s.lock = new(sync.RWMutex)
	return s
}

func (this *MemSession) Get(key string, data *string) bool {
	this.lock.RLock()
	if v, ok := this.cache[key]; ok {
		if time.Now().Unix() < v.ttl {
			this.lock.RUnlock()
			*data = v.obj
			return true
		}
	}
	this.lock.RUnlock()
	return false
}

func (this *MemSession) Set(key string, data string, ttl int) bool {
	this.lock.Lock()
	if _, ok := this.cache[key]; ok {
		delete(this.cache, key)
	}
	o := new(sobj)
	//超时最大为一小时
	if ttl > 60*60 {
		ttl = 60 * 60
	}
	o.ttl = time.Now().Unix() + int64(ttl)
	o.obj = data
	this.cache[key] = *o
	this.lock.Unlock()
	return true
}

func (this *MemSession) Delete(key string) bool {
	this.lock.Lock()
	if _, ok := this.cache[key]; ok {
		delete(this.cache, key)
	}
	this.lock.Unlock()
	return true
}

func (this *MemSession) Check(key string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if v, ok := this.cache[key]; ok {
		if time.Now().Unix() < v.ttl {
			return true
		} else {
			delete(this.cache, key)
			return false
		}
	}
	return false
}
