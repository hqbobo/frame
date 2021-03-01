package recover_handle

import (
	"fmt"
	"github.com/hqbobo/frame/common/log"
	"runtime/debug"
)

//Recover 回收
func RecoverHandle(v ...interface{}) {
	if err := recover(); err != nil {
		if len(v) > 0 {
			s := fmt.Sprintln(v...)
			fmt.Println("s : ", s, " err : ", err)
			log.Errorf("recover err %s : ", err)
			log.Errorf("%s\n%s\n", s, string(debug.Stack()))
		} else {
			fmt.Println("v : ", v)
			log.Errorf("%s\n", string(debug.Stack()))
		}
	}
}
