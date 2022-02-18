package main

import "github.com/hqbobo/frame/common/log"

/*
*测试Log函数
 */
func main() {
	log.SetLevel(log.TraceLevel)
	log.Debugln("121111113")
	log.Infoln("1211113")
	log.Warnln("123123")
	log.Errorln("12123")
	log.Fatalln("12123213")
}
