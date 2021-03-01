package log

import (
	"testing"
)

/*
*测试Log函数
 */
func Test_Log(t *testing.T) {
	SetLevel(TraceLevel)
	Debugln("123")
	Infoln("123")
	Warnln("123")
	Errorln("123")
	Fatalln("123")
}
