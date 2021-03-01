package common

import (
	"errors"

	"xorm.io/xorm"
	//"github.com/hqbobo/frame/common/log"
)

var (
	e *MysqlEngine
)

var ErrNotFound = errors.New("没找到")
var ErrUpdateFail = errors.New("没找到更新对象")
var ErrDeleteFail = errors.New("没找到删除对象")

type MysqlEngine struct {
	g *xorm.EngineGroup
}

func InitNewMysqlEngine(group *xorm.EngineGroup) {
	e = new(MysqlEngine)
	e.g = group
	return
}

type BaseRepositoryImpl struct {
}

//获取读写实例
func (this *BaseRepositoryImpl) GetEngine() *xorm.Engine {
	return e.g.Master()
}

//获取只读实例
func (this *BaseRepositoryImpl) GetReadEngine() *xorm.Engine {
	return e.g.Slave()
}
