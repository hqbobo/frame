package gongzhonghao

import (
	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/weixin/gongzhonghao/handle"
	"github.com/hqbobo/frame/common/weixin/model"
	"github.com/hqbobo/frame/common/weixin/pay"
)

//WeiXinSessions 微信主结构
type WeiXinSessions struct {
	sess []*WeiXinSession
	cfgs []model.WeixinCfg
}

//NewSessions 新的微信批量主结构
func NewSessions(cfg []model.WeixinCfg) *WeiXinSessions {
	s := new(WeiXinSessions)
	s.cfgs = cfg
	cfgs = cfg
	for _, v := range cfg {
		s.sess = append(s.sess, NewSession(v))
	}
	return s
}

//Run 启动批量模式
func (wss *WeiXinSessions) Run(cb handle.WeixinMsgHandleInterface, eventcb handle.WeixinEventHandleInterface) {
	var e error
	for _, v := range wss.sess {
		log.Debug("WeiXinSessions ", v.cfg)
		if _, e = v.getToken(); e != nil {
			log.Debug(e)
		}
		v.GetMaterial()
		v.publicKey, e = pay.GetPublickKey(v.cfg.MchId, v.cfg.KeyFile, v.cfg.CertFile, v.cfg.WXRoot, v.cfg.APIKey)
		if e != nil {
			log.Warn(e)
		}

	}
	svrsInit(wss.cfgs, cb, eventcb)
}

//GetConfig 获取配置
func (wss *WeiXinSessions) GetConfig(appid string) (w model.WeixinCfg) {
	for _, v := range wss.cfgs {
		if v.Appid == appid {
			return v
		}
	}
	return model.WeixinCfg{}
}

//GetSess 查询指定微信sess
func (wss *WeiXinSessions) GetSess(appid string) (*WeiXinSession, bool) {
	for _, v := range wss.sess {
		if v.cfg.Appid == appid {
			return v, true
		}
	}
	return nil, false
}

//AnySess 微信sess
func (wss *WeiXinSessions) AnySess() (*WeiXinSession, bool) {
	for _, v := range wss.sess {
		return v, true
	}
	return nil, false
}
