package xiaochengxu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/utils"
	wxmodel "github.com/hqbobo/frame/common/weixin/model"
	"github.com/hqbobo/frame/common/weixin/pay"
	"github.com/hqbobo/frame/common/weixin/xiaochengxu/model"
	minimodel "github.com/hqbobo/frame/common/weixin/xiaochengxu/model"
)

type MiniTpl struct {
	Tplid   string
	Context string
}

type weixntoken struct {
	Token   string
	TimeOut int64
}

type WeixinCacheInterface interface {
	Get(key string, data interface{}) bool
	Set(key string, data interface{}) bool
	Delete(key string) bool
}

// 微信小程序主结构
type WeiXinMiniSession struct {
	weixntoken
	cfg       wxmodel.WeixinCfg
	cache     WeixinCacheInterface
	publicKey string
}

// 获取配置
func (this *WeiXinMiniSession) CFG() wxmodel.WeixinCfg { return this.cfg }

func (sess *WeiXinMiniSession) stableToken() error {
	//优先缓存获取token
	var token weixntoken
	if sess.cache != nil {
		if sess.cache.Get(CacheTokenName+sess.cfg.Appid, &token) {
			//更新本地内存
			log.Debug("读取缓存token到本地:", token.Token, " 到期时间:", time.Now().Add(time.Second*time.Duration(token.TimeOut-time.Now().Unix())).String())
			sess.weixntoken = token
		}
	}
	log.Debug("time:", sess.TimeOut, "now:", time.Now().Unix())
	//还是超时或者没有
	if sess.Token == "" || sess.TimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp minimodel.AccessTokenRsp
		var req minimodel.AccessStableToken
		buf, err := json.Marshal(&req)
		if err != nil {
			log.Error(err)
			return err
		}
		r := utils.HttpsPostNoTLS(StableTokeUrl, buf)
		e := json.Unmarshal(r, &rsp)
		if e != nil {
			log.Warn(e)
			return e
		}
		sess.Token = rsp.Access_token
		sess.TimeOut = int64(rsp.Expires_in) + time.Now().Unix()
		log.Debug("APPID:["+sess.cfg.Appid+"]在线获取新的token :", rsp.Access_token, " 超时", time.Unix(sess.TimeOut, 0))

		if rsp.Errcode != 0 {
			return errors.New(rsp.ErrRsp.Errmsg)
		}
		//保存到缓存
		if sess.cache != nil {
			sess.cache.Set(CacheTokenName+sess.cfg.Appid, &sess.weixntoken)
		}
	}
	return nil
}

// 在线获取access_token
func (this *WeiXinMiniSession) reloadToken() error {
	//优先缓存获取token
	var token weixntoken
	if this.cache != nil {
		if this.cache.Get(CacheTokenName+this.cfg.Appid, &token) {
			//更新本地内存
			log.Debug("读取缓存token到本地:", token.Token, " 到期时间:", time.Now().Add(time.Second*time.Duration(token.TimeOut-time.Now().Unix())).String())
			this.weixntoken = token
		}
	}
	log.Debug("time:", this.TimeOut, "now:", time.Now().Unix())
	//还是超时或者没有
	if this.Token == "" || this.TimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp minimodel.AccessTokenRsp
		r, e := utils.HttpsGet(fmt.Sprintf(TokenURL, this.cfg.Appid, this.cfg.Appsecret))
		if e != nil {
			log.Warn(e)
			return e
		}
		e = json.Unmarshal(r, &rsp)
		if e != nil {
			log.Warn(e)
			return e
		}
		this.Token = rsp.Access_token
		this.TimeOut = int64(rsp.Expires_in) + time.Now().Unix()
		log.Debug("APPID:["+this.cfg.Appid+"]在线获取新的token :", rsp.Access_token, " 超时", time.Unix(this.TimeOut, 0))

		if rsp.Errcode != 0 {
			return errors.New(rsp.ErrRsp.Errmsg)
		}
		//保存到缓存
		if this.cache != nil {
			this.cache.Set(CacheTokenName+this.cfg.Appid, &this.weixntoken)
		}
	}

	return nil
}

func (this *WeiXinMiniSession) GetAccessToken() (string, error) {
	return this.getToken()
}

// 本地获取access_token
func (this *WeiXinMiniSession) froceReloadToken() (string, error) {
	//在线更新token
	var rsp minimodel.AccessTokenRsp
	r, e := utils.HttpsGet(fmt.Sprintf(TokenURL, this.cfg.Appid, this.cfg.Appsecret))
	if e != nil {
		log.Warn(e)
		return "", e
	}
	e = json.Unmarshal(r, &rsp)
	if e != nil {
		log.Warn(e)
		return "", e
	}
	this.Token = rsp.Access_token
	this.TimeOut = int64(rsp.Expires_in-100) + time.Now().Unix()
	log.Debug("APPID:["+this.cfg.Appid+"]在线获取新的token :", rsp.Access_token, " 超时", time.Unix(this.TimeOut, 0))

	if rsp.Errcode != 0 {
		return "", errors.New(rsp.ErrRsp.Errmsg)
	}
	//保存到缓存
	if this.cache != nil {
		this.cache.Set(CacheTokenName+this.cfg.Appid, &this.weixntoken)
	}
	return this.Token, nil

}

// 本地获取access_token
func (this *WeiXinMiniSession) getToken() (string, error) {
	//检查token是否已经获取和超时
	if this.Token == "" || this.TimeOut <= time.Now().Unix() {
		if e := this.reloadToken(); e != nil {
			return "", e
		}

	}
	return this.Token, nil
}

func (this *WeiXinMiniSession) SetCacheInterface(cache WeixinCacheInterface) {
	this.cache = cache

}

func (this *WeiXinMiniSession) Run() {
	var e error
	if _, e := this.getToken(); e != nil {
		log.Debug(e)
	}
	this.publicKey, e = pay.GetPublickKey(this.cfg.MchId, this.cfg.KeyFile, this.cfg.CertFile, this.cfg.WXRoot, this.cfg.APIKey)
	if e != nil {
		log.Warn(e)
	}
}

func NewSession(cfg wxmodel.WeixinCfg) *WeiXinMiniSession {
	O := new(WeiXinMiniSession)
	O.cfg = cfg
	return O
}

type WxMiniBody struct {
	OpenId      string `json:"openId"`
	Session_key string `json:"session_key"`
	UnionId     string `json:"unionid"`
}

func (sess *WeiXinMiniSession) weixinMiniProgramLogin(code string) (*WxMiniBody, error) {
	var err error
	if _, err = sess.getToken(); err != nil {
		log.Error(err)
		return nil, err
	}
	body, err := utils.HttpsGet("https://api.weixin.qq.com/sns/jscode2session?appid=" + sess.cfg.Appid + "&secret=" + sess.cfg.Appsecret + "&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		log.Debug(err, "--------->获取openid和sessionKey失败")
		return nil, err
	}
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	pre := string(body)
	if strings.Contains(pre, "errcode") {
		log.Warn(string(body))
		return nil, errors.New(string(body))
	}
	log.Debug("wx resp", string(body))
	o := new(WxMiniBody)
	err = json.Unmarshal(body, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (sess *WeiXinMiniSession) weixinMiniProgramGetPhone(code string) (*model.UserPhoneByCode, error) {
	var err error
	if _, err = sess.getToken(); err != nil {
		log.Error(err)
		return nil, err
	}

	body := utils.HttpPost("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="+sess.Token,
		[]byte("{\"code\": \""+code+"\"}"))
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	log.Debug("wx resp", string(body))
	o := new(model.UserPhoneByCode)
	err = json.Unmarshal(body, o)
	if err != nil {
		return o, err
	}
	if o.Errcode != 0 {
		return o, errors.New(o.Errmsg)
	}
	return o, nil
}
func (sess *WeiXinMiniSession) WxLogin(encryptedData, iv, code string) (*minimodel.UserWX, error) {
	session, err := sess.weixinMiniProgramLogin(code)
	if err != nil {
		return nil, err
	}
	log.Debug("iv:", iv, " code:", code, " session_key:", session.Session_key)
	pc := utils.WxBizDataCrypt{AppID: sess.cfg.Appid, SessionKey: session.Session_key}
	str, err := pc.Decrypt(encryptedData, iv, true) //第三个参数解释： 需要返回 JSON 数据类型时 使用 true, 需要返回 map 数据类型时 使用 false
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("%s", str.(string))
	user := new(minimodel.UserWX)
	json.Unmarshal([]byte(str.(string)), user)
	user.OpenId = session.OpenId
	user.Unionid = session.UnionId
	log.Debug(user)
	return user, nil
}

func (sess *WeiXinMiniSession) WxPhone(encryptedData, iv, code, phonecode string) (*minimodel.UserPhone, error) {
	session, err := sess.weixinMiniProgramLogin(code)
	if err != nil {
		return nil, err
	}
	log.Debug("iv:", iv, " code:", code, " session_key:", session.Session_key)
	try := 1
	//重试1次
ag:
	puser, err := sess.weixinMiniProgramGetPhone(phonecode)
	if err != nil {
		if puser != nil && puser.Errcode == 40001 && try > 0 {
			log.Warn("莫名过期重新获取:", sess.Token)
			sess.froceReloadToken()
			try--
			goto ag
		}
		return nil, err
	}
	user := new(minimodel.UserPhone)
	user.PhoneNumber = puser.Phone.PhoneNumber
	user.PurePhoneNumber = puser.Phone.PurePhoneNumber
	user.Watermark = puser.Phone.Watermark
	user.OpenId = session.OpenId
	user.Unionid = session.UnionId
	log.Debug(user)
	return user, nil
}

func (this *WeiXinMiniSession) QrcodeGetUnlimited(sense, page, ver string, c bool) (e error, t string, d []byte) {
	var err error
	if _, err = this.getToken(); err != nil {
		log.Error(err)
		return err, t, d
	}
	wtc := new(minimodel.QrcodeGetUnlimitedReq)
	wtc.Scene = sense
	wtc.Checkpath = c
	wtc.Page = page
	wtc.Envversion = ver
	buf, _ := json.Marshal(wtc)
	url := getwxacodeunlimit_url + this.Token
	d = utils.HttpPost(url, []byte(buf))
	wxError := new(minimodel.ErrRsp)
	if err := json.Unmarshal(d, wxError); err == nil {
		log.Debug(wxError)
		return errors.New(wxError.Errmsg), t, d
	}
	return nil, "image/jpeg", d
}

func (this *WeiXinMiniSession) SubscribeMessage(openid, templateid, page string, data interface{}) error {
	var err error
	if _, err = this.getToken(); err != nil {
		log.Error(err)
		return err
	}
	wtc := new(minimodel.SubscribeMessage)
	wtc.Template_id = templateid
	wtc.Touser = openid
	wtc.Data = data
	wtc.Miniprogram_state = "developer"
	wtc.Page = page
	buf, _ := json.Marshal(wtc)
	url := mini_subscribe_url + "?access_token=" + this.Token

	log.Debug(string(buf))
	body := utils.HttpPost(url, []byte(buf))
	wxError := new(minimodel.ErrRsp)
	if err := json.Unmarshal(body, wxError); err != nil {
		return err
	}

	log.Debug(wxError)
	if wxError.Errcode != 0 {
		body := utils.HttpPost(url, []byte(buf))
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

func computeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	//	fmt.Println(h.Sum(nil))
	sha := hex.EncodeToString(h.Sum(nil))
	//	fmt.Println(sha)

	//	hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func fillEmpty(s string) string {
	if s == "" {
		return "killEmpty"
	}
	return s
}

// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=4_3
// 计算APP签名
// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
// body 商品描述
// detail  商品详情
// attach  自定义字段
func (this *WeiXinMiniSession) GetAppUnifiedOrder(openid, orderid, body, detail, notifyurl,
	attach string, fee int, expire time.Time) (*wxmodel.AppUnifierOrderRsp, error) {
	return pay.GetAppUnifiedOrder(this.cfg.Appid, this.cfg.MchId, this.cfg.APIKey,
		openid, orderid, body, detail, notifyurl, attach,
		fee, expire, "mini")
}

func (this *WeiXinMiniSession) Refund(transaction_id, notifyurl, out_refund_order, desc string, fee, total int) (*wxmodel.RefundRsp, error) {
	return pay.Refund(this.cfg.Appid, this.cfg.MchId, this.cfg.APIKey,
		transaction_id, notifyurl, out_refund_order, desc, this.cfg.KeyFile, this.cfg.CertFile,
		this.cfg.WXRoot, fee, total)
}

func (this *WeiXinMiniSession) Draw(orderno, openid, desc, ip string, amount int) (*wxmodel.UserDrawRsp, error) {
	return pay.Draw(this.cfg.Appid, this.cfg.MchId, this.cfg.APIKey,
		orderno, openid, desc, ip, this.cfg.KeyFile, this.cfg.CertFile,
		this.cfg.WXRoot, amount)
}

// Refundquery 查询退款状态
func (this *WeiXinMiniSession) Refundquery(out_refund_no string, offset int) (*wxmodel.RefundBody, error) {
	return pay.Refundquery(this.cfg.Appid, this.cfg.MchId, out_refund_no, this.cfg.APIKey, offset)
}

// Draw 提現
func (this *WeiXinMiniSession) DrawBank(orderno, card, name, bankname, desc string, amount int) (*wxmodel.UserDrawBankRsp, error) {
	return pay.DrawBank(this.cfg.MchId, this.cfg.APIKey,
		orderno, card, name, bankname, desc, this.publicKey, this.cfg.KeyFile, this.cfg.CertFile,
		this.cfg.WXRoot, amount)
}
