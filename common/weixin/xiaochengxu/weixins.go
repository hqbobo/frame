package xiaochengxu

import (
	"github.com/hqbobo/frame/common/weixin/model"
)

//WeiXinMiniSession 微信主结构
type WeiXinMiniSessions struct {
	sess []*WeiXinMiniSession
	cfgs []model.WeixinCfg
}

func NewSessions(cfg []model.WeixinCfg) *WeiXinMiniSessions {
	s := new(WeiXinMiniSessions)
	s.cfgs = cfg
	for _, v := range cfg {
		sess := NewSession(v)
		s.sess = append(s.sess, sess)
		sess.Run()
	}
	return s
}

func (wss *WeiXinMiniSessions) GetConfig(appid string) (w model.WeixinCfg) {
	for _, v := range wss.cfgs {
		if v.Appid == appid {
			return v
		}
	}
	return model.WeixinCfg{}
}

func (wss *WeiXinMiniSessions) GetSess(appid string) (*WeiXinMiniSession, bool) {
	for _, v := range wss.sess {
		if v.cfg.Appid == appid {
			return v, true
		}
	}
	return nil, false
}
