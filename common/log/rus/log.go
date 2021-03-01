package rus

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	sysTime "time"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
	elogrus "gopkg.in/sohlich/elogrus.v2"
)

//日志等级
const (
	LEVEL_PANIC = iota
	LEVEL_FATA
	LEVEL_ERROR
	LEVEL_WARNING
	LEVEL_INFO
	LEVEL_DEBUG
	LEVEL_TRACE
)

const (
	C_Kafka_Log_Topic      = "all-server-log-test2"
	C_Kafka_Log_User_Topic = "all-user-log-test2"
)

//日志模式
const (
	MODEL_PRO = iota
	MODEL_INFO
	MODEL_DEV
	MODEL_TRACE
)

//调用log的服务器名字
var serverName string
var g_ip net.IP
var mylog = logrus.New()
var log = logrus.New()
var address string
var is_open = false

func init() {
	ip, err := externalIP()
	if err == nil {
		g_ip = ip
	} else {
		g_ip = nil
	}
	log.SetLevel(logrus.TraceLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	mylog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	mylog.SetOutput(os.Stdout)
	log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
}

func SetLogName(name string) {
	serverName = name
}

func OpenSendLog(name string, open bool, elkAddress string) {
	serverName = name
	is_open = open
	address = elkAddress
	if is_open {
		ConfigESLogger(address, g_ip.String(), serverName)
	}
}

var callerLevel = 2

func SetCallerLevel(level int) {
	callerLevel = level
}

func callerPrettyfier() (path string) {
	fname := ""
	pc, path, line, ok := runtime.Caller(callerLevel) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)
	return fmt.Sprintf("%s() %s:%d ", funcName, path, line)
}

func ConfigESLogger(esUrl string, esHOst string, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esUrl))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	esHook, err := elogrus.NewAsyncElasticHook(client, esHOst, logrus.WarnLevel, index)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.AddHook(esHook)

	esUserHook, err := elogrus.NewAsyncElasticHook(client, esHOst, logrus.WarnLevel, "用户操作记录测试")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	mylog.AddHook(esUserHook)
	return nil
}

//设置日志等级
func SetLogLevel(lev logrus.Level) {
	log.SetLevel(lev)
}

/*
*设置日志模式
*参数说明:
*@param:mod		模式 MODEL_PRO:只向日志服务器发送日志  MODEL_INFO:向日志服务器发送日志切输出到控制台 MODEL_DEV:只输出到控制台
 */
func SetLogModel(mod int) error {
	if mod <= MODEL_TRACE {
		mylog.SetOutput(os.Stdout)
		mylog.SetLevel(logrus.TraceLevel)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
		log.SetLevel(logrus.TraceLevel)
	}
	if mod <= MODEL_DEV {
		mylog.SetOutput(os.Stdout)
		mylog.SetLevel(logrus.DebugLevel)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
		log.SetLevel(logrus.DebugLevel)
	}
	if mod <= MODEL_INFO {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		writer := bufio.NewWriter(src)
		mylog.SetOutput(writer)
		mylog.SetLevel(logrus.InfoLevel)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
		log.SetLevel(logrus.InfoLevel)
	}
	if mod <= MODEL_PRO {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		writer := bufio.NewWriter(src)
		log.SetOutput(writer)
		log.SetLevel(logrus.WarnLevel)
		mylog.SetOutput(writer)
		mylog.SetLevel(logrus.WarnLevel)
	}
	return nil
}

//安全执行监听函数
func Listen(f interface{}, callback func(interface{}), param string) {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)
	timer := sysTime.NewTicker(sysTime.Millisecond * 500)
	success := make(chan bool)
	start := sysTime.Now()
	count := 0
	go func() {
		callback(f)
		close(success)
	}()
	for {
		select {
		case <-success:
			end := sysTime.Since(start).Nanoseconds() / 1000000.00
			timer.Stop()
			if end >= 500 && end < 1000 {
				log.Info(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", path, funcName, line, end, param))
			}
			if end >= 1000 && end < 2000 {
				log.Warn(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", path, funcName, line, end, param))
			}
			if end >= 2000 {
				log.Error(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", path, funcName, line, end, param))
			}
			return
		case <-timer.C:
			count++
			end := sysTime.Since(start).Nanoseconds() / 1000000.00
			if count >= 10 {
				log.Error(fmt.Sprintf("执行严重超时%d次提醒 %s %s %d (%dms) %s", count, path, funcName, line, end, param))
			} else {
				log.Info(fmt.Sprintf("执行严重超时%d次提醒 %s %s %d (%dms) %s", count, path, funcName, line, end, param))
			}
		}
	}
}

//计算函数所用时间
func TraceParam(param string) func() {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			Debug("执行严重超时提醒 %s %s %d (%dms) %s", path, funcName, line, end, param)
		}
		if end >= 1000 && end < 2000 {
			Warn("执行严重超时提醒 %s %s %d (%dms) %s", path, funcName, line, end, param)
		}
		if end >= 2000 {
			Error("执行严重超时提醒 %s %s %d (%dms) %s", path, funcName, line, end, param)
		}
	}
}

//计算函数所用时间
func RunTime() func() {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			Debug("执行严重超时提醒 %s %s %d (%dms)", path, funcName, line, end)
		}
		if end >= 1000 && end < 2000 {
			Warn("执行严重超时提醒 %s %s %d (%dms)", path, funcName, line, end)
		}
		if end >= 2000 {
			Error("执行严重超时提醒 %s %s %d (%dms)", path, funcName, line, end)
		}
	}
}

func GetFunctionName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

func UserInfoLog(userId int64, format string, args ...interface{}) {
	mylog.WithFields(logrus.Fields{
		"UserId": userId,
	}).Warn(callerPrettyfier() + fmt.Sprintf(format, args...))
}

//危险的
func Fatal(format string, args ...interface{}) {
	log.Fatal(callerPrettyfier() + fmt.Sprintf(format, args...))
}

//错误
func Error(format string, args ...interface{}) {
	log.Error(callerPrettyfier() + fmt.Sprintf(format, args...))
}

//警告
func Warn(format string, args ...interface{}) {
	log.Warn(callerPrettyfier() + fmt.Sprintf(format, args...))
}

//提示
func Info(format string, args ...interface{}) {
	log.Info(callerPrettyfier() + fmt.Sprintf(format, args...))
}

//调试
func Debug(format string, args ...interface{}) {
	log.Debug(callerPrettyfier() + fmt.Sprintf(format, args...))
}

func Trace(format string, args ...interface{}) {
	log.Trace(callerPrettyfier() + fmt.Sprintf(format, args...))
}

func Fatalf(args ...interface{}){
	log.Fatal(callerPrettyfier() + fmt.Sprint( args...))
}
func Errorf(args ...interface{}){
	log.Error(callerPrettyfier() + fmt.Sprint( args...))
}

func Warnf(args ...interface{}){
	log.Warn(callerPrettyfier() + fmt.Sprint( args...))
}

func Infof(args ...interface{}){
	log.Info(callerPrettyfier() + fmt.Sprint( args...))
}

func Debugf(args ...interface{}){
	log.Debug(callerPrettyfier() + fmt.Sprint( args...))
}

func Tracef(args ...interface{}){
	log.Trace(callerPrettyfier() + fmt.Sprint( args...))
}


func lastFname(fname string) string {
	flen := len(fname)
	n := strings.LastIndex(fname, ".")
	if n+1 < flen {
		return fname[n+1:]
	}
	return fname
}

func getFilePath(path string) string {
	s := strings.Split(path, "src")
	return s[1]
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
