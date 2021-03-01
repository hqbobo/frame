package main

import "github.com/hqbobo/frame/common/log"

/*
*测试Log函数
 */
func main() {
	log.SetLevel(log.TraceLevel)
	log.Debugln("123")
	log.Infoln("123")
	log.Warnln("123")
	log.Errorln("123")
	log.Fatalln("123")
}
