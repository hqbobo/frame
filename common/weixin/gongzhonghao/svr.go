package gongzhonghao

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/weixin/gongzhonghao/handle"
	wxmodel "github.com/hqbobo/frame/common/weixin/model"
)

var cfgs []wxmodel.WeixinCfg
var gSess *WeiXinSession

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("fromweixin"))
}
func index(w http.ResponseWriter, r *http.Request) {
	var s []string
	var str string
	var token string
	log.Trace("微信回调:", r.URL.String())
	appid := string([]byte(r.URL.Path)[strings.LastIndex(r.URL.Path, "/")+1:])
	//log.Debug(appid)

	for _, v := range cfgs {
		if v.Appid == appid {
			token = v.Token
			//log.Warn("token是", token)
		}
	}
	r.ParseForm() //解析参数, 默认是不会解析的
	s = append(s, token)
	signature := r.URL.Query().Get("signature")
	s = append(s, r.URL.Query().Get("timestamp"))
	s = append(s, r.URL.Query().Get("nonce"))
	echostr := r.URL.Query().Get("echostr")
	sort.Strings(s)
	for _, v := range s {
		str += v
	}
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(str))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	hex.EncodeToString(bs)

	//校验是否是微信服务器发来的消息
	if strings.Compare(signature, hex.EncodeToString(bs)) == 0 {
		//Get代表服务器做校验
		if strings.Compare(strings.ToLower(r.Method), "get") == 0 && (len(echostr) > 0) {
			w.Write([]byte(echostr))
		} else if strings.Compare(strings.ToLower(r.Method), "post") == 0 {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Error(err)
			}
			log.Trace("\n微信消息:\n", string(body))
			//解析头部
			var msg wxmodel.WeixinBase
			if e := xml.Unmarshal(body, &msg); e != nil {
				log.Warn(e)
			}
			rsp := handle.Entrance(appid, msg, body)

			r, e := xml.Marshal(rsp)
			if e != nil {
				log.Warn(string(r), " \n err:", e)
			}
			log.Trace("\n返回微信消息:\n", string(r))
			w.Write(r)

		} else {
			w.Write(nil)
		}
	}

}

func svrInit(sess *WeiXinSession, cb handle.WeixinMsgHandleInterface, event handle.WeixinEventHandleInterface) {
	gSess = sess
	handle.Init(cb, event)
	http.HandleFunc("/", index) //设置访问的路由
	//后台启动
	go func() {
		log.Info("微信启动在", sess.cfg.Addr)
		err := http.ListenAndServe(sess.cfg.Addr, nil) //设置监听的端口
		if err != nil {
			log.Error("ListenAndServe: ", err)
		}
	}()

}

func svrsInit(cfgs []wxmodel.WeixinCfg, cb handle.WeixinMsgHandleInterface, event handle.WeixinEventHandleInterface) {
	handle.Init(cb, event)
	http.HandleFunc("/", test) //设置访问的路由
	for _, v := range cfgs {
		http.HandleFunc("/"+v.Appid, index) //设置访问的路由
		log.Infof("注册微信路由:/" + v.Appid)
	}
	http.HandleFunc("/test", test) //设置访问的路由

	//后台启动
	go func() {
		log.Info("微信启动在80")
		err := http.ListenAndServe(":80", nil) //设置监听的端口
		if err != nil {
			log.Error("ListenAndServe: ", err)
		}
	}()

}
