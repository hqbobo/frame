package datakeeper

import (
	"testing"
	"time"
)

var testdata = []string{"test1", "test2"}
var testtime map[int]int64

func init() {
	testtime = make(map[int]int64, 0)
	testtime[time.Now().Day()] = time.Now().Unix()

}

type data struct {
}

func (rds *data) Load() map[interface{}]interface{} {
	out := make(map[interface{}]interface{}, 0)
	return out
}
func (rds *data) Add(key, val interface{}) {

}

func (rds *data) Update(key, val interface{}) {

}

func (rds *data) Delete(key interface{}) {

}

// TestFunc 基础功能测试
func TestFunc(t *testing.T) {
	dk := NewDataKeeper(&data{})
	now := time.Now().Unix()

	dk.Store(time.Now().Day(), now)
	val, ok := dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	time.Sleep(time.Second)
	dk.Store(time.Now().Day(), time.Now().Unix())
	val, ok = dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	if val == now {
		t.Error("更新失败")
		return
	}
	t.Log(now, "--", val)
	dk.Delete(time.Now().Day())
	val, ok = dk.Load(time.Now().Day())
	if ok {
		t.Error("删除失败找到数据 ", val)
		return
	}
}

type reloadData struct {
}

func (rds *reloadData) Add(key, val interface{}) {

}

func (rds *reloadData) Load() map[interface{}]interface{} {
	out := make(map[interface{}]interface{}, 0)
	out[time.Now().Day()] = time.Now().Unix()
	return out
}

func (rds *reloadData) Update(key, val interface{}) {

}

func (rds *reloadData) Delete(key interface{}) {

}

// TestTimer 测试定时触发load
func TestTimer(t *testing.T) {
	dk := NewDataKeeper(&reloadData{}, WithTriggerTimer(time.Second))
	val, ok := dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	time.Sleep(time.Second * 2)

	val1, ok := dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	if val == val1 {
		t.Error("定时更新失败")
		return
	}
	t.Log("val:", val, "--- val1:", val1)
}

// TestDaily 测试每天定时触发load
func TestDaily(t *testing.T) {
	dk := NewDataKeeper(&reloadData{}, WithTriggerDaily(time.Now().Hour(),
		time.Now().Minute(), time.Now().Second()+3))

	val, ok := dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	time.Sleep(time.Second * 4)
	val1, ok := dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	if val == val1 {
		t.Error("定时更新失败")
		return
	}
	t.Log("val:", val, "--- val1:", val1)
}

type updateData struct {
}

func (rds *updateData) Load() map[interface{}]interface{} {
	out := make(map[interface{}]interface{}, 0)
	for k, v := range testtime {
		out[k] = v
	}
	return out
}

func (rds *updateData) Add(key, val interface{}) {

}

func (rds *updateData) Update(key, val interface{}) {
	testtime[key.(int)] = val.(int64)
}

func (rds *updateData) Delete(key interface{}) {

}

// TestTimer 测试定时触发load
func TestUpdate(t *testing.T) {
	dk := NewDataKeeper(&updateData{}, WithUpdateTimer(time.Second*2))
	val, ok := dk.Load(time.Now().Day())
	if !ok {
		t.Error("没有找到数据")
		return
	}
	initdata := val
	t.Log("原始数据:", val)
	time.Sleep(time.Second)
	dk.Store(time.Now().Day(), time.Now().Unix())
	val, ok = testtime[time.Now().Day()]
	if !ok {
		t.Error("没有找到数据")
		return
	}
	t.Log("Store后数据:", val)
	if initdata == val {
		t.Error("原始数据不应更新还未到时间")
		return
	}
	time.Sleep(time.Second * 2)
	val, ok = testtime[time.Now().Day()]
	if !ok {
		t.Error("没有找到数据")
		return
	}
	if initdata != val {
		t.Error("原始数据更新失败")
		return
	}
	t.Log("Store后等待更新时间后数据:", val)
}

func BenchmarkLoad(b *testing.B) {
	dk := NewDataKeeper(&updateData{})
	day := time.Now().Day()
	for i := 0; i < b.N; i++ {
		_, ok := dk.Load(day)
		if !ok {
			b.Error("没有找到数据")
			return
		}
	}

}

func BenchmarkStore(b *testing.B) {
	dk := NewDataKeeper(&updateData{})
	day := time.Now().Day()
	data := time.Now().Unix()
	for i := 0; i < b.N; i++ {
		dk.Store(day, data)
	}
}
