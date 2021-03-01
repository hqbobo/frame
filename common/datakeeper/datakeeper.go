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
func (this *DataKeeper) Store(key, data interface{}) {
	this.pool.Store(key, data)
	//是否开启定时更新
	if !this.opt.upf {
		this.ds.Update(key, data)
	} else {
		this.update.Store(key, data)
	}
}

// Delete 删除数据
func (this *DataKeeper) Delete(key interface{}) {
	this.pool.Delete(key)
	this.ds.Delete(key)
}

// Reload 重载数据
func (this *DataKeeper) Reload() { this.reload() }

// Load 获取数据
func (this *DataKeeper) Load(key interface{}) (interface{}, bool) {
	return this.pool.Load(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (this *DataKeeper) Range(f func(key, value interface{}) bool) {
	this.pool.Range(f)
}

func (this *DataKeeper) Exit() {
	close(this.close)
}

func NewDataKeeper(ds DataSource, opts ...Option) *DataKeeper {
	dr := new(DataKeeper)
	dr.close = make(chan bool, 0)
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

func (this *DataKeeper) reload() {
	list := this.ds.Load()
	pool := new(sync.Map)
	for k, v := range list {
		pool.Store(k, v)
	}
	this.pool = pool
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

func (this *DataKeeper) tf() {
	timer := time.NewTimer(this.opt.timer)
	for {
		select {
		case <-timer.C:
			this.reload()
			timer.Reset(this.opt.timer)
		case <-this.close:
			break
		}
	}
}

func (this *DataKeeper) of() {
	var timer *time.Timer
	for {
		dur := getNextTime(this.opt.hour, this.opt.min, this.opt.sec)
		timer = time.NewTimer(dur)
		select {
		case <-timer.C:
			this.reload()
		case <-this.close:
			break
		}
	}
}

func (this *DataKeeper) uf() {
	timer := time.NewTimer(this.opt.uptimer)
	for {
		select {
		case <-timer.C:
			this.update.Range(func(key, value interface{}) bool {
				this.ds.Update(key, value)
				this.update.Delete(key)
				return true
			})
			timer.Reset(this.opt.uptimer)
		case <-this.close:
			break
		}
	}
}
