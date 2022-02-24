package gongzhonghao

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/utils"
	"github.com/hqbobo/frame/common/weixin/gongzhonghao/handle"
	wxmodel "github.com/hqbobo/frame/common/weixin/model"
	"github.com/hqbobo/frame/common/weixin/pay"
)

//WeixinCacheInterface 微信緩存接口
type WeixinCacheInterface interface {
	Get(key string, data interface{}) bool
	Set(key string, data interface{}) bool
	Delete(key string) bool
}

type weixntoken struct {
	Token   string
	TimeOut int64
}

type weixinticket struct {
	Ticket        string
	TicketTimeOut int64
}

//WeiXinSession 微信主结构
type WeiXinSession struct {
	weixntoken
	weixinticket
	cfg       wxmodel.WeixinCfg
	cache     WeixinCacheInterface
	publicKey string
}

//_HttpPost 封装
func (ws *WeiXinSession) _HttpPost(url string, data []byte) (body []byte) {
	var rsp wxmodel.ErrRsp
	var e error
	r := utils.HttpPost(url, data)
	e = json.Unmarshal(r, &rsp)
	if e != nil {
		log.Warn(e)
		return r
	}
	//过期了
	if rsp.Errcode == wxmodel.ErrornvalidCredential {
		ws.TokenReset()
	}
	return r
}

//_HttpsGet 封装
func (ws *WeiXinSession) _HttpsGet(url string) ([]byte, error) {
	var rsp wxmodel.ErrRsp
	var e error
	r, err := utils.HttpsGet(url)
	e = json.Unmarshal(r, &rsp)
	if e != nil {
		log.Warn(e)
		return r, e
	}
	//过期了
	if rsp.Errcode == wxmodel.ErrornvalidCredential {
		ws.TokenReset()
	}
	return r, err
}

//_HttpGet 封装
func (ws *WeiXinSession) _HttpGet(url string) []byte {
	var rsp wxmodel.ErrRsp
	var e error
	r := utils.HttpGet(url)
	e = json.Unmarshal(r, &rsp)
	if e != nil {
		log.Warn(e)
		return r
	}
	//过期了
	if rsp.Errcode == wxmodel.ErrornvalidCredential {
		ws.TokenReset()
	}
	return r
}

//CFG 配置
func (ws *WeiXinSession) CFG() wxmodel.WeixinCfg { return ws.cfg }

//在线获取access_token
func (ws *WeiXinSession) reloadToken() error {
	//优先缓存获取token
	var token weixntoken
	if ws.cache != nil {
		if ws.cache.Get(CacheTokenName+ws.cfg.Appid, &token) {
			//更新本地内存
			log.Debug("["+ws.cfg.Appid+"]读取缓存token到本地:", token.Token, " 到期时间:", time.Now().Add(time.Second*time.Duration(token.TimeOut-time.Now().Unix())).String())
			ws.weixntoken = token
		}
	}
	log.Debug("time:", ws.TimeOut, "now:", time.Now().Unix(), ws.cfg.Appid, ws.cfg.Appsecret)
	//还是超时或者没有
	if ws.Token == "" || ws.TimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp wxmodel.TokenRsp
		r, e := ws._HttpsGet(fmt.Sprintf(TokenURL, ws.cfg.Appid, ws.cfg.Appsecret))
		if e != nil {
			log.Warn(e)
			return e
		}
		e = json.Unmarshal(r, &rsp)
		if e != nil {
			log.Warn(e)
			return e
		}
		log.Debug("在线获取新的token :", rsp.Access_token)
		ws.Token = rsp.Access_token
		ws.TimeOut = int64(rsp.Expires_in-10) + time.Now().Unix()
		if rsp.Errcode != 0 {
			return errors.New(rsp.ErrRsp.Errmsg)
		}
		//保存到缓存
		if ws.cache != nil {
			ws.cache.Set(CacheTokenName+ws.cfg.Appid, &ws.weixntoken)
		}
	}

	return nil
}

//GetAccessToken 获取token
func (ws *WeiXinSession) GetAccessToken() (string, error) {
	return ws.getToken()
}

//本地获取access_token
func (ws *WeiXinSession) getToken() (string, error) {
	//检查token是否已经获取和超时
	if ws.Token == "" || ws.TimeOut <= time.Now().Unix() {
		if e := ws.reloadToken(); e != nil {
			return "", e
		}

	}
	return ws.Token, nil
}

//Ticket
func (ws *WeiXinSession) reloadTicket() error {
	ws.getToken()

	//优先缓存获取Ticket
	var ticket weixinticket
	if ws.cache != nil {
		if ws.cache.Get(CacheTicketName+ws.cfg.Appid, &ticket) {
			//更新本地内存
			log.Debug("读取缓存ticket到本地:", ticket.Ticket, " 到期时间:", time.Now().Add(time.Second*time.Duration(ticket.TicketTimeOut-time.Now().Unix())).String())
			ws.weixinticket = ticket
		}
	}
	log.Debug("time:", ws.TicketTimeOut, "now:", time.Now().Unix())
	//还是超时或者没有
	if ws.Ticket == "" || ws.TicketTimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp wxmodel.TicketRsp
		r, e := ws._HttpsGet(fmt.Sprintf(GetTicketStr, ws.Token))
		log.Debug(fmt.Sprintf(GetTicketStr, ws.Token), ": ", string(r))
		if e != nil {
			log.Warn(e)
			return e
		}
		e = json.Unmarshal(r, &rsp)
		if e != nil {
			log.Warn(e)
			return e
		}
		log.Debug("在线获取新的ticket:", rsp.Ticket)
		ws.Ticket = rsp.Ticket
		ws.TicketTimeOut = int64(rsp.Expires_in) + time.Now().Unix()
		if rsp.Errcode != 0 {
			return errors.New(rsp.ErrRsp.Errmsg)
		}
		//保存到缓存
		if ws.cache != nil {
			ws.cache.Set(CacheTicketName+ws.cfg.Appid, &ws.weixinticket)
		}
	}

	return nil
}

//本地获取ticket
func (ws *WeiXinSession) getTicket() (string, error) {
	//检查ticket是否已经获取和超时
	if ws.Ticket == "" || ws.TicketTimeOut <= time.Now().Unix() {
		if e := ws.reloadTicket(); e != nil {
			return "", e
		}

	}
	return ws.Ticket, nil
}

//获取微信服务器IP地址
func (ws *WeiXinSession) getWeixinSvrIP() error {
	//更新token
	var rsp wxmodel.ServerListRsp
	r, e := ws._HttpsGet(fmt.Sprintf(GetServerURL, ws.Token))
	if e != nil {
		log.Warn(e)
		return e
	}
	e = json.Unmarshal(r, &rsp)
	if e != nil {
		log.Warn(e)
		return e
	}
	log.Debug(rsp)
	if rsp.Errcode != 0 {
		return errors.New(rsp.ErrRsp.Errmsg)
	}
	return nil
}

//SetCacheInterface 设置缓存接口
func (ws *WeiXinSession) SetCacheInterface(cache WeixinCacheInterface) {
	ws.cache = cache
}

//Run 运行
func (ws *WeiXinSession) Run(cb handle.WeixinMsgHandleInterface, eventcb handle.WeixinEventHandleInterface) {
	var e error
	if _, e = ws.getToken(); e != nil {
		log.Warn(e)
	}
	ws.publicKey, e = pay.GetPublickKey(ws.cfg.MchId, ws.cfg.KeyFile, ws.cfg.CertFile, ws.cfg.WXRoot, ws.cfg.APIKey)
	if e != nil {
		log.Info(e)
	}
	svrInit(ws, cb, eventcb)
}

//TokenReset token重置
func (ws *WeiXinSession) TokenReset() {
	ws.Token = ""
	ws.Ticket = ""
	ws.cache.Delete(CacheTicketName + ws.cfg.Appid)
	ws.cache.Delete(CacheTokenName + ws.cfg.Appid)
	ws.reloadToken()
	ws.reloadTicket()
}

//NewSession 创建一个新的微信实例
func NewSession(cfg wxmodel.WeixinCfg) *WeiXinSession {
	O := new(WeiXinSession)
	O.cfg = cfg
	return O
}

//GetUserToken 根据code获取用户token
func (ws *WeiXinSession) GetUserToken(code string) (*wxmodel.UserToken, error) {
	url := fmt.Sprintf(weixin_user_token, ws.cfg.Appid, ws.cfg.Appsecret, code)
	body := ws._HttpGet(url)
	log.Info(string(body))
	usr := new(wxmodel.UserToken)
	if err := json.Unmarshal(body, usr); err != nil {
		return nil, err
	}
	if usr.Errcode != 0 {
		return nil, errors.New(usr.Errmsg)
	}
	log.Debug("GetUserToken %+v", usr)
	return usr, nil
}

//GetUInfo 获取用户信息
func (ws *WeiXinSession) GetUInfo(openid, token string) (*wxmodel.UInfo, error) {
	if _, err := ws.getToken(); err != nil {
		log.Error(err)
		return nil, err
	}

	url := fmt.Sprintf(weixin_user_info, token, openid)
	body := ws._HttpGet(url)
	log.Info(string(body))
	usr := new(wxmodel.UInfo)
	if err := json.Unmarshal(body, usr); err != nil {
		return nil, err
	}
	log.Debug("UInfo %+v", usr)
	return usr, nil
}

var topColor = "#FF0000"

//SendTemplateMsg 发送模板消息
func (ws *WeiXinSession) SendTemplateMsg(openid, templateid string, data interface{}) error {
	if _, err := ws.getToken(); err != nil {
		log.Error(err)
		return err
	}
	wtc := wxmodel.WeixinTemp{}
	wtc.TemplateId = templateid
	wtc.ToUser = openid
	wtc.Topcolor = topColor
	wtc.Data = data
	buf, _ := json.Marshal(wtc)
	url := weixin_template_url + "?access_token=" + ws.Token

	log.Trace(string(buf))
	body := ws._HttpPost(url, []byte(buf))
	wxError := new(wxmodel.ErrRsp)
	if err := json.Unmarshal(body, wxError); err != nil {
		return err
	}

	log.Trace(wxError)
	if wxError.Errcode != 0 {
		body := ws._HttpPost(url, []byte(buf))
		if err := json.Unmarshal(body, wxError); err != nil {
			return err
		}
		if wxError.Errcode != 0 {
			fmt.Println(wxError)
			return errors.New(wxError.Errmsg)
		}
	}
	return nil
}

//SendTemplateMsg 发送模板消息
func (ws *WeiXinSession) SendTemplateMsgJumpUrl(openid, templateid string, data interface{}, jurl string) error {
	if _, err := ws.getToken(); err != nil {
		log.Error(err)
		return err
	}
	wtc := wxmodel.WeixinTemp{}
	wtc.TemplateId = templateid
	wtc.ToUser = openid
	wtc.Topcolor = topColor
	wtc.Data = data
	wtc.Url = jurl
	buf, _ := json.Marshal(wtc)
	url := weixin_template_url + "?access_token=" + ws.Token

	log.Trace(string(buf))
	body := ws._HttpPost(url, []byte(buf))
	wxError := new(wxmodel.ErrRsp)
	if err := json.Unmarshal(body, wxError); err != nil {
		return err
	}

	log.Trace(wxError)
	if wxError.Errcode != 0 {
		body := ws._HttpPost(url, []byte(buf))
		if err := json.Unmarshal(body, wxError); err != nil {
			return err
		}
		if wxError.Errcode != 0 {
			fmt.Println(wxError)
			return errors.New(wxError.Errmsg)
		}
	}
	return nil
}

//SendTemplateMsgJumpMini 发送模板消息
func (ws *WeiXinSession) SendTemplateMsgJumpMini(openid, templateid string, data interface{}, appid, jump string) error {
	if _, err := ws.getToken(); err != nil {
		log.Error(err)
		return err
	}
	wtc := wxmodel.WeixinTemp{}
	wtc.TemplateId = templateid
	wtc.ToUser = openid
	wtc.Topcolor = topColor
	wtc.Data = data
	wtc.Mini.Appid = appid
	wtc.Mini.Pagepath = jump
	buf, _ := json.Marshal(wtc)
	url := weixin_template_url + "?access_token=" + ws.Token

	log.Trace(string(buf))
	body := ws._HttpPost(url, []byte(buf))
	wxError := new(wxmodel.ErrRsp)
	if err := json.Unmarshal(body, wxError); err != nil {
		return err
	}

	log.Trace(wxError)
	if wxError.Errcode != 0 {
		body := ws._HttpPost(url, []byte(buf))
		if err := json.Unmarshal(body, wxError); err != nil {
			return err
		}
		if wxError.Errcode != 0 {
			fmt.Println(wxError)
			return errors.New(wxError.Errmsg)
		}
	}
	return nil
}

// SendKfText 发送客服消息
func (ws *WeiXinSession) SendKfText(openid, msg string) error {
	if _, err := ws.getToken(); err != nil {
		log.Error(err)
		return err
	}
	wtc := wxmodel.WeixinKfText{}
	wtc.ToUser = openid
	wtc.Msgtype = "text"
	wtc.Text.Content = msg
	buf, _ := json.Marshal(wtc)
	url := weixin_kf_url + "?access_token=" + ws.Token

	log.Debug(string(buf))
	body := ws._HttpPost(url, []byte(buf))
	wxError := new(wxmodel.ErrRsp)
	if err := json.Unmarshal(body, wxError); err != nil {
		return err
	}

	log.Debug(wxError)
	if wxError.Errcode != 0 {
		body := ws._HttpPost(url, []byte(buf))
		if err := json.Unmarshal(body, wxError); err != nil {
			return err
		}
		if wxError.Errcode != 0 {
			fmt.Println(wxError)
			return errors.New(wxError.Errmsg)
		}
	}
	return nil
}

// GetShareSignature 获取微信h5签名
func (ws *WeiXinSession) GetShareSignature(url string) (interface{}, error) {
	count := 0
ag:
	if _, err := ws.getTicket(); err != nil {
		ws.TokenReset()
		log.Error(err)
		//return nil, err
		if count > 10 {
			return nil, err
		}
		count++
		goto ag
	}
	s := new(wxmodel.SignatureResponse)
	s.NonceStr = randStr(16)
	s.Timestamp = time.Now().Unix()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ws.Ticket, s.NonceStr, s.Timestamp, url)
	// log.Debug("jsapi: ", str)
	h := sha1.New()
	h.Write([]byte(str))
	s.Signature = hex.EncodeToString(h.Sum(nil))
	return s, nil
}

// GetUserInfo 获取微信用户信息
func (ws *WeiXinSession) GetUserInfo(openid string) *wxmodel.UserInfoRsp {
	var err error
	if _, err = ws.getToken(); err != nil {
		log.Error(err)
		return nil
	}
	url := fmt.Sprintf(weixin_userinfo_url, ws.Token, openid)

	body := ws._HttpGet(url)
	log.Debug(string(body))
	rsp := new(wxmodel.UserInfoRsp)
	if err := json.Unmarshal(body, rsp); err != nil {
		log.Debug(err)
		return nil
	}

	return rsp
}

func fillEmpty(s string) string {
	if s == "" {
		return "killEmpty"
	}
	return s
}

const timeFormat = "20060102150405"

//GetAppUnifiedOrder UnifiedOrder
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=4_3
//计算APP签名
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
//body 商品描述
//detail  商品详情
//attach  自定义字段
func (ws *WeiXinSession) GetAppUnifiedOrder(openid, orderid, body, detail, notifyurl,
	attach string, fee int, expire time.Time) (*wxmodel.AppUnifierOrderRsp, error) {
	return pay.GetAppUnifiedOrder(ws.cfg.Appid, ws.cfg.MchId, ws.cfg.APIKey,
		openid, orderid, body, detail, notifyurl, attach,
		fee, expire, "WEB")
}

//Qrpay 扫码支付
func (ws *WeiXinSession) Qrpay(authcode, orderid, body, detail, attach string, fee int) (*wxmodel.QrPayRsp, error) {
	order := new(wxmodel.QrPay)
	order.Appid = ws.cfg.Appid
	order.Mchid = ws.cfg.MchId
	order.NonceStr = utils.GetRandomString(16)
	order.Deviceinfo = "qrcode"
	order.Body = fillEmpty(body)
	order.Detail = fillEmpty(detail)
	order.Attach = fillEmpty(attach)
	order.Outtradeno = fillEmpty(orderid)
	order.TotalFee = fee
	order.AuthCode = authcode
	stringTemp := fmt.Sprintf("appid=%s&attach=%s&auth_code=%s&body=%s&detail=%s&device_info=%s&mch_id=%s&nonce_str=%s&out_trade_no=%s&total_fee=%d&key=%s",
		ws.cfg.Appid,
		order.Attach,
		order.AuthCode,
		order.Body,
		order.Detail,
		order.Deviceinfo,
		order.Mchid,
		order.NonceStr,
		order.Outtradeno,
		order.TotalFee,
		ws.cfg.APIKey,
	)
	log.Info(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Trace("--app_pay_qrcode--", string(buf))
	r := ws._HttpPost(app_pay_qrcode, []byte(buf))
	log.Trace("--Rsp--", string(r))
	rsp := new(wxmodel.QrPayRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Trace("--AppUnifierOrderRsp--", rsp)
	return rsp, nil
}

//Orderquery 查询订单状态
func (ws *WeiXinSession) Orderquery(transactionid, outtradeno string) (*wxmodel.OrderqueryRsp, error) {
	var stringTemp string
	order := new(wxmodel.Orderquery)
	order.Appid = ws.cfg.Appid
	order.Mchid = ws.cfg.MchId
	order.TransactionID = transactionid
	order.OutTradeNo = outtradeno
	order.NonceStr = utils.GetRandomString(16)
	if order.OutTradeNo != "" {
		stringTemp = fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&out_trade_no=%s&key=%s",
			ws.cfg.Appid,
			order.Mchid,
			order.NonceStr,
			order.OutTradeNo,
			ws.cfg.APIKey,
		)
	} else {
		stringTemp = fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&transaction_id=%s&key=%s",
			ws.cfg.Appid,
			order.Mchid,
			order.NonceStr,
			order.TransactionID,
			ws.cfg.APIKey,
		)
	}
	log.Debug(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r := ws._HttpPost(app_pay_orderquery, []byte(buf))
	log.Trace("--Rsp--", string(r))
	rsp := new(wxmodel.OrderqueryRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debug("--OrderqueryRsp--", rsp)
	return rsp, nil
}

//Authcodetoopenid 付款码查询openid
func (ws *WeiXinSession) Authcodetoopenid(AuthCode string) (*wxmodel.AuthcodetoopenidRsp, error) {
	var stringTemp string
	order := new(wxmodel.Authcodetoopenid)
	order.Appid = ws.cfg.Appid
	order.Mchid = ws.cfg.MchId
	order.NonceStr = utils.GetRandomString(16)
	order.AuthCode = AuthCode
	stringTemp = fmt.Sprintf("appid=%s&auth_code=%s&mch_id=%s&nonce_str=%s&key=%s",
		ws.cfg.Appid,
		order.AuthCode,
		order.Mchid,
		order.NonceStr,
		ws.cfg.APIKey,
	)
	log.Debug(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r := ws._HttpPost(app_pay_authcodetoopenid, []byte(buf))
	log.Trace("--Rsp--", string(r))
	rsp := new(wxmodel.AuthcodetoopenidRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Trace("--AuthcodetoopenidRsp--", rsp)
	return rsp, nil
}

//Refund 退款
func (ws *WeiXinSession) Refund(transactionid, notifyurl, outrefundorder, desc string, fee, total int) (*wxmodel.RefundRsp, error) {
	return pay.Refund(ws.cfg.Appid, ws.cfg.MchId, ws.cfg.APIKey,
		transactionid, notifyurl, outrefundorder, desc, ws.cfg.KeyFile, ws.cfg.CertFile,
		ws.cfg.WXRoot, fee, total)
}

//WxConfig WxConfig配置
func (ws *WeiXinSession) WxConfig(url string) (conf wxmodel.WxConfig) {
	count := 0
ag:
	if _, err := ws.getTicket(); err != nil {
		ws.TokenReset()
		log.Error(err)
		//return nil, err
		if count > 10 {
			return conf
		}
		count++
		goto ag
	}

	conf.AppID = ws.cfg.Appid
	conf.Timestamp = time.Now().Unix()
	conf.NonceStr = randStr(16)
	sh := sha1.New()
	s := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ws.Ticket,
		conf.NonceStr,
		conf.Timestamp,
		url,
	)
	log.Trace(s)
	sh.Write([]byte(s))
	conf.Signature = fmt.Sprintf("%x", sh.Sum(nil))
	return conf
}

//Draw 提現
func (ws *WeiXinSession) Draw(orderno, openid, desc, ip string, amount int) (*wxmodel.UserDrawRsp, error) {
	return pay.Draw(ws.cfg.Appid, ws.cfg.MchId, ws.cfg.APIKey,
		orderno, openid, desc, ip, ws.cfg.KeyFile, ws.cfg.CertFile,
		ws.cfg.WXRoot, amount)
}

//Draw 提現
func (ws *WeiXinSession) DrawBank(orderno, card, name, bankname, desc string, amount int) (*wxmodel.UserDrawBankRsp, error) {
	return pay.DrawBank(ws.cfg.MchId, ws.cfg.APIKey,
		orderno, card, name, bankname, desc, ws.publicKey, ws.cfg.KeyFile, ws.cfg.CertFile,
		ws.cfg.WXRoot, amount)
}

//Refundquery 查询退款状态
func (ws *WeiXinSession) Refundquery(outrefundno string, offset int) (*wxmodel.RefundBody, error) {
	return pay.Refundquery(ws.cfg.Appid, ws.cfg.MchId, outrefundno, ws.cfg.APIKey, offset)
}

//GetMaterial 獲取素材列表
func (ws *WeiXinSession) GetMaterial() error {
	if _, err := ws.getToken(); err != nil {
		log.Error(err)
		return err
	}
	type req struct {
		Type   string `json:"type"`
		Offset int    `json:"offset"`
		Count  int    `json:"count"`
	}
	var r req
	r.Type = "image"
	r.Offset = 0
	r.Count = 20
	data, _ := json.Marshal(&r)
	url := fmt.Sprintf(GetMaterialServerURL, ws.Token)
	body := ws._HttpPost(url, data)
	log.Debug(string(body))
	return nil
}
