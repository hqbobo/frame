package handle

import (
	"encoding/xml"
	"github.com/hqbobo/frame/common/log"
	wxmodel "github.com/hqbobo/frame/common/weixin/model"
)

type WeixinMsgHandleInterface interface {
	OnText(appid string, msg *wxmodel.WexinTextMessage) interface{}
	OnImage(appid string, msg *wxmodel.WexinImageMessage) interface{}
	OnVoice(appid string, msg *wxmodel.WexinVoiceMessage) interface{}
	OnVideo(appid string, msg *wxmodel.WexinVideoMessage) interface{}
	OnLink(appid string, msg *wxmodel.WexinLinkMessage) interface{}
	OnLocation(appid string, msg *wxmodel.WexinLocationMessage) interface{}
}

var weixinMsgHandleImpl WeixinMsgHandleInterface
var weixinEventHanldeImpl WeixinEventHandleInterface

type WeixinEventHandleInterface interface {
	OnSubscribe(appid string, event *wxmodel.WeixinSubscribeEvent) interface{}
	OnScanner(appid string, event *wxmodel.WeixinScannerEvent) interface{}
	OnLocationEvent(appid string, event *wxmodel.WeixinLocationEvent) interface{}
	OnMenu(appid string, event *wxmodel.WeixinMenuEvent) interface{}
}

type handle interface {
	GetKey() string
	Handle(appid string, data []byte) interface{}
}

var handles map[string]handle
var eventhandles map[string]handle

//处理主入口
func Entrance(appid string, head wxmodel.WeixinBase, data []byte) interface{} {
	if head.MsgType.Value == "event" {
		var msg wxmodel.WeixinEventBase
		if e := xml.Unmarshal(data, &msg); e != nil {
			log.Warn(e)
		}
		if h, ok := eventhandles[msg.Event.Value]; ok {
			return h.Handle(appid, data)
		} else {
			log.Warn(msg.Event.Value, " event not found  ")
		}
	} else {
		if h, ok := handles[head.MsgType.Value]; ok {
			return h.Handle(appid, data)
		} else {
			log.Warn(head.MsgType.Value, " not found  ")
		}
	}
	return nil
}

func register(h handle) {
	handles[h.GetKey()] = h
}
func registerevent(h handle) {
	eventhandles[h.GetKey()] = h
}
func Init(cb WeixinMsgHandleInterface, event WeixinEventHandleInterface) {
	weixinMsgHandleImpl = cb
	weixinEventHanldeImpl = event
	handles = make(map[string]handle, 0)
	eventhandles = make(map[string]handle, 0)

	register(new(TextHandle))
	register(new(ImageHandle))
	register(new(VoiceHandle))
	register(new(LocationHandle))
	register(new(VideoHandle))
	register(new(ShortVideoHandle))
	register(new(LinkHandle))
	registerevent(new(SubscribeEventHandle))
	registerevent(new(ScannerEventHandle))
	registerevent(new(LocationEventHandle))
	registerevent(new(MenuEventHandle))
}
