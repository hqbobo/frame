package datakeeper

import (
	"sync"
	"time"
)

// DataSource 数据源的实现
type DataSource interface {
	Load() map[interface{}]interface{}
	Update(key, val interface{})
	Delete(key interface{})
}

// DataKeeper 数据维护管理模块
type DataKeeper struct {
	opt    options
	pool   *sync.Map
	ds     DataSource
	close  chan bool
	update *sync.Map
}

// Store 保存数据
func (dk *DataKeeper) Store(key, data interface{}) {
	dk.pool.Store(key, data)
	//是否开启定时更新
	if !dk.opt.upf {
		dk.ds.Update(key, data)
	} else {
		dk.update.Store(key, data)
	}
}

// Delete 删除数据
func (dk *DataKeeper) Delete(key interface{}) {
	dk.pool.Delete(key)
	dk.ds.Delete(key)
}

// Reload 重载数据
func (dk *DataKeeper) Reload() { dk.reload() }

// Load 获取数据
func (dk *DataKeeper) Load(key interface{}) (interface{}, bool) {
	return dk.pool.Load(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (dk *DataKeeper) Range(f func(key, value interface{}) bool) {
	dk.pool.Range(f)
}

func (dk *DataKeeper) Exit() {
	close(dk.close)
}

func NewDataKeeper(ds DataSource, opts ...Option) *DataKeeper {
	dr := new(DataKeeper)
	dr.close = make(chan bool)
	dr.update = new(sync.Map)
	dr.opt = options{
		of:  false,
		tf:  false,
		upf: false,
	}
	for _, opt := range opts {
		opt.apply(&dr.opt)
	}
	dr.ds = ds
	if dr.opt.tf {
		go dr.tf()
	}
	if dr.opt.of {
		go dr.of()
	}
	if dr.opt.upf {
		go dr.uf()
	}
	dr.reload()
	return dr
}

func (dk *DataKeeper) reload() {
	list := dk.ds.Load()
	pool := new(sync.Map)
	for k, v := range list {
		pool.Store(k, v)
	}
	dk.pool = pool
}

//根据指定时间计算下一个时间差距
func getNextTime(h, m, s int) time.Duration {
	now := time.Now()
	//算出当天运行时间
	target := now.Add(-time.Hour*time.Duration(now.Hour()) -
		time.Minute*time.Duration(now.Minute()) -
		time.Second*time.Duration(now.Second()) -
		time.Nanosecond*time.Duration(now.Nanosecond()))
	target = target.Add(time.Hour*time.Duration(h) +
		time.Minute*time.Duration(m) +
		time.Second*time.Duration(s))
	//今天运行
	if target.After(now) {
		return target.Sub(now)
	}
	return target.AddDate(0, 0, 1).Sub(now)
}

func (dk *DataKeeper) tf() {
	timer := time.NewTimer(dk.opt.timer)
	for {
		select {
		case <-timer.C:
			dk.reload()
			timer.Reset(dk.opt.timer)
		case <-dk.close:
		}
	}
}

func (dk *DataKeeper) of() {
	var timer *time.Timer
	for {
		dur := getNextTime(dk.opt.hour, dk.opt.min, dk.opt.sec)
		timer = time.NewTimer(dur)
		select {
		case <-timer.C:
			dk.reload()
		case <-dk.close:
		}
	}
}

func (dk *DataKeeper) uf() {
	timer := time.NewTimer(dk.opt.uptimer)
	for {
		select {
		case <-timer.C:
			dk.update.Range(func(key, value interface{}) bool {
				dk.ds.Update(key, value)
				dk.update.Delete(key)
				return true
			})
			timer.Reset(dk.opt.uptimer)
		case <-dk.close:
		}
	}
}
