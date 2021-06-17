package conf

import (
	"encoding/json"

	"github.com/hqbobo/frame/common/log"

	//"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"runtime"
	"testing"
)

//
type registry struct {
	Type  string
	Addrs []string
}

type broker struct {
	Type  string
	Addrs []string
}

type dlock struct {
	Type string
	Addr string
}

type mysqlConn struct {
	User     string `json:"user"`
	Password string `json:"pwd"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       string `json:"db"`
	Size     int    `json:"size"`
	SqlLog   bool   `json:"sqlLog"`
	Sync     bool   `json:"sync"`
	Cache    bool   `json:"cache"`
}

type mysqlConfig struct {
	Read  []mysqlConn
	Write []mysqlConn
}

type dbConfig struct {
	Types       string
	MysqlConfig mysqlConfig
}

type cache struct {
	Type  string
	Pass  string
	Addrs []string
}

type service struct {
	Ttl      int
	Interval int
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	Host      string
	Registry  registry
	Broker    broker
	Cache     cache
	Tracelink string
	Dlock     dlock
	Service   service
	DbConfig  dbConfig
	Version   string
	Pprof     string
	Rabbitmq  string
}

var (
	config GlobalConfig
)

// CFG 获取配置
func CFG() *GlobalConfig {
	return &config
}

// ParseConfig 解析配置文件
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Warnln("use -c to specify configuration file")
	}

	configContent, err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Warnln("read config file: %s err %s", cfg, err)
	}

	err = json.Unmarshal([]byte(configContent), &config)
	if err != nil {
		log.Warnln("parse config file: %s err %s", cfg, err)
	}
	log.Debugln("read config file: %s successfully", cfg)
}

// ParseStr 解析配置字符串
func ParseStr(str string) {
	err := json.Unmarshal([]byte(str), &config)
	if err != nil {
		log.Warnln("parse config string: %s err %s", str, err)
	}
}

var _ = func() bool {
	testing.Init()
	return true
}()

type runMode string

func (rm *runMode) IsTest() bool {
	if *rm == "test" {
		return true
	}
	return false
}

func (rm *runMode) IsLocal() bool {
	if *rm == "local" {
		return true
	}
	return false
}
func (rm *runMode) Set(mode string) { *rm = runMode(mode) }
func (rm *runMode) Get() string     { return string(*rm) }

// Mode 运行模式
var Mode runMode

//配置文件初始化
func init() {
	var filepath string
	line := "\n"
	Mode.Set(os.Getenv("RUNMODE"))
	if Mode.Get() == "develop" {
		Mode.Set("local")
	}
	//log.SetLogLevel(log.LevelDebug)
	//log.SetLogModel(log.ModelDev)
	if s := os.Getenv("config"); s == "" {
		envStr := Mode.Get() + ".json"
		fmt.Println("envStr >>> ", envStr)
		_, filename, _, _ := runtime.Caller(0)
		filepath = path.Join(path.Dir(filename), envStr)
		//cfg := flag.String("c", "", "配置文件")
		////配置文件
		//help := flag.Bool("h", false, "help")
		//flag.Parse()
		//if *help {
		//	flag.Usage()
		//	os.Exit(0)
		//}
		//if cfg != nil && len(*cfg) > 0 {
		//	filepath = *cfg
		//	ParseConfig(filepath)
		//} else {
		ParseConfig(filepath)
		//}
	} else {
		line = line + addline("配置来源", "环境变量")
		//fmt.Println("配置信息：：：", s)
		ParseStr(s)
	}
	line = line + addline("运行模式", Mode.Get())
	line = line + addline("配置文件", filepath)
	line = line + addline("当前版本", config.Version)
	line = line + addline("注册中心", config.Registry.Type)
	line = line + addlineArray("注册地址", config.Registry.Addrs)
	line = line + addline("消息队列", config.Broker.Type)
	line = line + addlineArray("队列地址", config.Broker.Addrs)
	line = line + addline("缓存数据", config.Cache.Type)
	line = line + addlineArray("缓存地址", config.Cache.Addrs)
	line = line + addline("链路追踪", config.Tracelink)
	line = line + addline("持久数据", config.DbConfig.Types)
	line = line + addline("分布式锁", config.Dlock.Type)
	line = line + addline("锁地址", config.Dlock.Addr)
	line = line + addline("pprof", config.Pprof)
	log.Infof(line)
	go func() {
		if config.Pprof == "" {
			config.Pprof = ":9876"
		}
		fmt.Println("pprof start at :", config.Pprof)
		fmt.Println(http.ListenAndServe(config.Pprof, nil))
	}()
}

const (
	fmtstr = "%30s:%-50s\n"
)

func addlineArray(types string, value []string) string {
	var s string
	if len(value) > 0 {
		for k, v := range value {
			if k == 0 {
				s = fmt.Sprintf(fmtstr, types, v)
			} else {
				s = s + fmt.Sprintf(fmtstr, "", value)
			}
		}
		return s
	}
	return ""
}

func addline(types, value string) string {
	if len(value) > 0 {
		return fmt.Sprintf(fmtstr, types, value)
	}
	return ""
}
