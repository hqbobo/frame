package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/hqbobo/frame/common/conf"
	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/db/common"

	//初始化mysql
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/names"
)

//MysqlConn 数据库
type MysqlConn struct {
	User     string `json:"user"`
	Password string `json:"pwd"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       string `json:"db"`
	Size     int    `json:"size"`
	SQLLog   bool   `json:"sqlLog"`
	Sync     bool   `json:"sync"`
}

//MysqlConfig MYSQL数据库配置
type MysqlConfig struct {
	Read  []MysqlConn
	Write []MysqlConn
}

//DbConfig 数据库配置
type DbConfig struct {
	Types       string
	MysqlConfig MysqlConfig
}

func mysqlInit(config *conf.GlobalConfig) error {
	w := make([]*xorm.Engine, 0)
	r := make([]*xorm.Engine, 0)
	//分页缓存
	cacher := caches.NewLRUCacher(newCacheStore(config.Cache.Addrs, config.Cache.Pass), 5000)
	for _, v := range config.DbConfig.MysqlConfig.Write {
		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=%s", v.User, v.Password,
			v.Host, v.Port, v.DB, "Local")
		e, err := xorm.NewEngine("mysql", dbURL)
		if err != nil {
			log.Warnln(err)
		}
		log.Debugln(v)

		e.ShowSQL(v.SqlLog)
		err = e.Ping()

		loger := new(Logger)
		e.SetLogger(loger)
		if err != nil {
			log.Warnln(err)
		}
		log.Debugln("mysql 实例 ", dbURL, " 连接成功")
		e.SetMapper(names.SameMapper{})
		e.SetMaxOpenConns(v.Size)
		e.ShowSQL(v.SqlLog)
		if v.Sync == true {
			log.Debugln("mysql write sync")
		}
		if v.Cache || os.Getenv("DBWCACHE") != "" {
			log.Info("mysql 实例 启动缓存")
			e.SetDefaultCacher(cacher)
		}
		w = append(w, e)
	}
	for _, v := range config.DbConfig.MysqlConfig.Read {
		dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=%s", v.User, v.Password,
			v.Host, v.Port, v.DB, "Local")
		e, err := xorm.NewEngine("mysql", dbURL)
		if err != nil {
			log.Warnln(err)
		}
		err = e.Ping()
		log.Debugln(v)
		loger := new(Logger)
		e.SetLogger(loger)
		if err != nil {
			log.Warnln(err)
		}
		log.Debugln("mysql 只读实例 ", dbURL, " 连接成功")
		log.Debug(e)
		e.ShowSQL(v.SqlLog)
		//设置字段映射规则
		e.SetMapper(names.SameMapper{})
		e.SetMaxOpenConns(v.Size)
		if v.Cache || os.Getenv("DBRCACHE") != "" {
			log.Info("mysql 只读实例 启动缓存")
			e.SetDefaultCacher(cacher)
		}
		r = append(r, e)
	}
	group, err := xorm.NewEngineGroup(w[0], r)
	if err != nil {
		log.Warn(err)
	}
	common.InitNewMysqlEngine(group)
	return nil
}

//InitDal 初始化数据库
func InitDal(config *conf.GlobalConfig) error {
	if config == nil {
		return nil
	}
	switch config.DbConfig.Types {
	case "mysql":
		return mysqlInit(config)
	default:
		return errors.New("未知数据库类型")
	}
}
