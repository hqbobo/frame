package handle

import (
	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/weixin/model"
)

type WeixinReplyMsgHandleInterfce interface {
	OnTextMsg(msg *model.ReplyTextMessage) interface{}
	OnImageMsg(msg *model.ReplyImageMessage) interface{}
	OnVoiceMsg(msg *model.ReplyVoiceMessage) interface{}
	OnVideoMsg(msg *model.ReplyVideoMessage) interface{}
	OnMusicMsg(msg *model.ReplyMusicMessage) interface{}
	OnImageTextMsg(msg *model.ReplyImageTextMessage) interface{}
}

var WeixinReplyHandleImpl WeixinReplyMsgHandleInterfce

type replyHandle interface {
	GetKey() string
	ReplyHandle(data []byte) interface{}
}

var replyHandles map[string]replyHandle

func ReplyEntrace(head model.ReplyMessageBase, data []byte) interface{} {
	if h, ok := replyHandles[head.MsgType.Value]; ok {
		return h.ReplyHandle(data)
	} else {
		log.Warn(head.MsgType.Value, "未知类型")
	}
	return nil

}

func replyRegister(h replyHandle) {
	replyHandles[h.GetKey()] = h
}

func ReplyInit(rh WeixinReplyMsgHandleInterfce) {
	WeixinReplyHandleImpl = rh
	replyHandles = make(map[string]replyHandle, 0)
	replyRegister(new(ReplyTextHandle))
	replyRegister(new(ReplyImageHandle))
	replyRegister(new(ReplyVoicehandle))
	replyRegister(new(ReplyImageTextHandle))
	replyRegister(new(ReplyMusicHandle))
	replyRegister(new(ReplyVideoHandle))

}
