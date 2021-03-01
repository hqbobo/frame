package errors

// Errordef 错误定义
type Errordef struct {
	Code int32
	Msg  string
}

//SvcErrorDef 错误定义
type SvcErrorDef []Errordef

func (edef SvcErrorDef) getmsg(code int32) string {
	for _, v := range edef {
		if v.Code == code {
			return v.Msg
		}
	}
	return "未知错误"
}

//错误定义中文描述
var defs SvcErrorDef

//SetDefs 错误定义中文描述
func SetDefs(def SvcErrorDef) { defs = def }
