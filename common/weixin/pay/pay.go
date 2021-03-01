package pay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/hqbobo/frame/common/log"
	"github.com/hqbobo/frame/common/utils"
	wxmodel "github.com/hqbobo/frame/common/weixin/model"
	"strconv"
	"strings"
	"time"
)

func fillEmpty(s string) string {
	if s == "" {
		return "默认文本"
	}
	return s
}

// GetAppUnifiedOrder 获取UnifiedOrder
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=4_3
//计算APP签名
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
//body 商品描述
//detail  商品详情
//attach  自定义字段
func GetAppUnifiedOrder(appid, mchid, apikey, openid,
	orderid, body, detail, notifyurl, attach string, fee int, expire time.Time, device string) (*wxmodel.AppUnifierOrderRsp, error) {
	order := new(wxmodel.AppUnifierOrder)
	order.Appid = appid
	order.Mchid = mchid
	order.NonceStr = utils.GetRandomString(16)
	order.Deviceinfo = device
	order.Body = fillEmpty(body)
	order.Detail = fillEmpty(detail)
	order.Attach = fillEmpty(attach)
	order.Outtradeno = fillEmpty(orderid)
	order.TotalFee = fee
	order.Notifyurl = fillEmpty(notifyurl)
	order.Tradetype = "JSAPI"
	order.Openid = fillEmpty(openid)
	order.Timeexpire = expire.Format(wxmodel.TimeFormat)
	stringTemp := fmt.Sprintf("appid=%s&attach=%s&body=%s&detail=%s&device_info=%s&mch_id=%s&nonce_str=%s&notify_url=%s&openid=%s&out_trade_no=%s&time_expire=%s&total_fee=%d&trade_type=%s&key=%s",
		order.Appid,
		order.Attach,
		order.Body,
		order.Detail,
		order.Deviceinfo,
		order.Mchid,
		order.NonceStr,
		order.Notifyurl,
		order.Openid,
		order.Outtradeno,
		order.Timeexpire,
		order.TotalFee,
		order.Tradetype,
		apikey,
	)
	log.Info(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug("--GetAppUnifiedOrder--", string(buf))
	r := utils.HttpPost(wxmodel.PayUnifiedorder, []byte(buf))
	log.Debug("--Rsp--", string(r))
	rsp := new(wxmodel.AppUnifierOrderRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debug("--AppUnifierOrderRsp--", rsp)
	return rsp, nil
}

//Refundquery 查询单个退款状态
func Refundquery(appid, mchid, outrefundno, apikey string, offset int) (*wxmodel.RefundBody, error) {
	var stringTemp string
	order := new(wxmodel.Refundquery)
	order.Appid = appid
	order.Mchid = mchid
	order.OutRefundNo = outrefundno
	order.NonceStr = utils.GetRandomString(16)
	if offset != 0 {
		order.Offset = offset
		stringTemp = fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&offset=%d&out_refund_no=%s&key=%s",
			order.Appid,
			order.Mchid,
			order.NonceStr,
			order.Offset,
			order.OutRefundNo,
			apikey,
		)
	} else {
		stringTemp = fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&out_refund_no=%s&key=%s",
			order.Appid,
			order.Mchid,
			order.NonceStr,
			order.OutRefundNo,
			apikey,
		)
	}
	log.Debug(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r := utils.HttpPost(wxmodel.PayRefundqueryURL, []byte(buf))
	log.Debug("--Rsp--", string(r))
	rsp := new(wxmodel.RefundqueryRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	arr := make(map[string]string, 0)

	if err := xml.Unmarshal(r, arr); err != nil {
		return nil, err
	}
	ret := new(wxmodel.RefundBody)

	//查询第一个
	if rsp.RefundCount == 1 {
		ret.Returncode = rsp.Returncode
		ret.Returnmsg = rsp.Returnmsg
		ret.Appid = rsp.Appid
		ret.MchID = rsp.MchID
		ret.OutRefundNo = arr["out_refund_no_0"]
		ret.OutTradeNo = rsp.OutTradeNo
		ret.TotalFee = rsp.TotalFee
		ret.RefundID = arr["refund_id_0"]
		ret.RefundFee, err = strconv.ParseInt(arr["refund_fee_0"], 10, 64)
		if err != nil {
			log.Warn(err)
		}
		ret.RefundStatus = arr["refund_status_0"]
	}

	log.Debug("--Refundquery--", ret)
	return ret, nil
}

// Draw 提现
func Draw(appid, mchid, apikey, orderno, openid, desc, ip, key, cert, root string, amount int) (*wxmodel.UserDrawRsp, error) {
	order := new(wxmodel.UserDraw)
	order.Mchid = mchid
	order.NonceStr = utils.GetRandomString(16)
	order.PartnerTradeNo = orderno
	order.MchAppID = appid
	order.Openid = openid
	order.CheckName = "NO_CHECK"
	order.Amount = amount
	order.Desc = desc
	order.SpbillCreateIP = ip
	stringTemp := fmt.Sprintf("amount=%d&check_name=%s&desc=%s&mch_appid=%s&mchid=%s&nonce_str=%s&openid=%s&partner_trade_no=%s&spbill_create_ip=%s&key=%s",
		order.Amount,
		order.CheckName,
		order.Desc,
		order.MchAppID,
		order.Mchid,
		order.NonceStr,
		order.Openid,
		order.PartnerTradeNo,
		order.SpbillCreateIP,
		apikey,
	)
	log.Info(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug("--draw--", string(buf))
	r := utils.HttpsPost(wxmodel.PayDrawURL, key, cert, root, []byte(buf))
	log.Debug("--draw--", string(r))
	rsp := new(wxmodel.UserDrawRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debug("--draw--", rsp)
	return rsp, nil
}

//GetPublickKey 获取证书
func GetPublickKey(mchid, key, cert, root, apikey string) (string, error) {
	req := new(wxmodel.GetPublicKey)
	req.Mchid = mchid
	req.NonceStr = utils.GetRandomString(16)
	req.Signtype = "MD5"
	stringTemp := fmt.Sprintf("mch_id=%s&nonce_str=%s&sign_type=%s&key=%s",
		req.Mchid,
		req.NonceStr,
		req.Signtype,
		apikey,
	)
	log.Info(stringTemp)
	req.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(req)
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Debug("--GetPublickKey--", string(buf))
	r := utils.HttpsPost(wxmodel.PayGetPublicKey, key, cert, root, []byte(buf))
	rsp := new(wxmodel.GetPublicKeyRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return "", err
	}
	log.Debug("--GetPublickKey--", rsp)
	return rsp.PubKey, nil
}

const (
	PKCS1 PKCSType = 1 // 非java适用
	PKCS8 PKCSType = 2 // java适用
)

type PKCSType uint8

// rsaEncryptData RSA加密数据
//	t：PKCS1 或 PKCS8
//	originData：原始字符串byte数组
//	publicKey：公钥
func rsaEncryptData(t PKCSType, originData []byte, publicKey string) (cipherData string, err error) {
	var (
		key *rsa.PublicKey
	)

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("publicKey decode error")
	}

	switch t {
	case PKCS1:
		pkcs1Key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return "", err
		}
		key = pkcs1Key
	case PKCS8:
		pkcs8Key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", err
		}
		pk8, ok := pkcs8Key.(*rsa.PublicKey)
		if !ok {
			return "", errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		pkcs1Key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", err
		}
		pk1, ok := pkcs1Key.(*rsa.PublicKey)
		if !ok {
			return "", errors.New("publicKey parse error")
		}
		key = pk1
	}

	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, originData)
	if err != nil {
		return "", fmt.Errorf("xrsa.EncryptPKCS1v15：%w", err)
	}

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

//bankCode 转化银行名字到code
//https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=24_4&index=5
func bankCodeParse(name string) string {
	switch name {
	case "工商银行":
		return "1002"
	case "农业银行":
		return "1005"
	case "建设银行":
		return "1003"
	case "中国银行":
		return "1026"
	case "交通银行":
		return "1020"
	case "招商银行":
		return "1001"
	case "邮储银行":
		return "1066"
	case "民生银行":
		return "1006"
	case "平安银行":
		return "1010"
	case "中信银行":
		return "1021"
	case "浦发银行":
		return "1004"
	case "兴业银行":
		return "1009"
	case "光大银行":
		return "1022"
	case "广发银行":
		return "1027"
	case "华夏银行":
		return "1025"
	}
	return name
}

// DrawBank 提现到银行
func DrawBank(mchid, apikey, orderno, encbankno, enctruename, bankname, desc, publickey, key, cert, root string, amount int) (*wxmodel.UserDrawBankRsp, error) {
	order := new(wxmodel.UserDrawBank)
	order.Mchid = mchid
	order.NonceStr = utils.GetRandomString(16)
	order.PartnerTradeNo = orderno
	order.Amount = amount
	order.Desc = fillEmpty(desc)
	log.Warn(encbankno, "-", enctruename, "-", bankname)
	endata, err := rsaEncryptData(PKCS1, []byte(encbankno), publickey)
	if err != nil {
		return nil, err
	}
	order.EncBankNo = endata
	endata, err = rsaEncryptData(PKCS1, []byte(enctruename), publickey)
	if err != nil {
		return nil, err
	}
	order.EncTrueName = endata
	order.BankCode = bankCodeParse(bankname)
	stringTemp := fmt.Sprintf("amount=%d&bank_code=%s&desc=%s&enc_bank_no=%s&enc_true_name=%s&mch_id=%s&nonce_str=%s&partner_trade_no=%s&key=%s",
		order.Amount,
		order.BankCode,
		order.Desc,
		order.EncBankNo,
		order.EncTrueName,
		order.Mchid,
		order.NonceStr,
		order.PartnerTradeNo,
		apikey,
	)
	log.Trace(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Trace("--DrawBank--", string(buf))
	r := utils.HttpsPost(wxmodel.PayDrawBankURL, key, cert, root, []byte(buf))
	log.Trace("--DrawBank--", string(r))
	rsp := new(wxmodel.UserDrawBankRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Trace("--DrawBank--", rsp)
	return rsp, nil
}

// Refund 退款
func Refund(appid, mchid, apikey, transactionid, notifyurl, outrefundorder, desc, key, cert, root string, fee, total int) (*wxmodel.RefundRsp, error) {
	order := new(wxmodel.RefundReq)
	order.Appid = appid
	order.Mchid = mchid
	order.NonceStr = utils.GetRandomString(16)
	order.TotalFee = total
	order.RefundFee = fee
	//order.TransactionID = transactionid
	order.Notifyurl = fillEmpty(notifyurl)
	order.RefundDesc = desc
	order.OutRefundNo = outrefundorder
	order.OutTradeNo = transactionid
	stringTemp := fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&notify_url=%s&out_refund_no=%s&out_trade_no=%s&refund_desc=%s&refund_fee=%d&total_fee=%d&key=%s",
		order.Appid,
		order.Mchid,
		order.NonceStr,
		order.Notifyurl,
		order.OutRefundNo,
		order.OutTradeNo,
		order.RefundDesc,
		order.RefundFee,
		order.TotalFee,
		apikey,
	)
	log.Trace(stringTemp)
	order.Sign = strings.ToUpper(utils.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r := utils.HttpsPost(wxmodel.PayRefundURL, key, cert, root, []byte(buf))
	rsp := new(wxmodel.RefundRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Trace("--Refund--", rsp)
	return rsp, nil
}
