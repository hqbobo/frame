package handle

import (
	"encoding/xml"
	wxlog "github.com/hqbobo/frame/common/log"
	wxmodel "github.com/hqbobo/frame/common/weixin/model"
)

//处理接收普通消息
type TextHandle struct {
}

func (this *TextHandle) GetKey() string {
	return "text"
}

func (this *TextHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinTextMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnText(appid, msg)
	}

	return nil
}

//处理接收图片消息
type ImageHandle struct {
}

func (this *ImageHandle) GetKey() string {
	return "image"
}

func (this *ImageHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinImageMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnImage(appid, msg)
	}

	return nil
}

//处理接收语音消息
type VoiceHandle struct {
}

func (this *VoiceHandle) GetKey() string {
	return "voice"
}

func (this *VoiceHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinVoiceMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnVoice(appid, msg)
	}

	return nil
}

//处理接收视频消息
type VideoHandle struct {
}

func (this *VideoHandle) GetKey() string {
	return "video"
}

func (this *VideoHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinVideoMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnVideo(appid, msg)
	}

	return nil
}

//处理接小收视频消息
type ShortVideoHandle struct {
}

func (this *ShortVideoHandle) GetKey() string {
	return "shortvideo"
}

func (this *ShortVideoHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinVideoMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnVideo(appid, msg)
	}

	return nil
}

//处理链接消息
type LinkHandle struct {
}

func (this *LinkHandle) GetKey() string {
	return "link"
}

func (this *LinkHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinLinkMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnLink(appid, msg)
	}

	return nil
}

//处理位置消息
type LocationHandle struct {
}

func (this *LocationHandle) GetKey() string {
	return "location"
}

func (this *LocationHandle) Handle(appid string, data []byte) interface{} {
	msg := new(wxmodel.WexinLocationMessage)

	//结构转换
	e := xml.Unmarshal(data, msg)
	if e != nil {
		wxlog.Warn(e)
		return nil
	}

	if weixinMsgHandleImpl != nil {
		return weixinMsgHandleImpl.OnLocation(appid, msg)
	}

	return nil
}
