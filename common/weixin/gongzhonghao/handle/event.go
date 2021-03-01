package handle

import (
	"encoding/xml"
	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/weixin/model"
)

//处理关注和取消关注事件
type SubscribeEventHandle struct {
}

func (this *SubscribeEventHandle) GetKey() string {
	return "subscribe"
}

func (this *SubscribeEventHandle) Handle(appid string, data []byte) interface{} {
	event := new(model.WeixinSubscribeEvent)

	e := xml.Unmarshal(data, event)

	if e != nil {
		log.Warn(e)
		return nil
	}

	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnSubscribe(appid, event)
	}

	return nil
}

//处理扫描二维码事件
type ScannerEventHandle struct {
}

func (this *ScannerEventHandle) GetKey() string {
	return "Scanner"
}
func (this *ScannerEventHandle) Handle(appid string, data []byte) interface{} {
	event := new(model.WeixinScannerEvent)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warn(e)
		return nil
	}
	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnScanner(appid, event)
	}
	return nil
}

//处理上报地理位置事件
type LocationEventHandle struct {
}

func (this *LocationEventHandle) GetKey() string {
	return "location"
}
func (this *LocationEventHandle) Handle(appid string, data []byte) interface{} {
	event := new(model.WeixinLocationEvent)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warn(e)
		return nil
	}
	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnLocationEvent(appid, event)
	}
	return nil
}

//处理自定义菜单事件
type MenuEventHandle struct {
}

func (this *MenuEventHandle) GetKey() string {
	return "CLICK"
}
func (this *MenuEventHandle) Handle(appid string, data []byte) interface{} {
	event := new(model.WeixinMenuEvent)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warn(e)
		return nil
	}
	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnMenu(appid, event)
	}
	return nil
}

//推送完成
type TEMPLATESENDJOBFINISHHandle struct {
}

func (this *TEMPLATESENDJOBFINISHHandle) GetKey() string {
	return "TEMPLATESENDJOBFINISH"
}
func (this *TEMPLATESENDJOBFINISHHandle) Handle(appid string, data []byte) interface{} {
	event := new(model.WeixinEventBase)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warn(e)
		return nil
	}
	log.Trace("do nonthing")
	return nil
}
