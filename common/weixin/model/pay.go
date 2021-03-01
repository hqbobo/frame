package model

import "encoding/xml"

//PayOrderqueryURL 查询订单
const PayOrderqueryURL = "https://api.mch.weixin.qq.com/pay/orderquery"

//PayRefundqueryURL 查询退款
const PayRefundqueryURL = "https://api.mch.weixin.qq.com/pay/refundquery"

//PayDrawURL 提现URL
const PayDrawURL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"

//PayDrawBankURL 提现银行URL
const PayDrawBankURL = "https://api.mch.weixin.qq.com/mmpaysptrans/pay_bank"

//PayGetPublicKey 获取RSA加密公钥API
const PayGetPublicKey = "https://fraud.mch.weixin.qq.com/risk/getpublickey"

//PayRefundURL 退款url
const PayRefundURL = "https://api.mch.weixin.qq.com/secapi/pay/refund"

//PayUnifiedorder Unifiedorderurl
const PayUnifiedorder = "https://api.mch.weixin.qq.com/pay/unifiedorder"

//TimeFormat 时间格式
const TimeFormat = "20060102150405"

//RefundReq 退款请求
type RefundReq struct {
	Appid         string `xml:"appid,omitempty"`
	Mchid         string `xml:"mch_id,omitempty"`
	NonceStr      string `xml:"nonce_str,omitempty"`
	Sign          string `xml:"sign,omitempty"`
	Signtype      string `xml:"sign_type,omitempty"`
	OutTradeNo    string `xml:"out_trade_no,omitempty"`
	//TransactionID string `xml:"transaction_id,omitempty"`
	OutRefundNo   string `xml:"out_refund_no,omitempty"`
	TotalFee      int    `xml:"total_fee,omitempty"`
	RefundFee     int    `xml:"refund_fee,omitempty"`
	RefundFeeType string `xml:"refund_fee_type,omitempty"`
	RefundDesc    string `xml:"refund_desc,omitempty"`
	Notifyurl     string `xml:"notify_url,omitempty"`
}

//RefundRsp 退款返回
type RefundRsp struct {
	XMLName       xml.Name    `xml:"xml,omitempty"`
	Appid         CdataString `xml:"appid,omitempty"`
	Mchid         CdataString `xml:"mch_id,omitempty"`
	NonceStr      CdataString `xml:"nonce_str,omitempty"`
	Sign          CdataString `xml:"sign,omitempty"`
	Resultcode    CdataString `xml:"result_code,omitempty"`
	Errcode       CdataString `xml:"err_code,omitempty"`
	Errcodedes    CdataString `xml:"err_code_des,omitempty"`
	TransactionID CdataString `xml:"transaction_id,omitempty"`
	OutTradeNo    CdataString `xml:"out_trade_no,omitempty"`
	OutRefundNo   CdataString `xml:"out_refund_no,omitempty"`
	RefundID      CdataString `xml:"refund_id,omitempty"`
	RefundFee     int         `xml:"refund_fee,omitempty"`
	TotalFee      int         `xml:"total_fee,omitempty"`
	CashFee       int         `xml:"cash_fee,omitempty"`
	ReturnCode
}

//AppUnifierOrder UnifierOrder请求
type AppUnifierOrder struct {
	XMLName        xml.Name `xml:"xml,omitempty"`
	Appid          string   `xml:"appid,omitempty"`
	Mchid          string   `xml:"mch_id,omitempty"`
	Deviceinfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	Signtype       string   `xml:"sign_type,omitempty"`
	Body           string   `xml:"body,,omitempty"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"`
	Outtradeno     string   `xml:"out_trade_no,omitempty"`
	Openid         string   `xml:"openid,omitempty"`
	FeeType        string   `xml:"fee_type,omitempty"`
	TotalFee       int      `xml:"total_fee,omitempty"`
	Spbillcreateip string   `xml:"spbill_create_ip,omitempty"`
	Timestart      string   `xml:"time_start,omitempty"`
	Timeexpire     string   `xml:"time_expire,omitempty"`
	Goodstag       string   `xml:"goods_tag,omitempty"`
	Notifyurl      string   `xml:"notify_url,omitempty"`
	Tradetype      string   `xml:"trade_type,omitempty"`
	Limitpay       string   `xml:"limit_pay,omitempty"`
	Sceneinfo      string   `xml:"scene_info,omitempty"`
}

//ReturnCode 请求返回
type ReturnCode struct {
	Returncode CdataString `xml:"return_code,omitempty"`
	Returnmsg  CdataString `xml:"return_msg,omitempty"`
}

//AppUnifierOrderRsp UnifierOrder请求返回
type AppUnifierOrderRsp struct {
	XMLName    xml.Name    `xml:"xml,omitempty"`
	Appid      CdataString `xml:"appid,omitempty"`
	Mchid      CdataString `xml:"mch_id,omitempty"`
	Deviceinfo CdataString `xml:"device_info,omitempty"`
	NonceStr   CdataString `xml:"nonce_str,omitempty"`
	Sign       CdataString `xml:"sign,omitempty"`
	Resultcode CdataString `xml:"result_code,omitempty"`
	Errcode    CdataString `xml:"err_code,omitempty"`
	Errcodedes CdataString `xml:"err_code_des,omitempty"`
	Tradetype  CdataString `xml:"trade_type,omitempty"`
	Prepayid   CdataString `xml:"prepay_id,omitempty"`
	ReturnCode
}

//NotifyBody 支付回调
//！！！如果微信修改返回字段，验签会失败！！！
//result_code Return_code  both SUCCESS
//代金券ID	coupon_id_$n	否	String(20)	10000	代金券ID,$n为下标，从0开始编号
//单个代金券支付金额	coupon_fee_$n	否	Int	100	单个代金券支付金额,$n为下标，从0开始编号
type NotifyBody struct {
	XMLName            xml.Name `xml:"xml,omitempty"`
	Returncode         string   `xml:"return_code,omitempty"`
	Appid              string   `xml:"appid,omitempty"`
	BankType           string   `xml:"bank_type,omitempty"`
	DeviceInfo         string   `xml:"device_info,omitempty"`
	CashFee            string   `xml:"cash_fee,omitempty"`
	FeeType            string   `xml:"fee_type,omitempty"`
	IsSubscribe        string   `xml:"is_subscribe,omitempty"`
	MchID              string   `xml:"mch_id,omitempty"`
	NonceStr           string   `xml:"nonce_str,omitempty"`
	Openid             string   `xml:"openid,omitempty"`
	OutTradeNo         string   `xml:"out_trade_no,omitempty"`
	ResultCode         string   `xml:"result_code,omitempty"`
	Sign               string   `xml:"sign,omitempty"`
	TimeEnd            string   `xml:"time_end,omitempty"`
	TotalFee           int64    `xml:"total_fee,omitempty"`
	TradeType          string   `xml:"trade_type,omitempty"`
	TransactionID      string   `xml:"transaction_id,omitempty"`
	SignType           string   `xml:"sign_type,omitempty"`
	SettlementTotalFee int64    `xml:"settlement_total_fee,omitempty"`
	ErrCode            string   `xml:"err_code,omitempty"`
	ErrCodeDes         string   `xml:"err_code_des,omitempty"`
	CashFeeType        string   `xml:"cash_fee_type,omitempty"`
	CouponFee          string   `xml:"coupon_fee,omitempty"`
	CouponCount        string   `xml:"coupon_count,omitempty"`
	Attach             string   `xml:"attach,omitempty"`
	CouponID0          string   `xml:"coupon_id_0,omitempty"`
	CouponID1          string   `xml:"coupon_id_1,omitempty"`
	CouponType0        string   `xml:"coupon_type_0,omitempty"`
	CouponType1        string   `xml:"coupon_type_1,omitempty"`
	CouponFee0         string   `xml:"coupon_fee_0,omitempty"`
	CouponFee1         string   `xml:"coupon_fee_1,omitempty"`
}

//RefundBody 退款回执
//！！！如果微信修改返回字段，验签会失败！！！
//https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_16&index=10#menu1
type RefundBody struct {
	XMLName             xml.Name `xml:"xml,omitempty"`
	Appid               string   `xml:"appid,omitempty"`
	MchID               string   `xml:"mch_id,omitempty"`
	NonceStr            string   `xml:"nonce_str,omitempty"`
	ReqInfo             string   `xml:"req_info,omitempty"`
	TransactionID       string   `xml:"transaction_id,omitempty"` //微信订单号
	OutTradeNo          string   `xml:"out_trade_no,omitempty"`   //商户系统内部的订单号
	RefundID            string   `xml:"refund_id,omitempty"`      //微信退款单号
	OutRefundNo         string   `xml:"out_refund_no,omitempty"`  //商户退款单号
	TotalFee            int64    `xml:"total_fee,omitempty"`
	SettlementTotalFee  int64    `xml:"settlement_total_fee,omitempty"`
	RefundFee           int64    `xml:"refund_fee,omitempty"`
	SettlementRefundFee int64    `xml:"settlement_refund_fee,omitempty"`
	RefundStatus        string   `xml:"refund_status,omitempty"` //SUCCESS-退款成功 CHANGE-退款异常REFUNDCLOSE—退款关闭
	ReturnCode
}

//RefundNotifyBody 基础数据
//！！！如果微信修改返回字段，验签会失败！！！
//https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_16&index=10#menu1
type RefundNotifyBody struct {
	XMLName  xml.Name `xml:"xml,omitempty"`
	Appid    string   `xml:"appid,omitempty"`
	MchID    string   `xml:"mch_id,omitempty"`
	NonceStr string   `xml:"nonce_str,omitempty"`
	ReqInfo  string   `xml:"req_info,omitempty"` //加密数据
	ReturnCode
}

//RefundNotifyRoot 加密数据
type RefundNotifyRoot struct {
	TransactionID       string `xml:"transaction_id,omitempty"` //微信订单号
	OutTradeNo          string `xml:"out_trade_no,omitempty"`   //商户系统内部的订单号
	RefundID            string `xml:"refund_id,omitempty"`      //微信退款单号
	OutRefundNo         string `xml:"out_refund_no,omitempty"`  //商户退款单号
	TotalFee            int64  `xml:"total_fee,omitempty"`
	SettlementTotalFee  int64  `xml:"settlement_total_fee,omitempty"`
	RefundFee           int64  `xml:"refund_fee,omitempty"`
	SettlementRefundFee int64  `xml:"settlement_refund_fee,omitempty"`
	RefundStatus        string `xml:"refund_status,omitempty"` //SUCCESS-退款成功 CHANGE-退款异常REFUNDCLOSE—退款关闭
	SuccessTime         string `xml:"success_time,omitempty"`
	RefundRequestSource string `xml:"refund_request_source,omitempty"`
}

//UserDraw 提现请求
type UserDraw struct {
	MchAppID       string `xml:"mch_appid"`
	Mchid          string `xml:"mchid"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	Openid         string `xml:"openid"`
	CheckName      string `xml:"check_name"`
	Amount         int    `xml:"amount"`
	Desc           string `xml:"desc"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
}

//UserDrawRsp 提现返回
type UserDrawRsp struct {
	XMLName        xml.Name    `xml:"xml,omitempty"`
	AppID          CdataString `xml:"mch_appid,omitempty"`
	MchID          CdataString `xml:"mchid,omitempty"`
	DeviceInfo     string      `xml:"device_info,omitempty"`
	NonceStr       CdataString `xml:"nonce_str,omitempty"`
	Sign           CdataString `xml:"sign,omitempty"`
	Resultcode     CdataString `xml:"result_code,omitempty"`
	Errcode        CdataString `xml:"err_code,omitempty"`
	Errcodedes     CdataString `xml:"err_code_des,omitempty"`
	PartnerTradeNo CdataString `xml:"partner_trade_no,omitempty"` //商户订单号，需保持历史全局唯一性(只能是字母或者数字，不能包含有其它字符)
	PaymentNo      CdataString `xml:"payment_no,omitempty"`       //企业付款成功，返回的微信付款单号
	PaymentTime    CdataString `xml:"payment_time,omitempty"`     //企业付款成功时间
	ReturnCode
}

//UserDrawBank 提现银行请求
type UserDrawBank struct {
	Mchid          string `xml:"mch_id"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	EncBankNo      string `xml:"enc_bank_no"`
	EncTrueName    string `xml:"enc_true_name"`
	BankCode       string `xml:"bank_code"`
	Amount         int    `xml:"amount"`
	Desc           string `xml:"desc"`
}

//UserDrawBankRsp 提现银行返回
type UserDrawBankRsp struct {
	XMLName        xml.Name    `xml:"xml,omitempty"`
	Resultcode     CdataString `xml:"result_code,omitempty"`
	Errcode        CdataString `xml:"err_code,omitempty"`
	Errcodedes     CdataString `xml:"err_code_des,omitempty"`
	MchID          CdataString `xml:"mchid,omitempty"`
	PartnerTradeNo CdataString `xml:"partner_trade_no,omitempty"` //商户订单号，需保持历史全局唯一性(只能是字母或者数字，不能包含有其它字符)
	Amount         int         `xml:"amount"`
	NonceStr       string      `xml:"nonce_str"`
	Sign           CdataString `xml:"sign,omitempty"`
	//以下字段在return_code 和result_code都为SUCCESS的时候有返回
	PaymentNo CdataString `xml:"payment_no,omitempty"` //企业付款成功，返回的微信付款单号
	CmmsAmt   int         `xml:"cmms_amt"`             //手续费金额 RMB：分
	ReturnCode
}

// WxConfig WxConfig配置
type WxConfig struct {
	AppID     string
	Timestamp int64
	NonceStr  string
	Signature string
}

// QrPay 二维码支付请求
type QrPay struct {
	XMLName        xml.Name `xml:"xml,omitempty"`
	Appid          string   `xml:"appid,omitempty"`
	Mchid          string   `xml:"mch_id,omitempty"`
	Deviceinfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	Signtype       string   `xml:"sign_type,omitempty"`
	Body           string   `xml:"body,,omitempty"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"` //商家数据包，原样返回 String(128)
	Outtradeno     string   `xml:"out_trade_no,omitempty"`
	TotalFee       int      `xml:"total_fee,omitempty"`
	FeeType        string   `xml:"fee_type,omitempty"`
	Spbillcreateip string   `xml:"spbill_create_ip,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty"`
	Limitpay       string   `xml:"limit_pay,omitempty"` //no_credit--指定不能使用信用卡支付
	TimeStart      string   `xml:"time_start,omitempty"`
	Timeexpire     string   `xml:"time_expire,omitempty"`
	Receipt        string   `xml:"receipt,omitempty"`
	AuthCode       string   `xml:"auth_code,omitempty"` //付款码
	Sceneinfo      string   `xml:"scene_info,omitempty"`
}

// QrPayRsp 二维码支付返回
type QrPayRsp struct {
	XMLName            xml.Name `xml:"xml,omitempty"`
	Appid              string   `xml:"appid,omitempty"`
	MchID              string   `xml:"mch_id,omitempty"`
	DeviceInfo         string   `xml:"device_info,omitempty"`
	NonceStr           string   `xml:"nonce_str,omitempty"`
	Sign               string   `xml:"sign,omitempty"`
	ResultCode         string   `xml:"result_code,omitempty"`
	ErrCode            string   `xml:"err_code,omitempty"`
	ErrCodeDes         string   `xml:"err_code_des,omitempty"`
	Openid             string   `xml:"openid,omitempty"`
	IsSubscribe        string   `xml:"is_subscribe,omitempty"`
	TradeType          string   `xml:"trade_type,omitempty"`
	BankType           string   `xml:"bank_type,omitempty"`
	FeeType            string   `xml:"fee_type,omitempty"`
	TotalFee           int      `xml:"total_fee,omitempty"`
	SettlementTotalFee int      `xml:"settlement_total_fee,omitempty"`
	CouponFee          string   `xml:"coupon_fee,omitempty"`
	CashFeeType        string   `xml:"cash_fee_type,omitempty"`
	CashFee            string   `xml:"cash_fee,omitempty"`
	TransactionID      string   `xml:"transaction_id,omitempty"`
	OutTradeNo         string   `xml:"out_trade_no,omitempty"`
	Attach             string   `xml:"attach,omitempty"`
	TimeEnd            string   `xml:"time_end,omitempty"`
	PromotionDetail    string   `xml:"promotion_detail,omitempty"`
	ReturnCode
}

// Orderquery 订单查询
type Orderquery struct {
	XMLName       xml.Name `xml:"xml,omitempty"`
	Appid         string   `xml:"appid,omitempty"`
	Mchid         string   `xml:"mch_id,omitempty"`
	TransactionID string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	Signtype      string   `xml:"sign_type,omitempty"`
}

// OrderqueryRsp 订单查询返回
type OrderqueryRsp struct {
	XMLName     xml.Name `xml:"xml,omitempty"`
	Appid       string   `xml:"appid,omitempty"`
	MchID       string   `xml:"mch_id,omitempty"`
	NonceStr    string   `xml:"nonce_str,omitempty"`
	Sign        string   `xml:"sign,omitempty"`
	ResultCode  string   `xml:"result_code,omitempty"`
	ErrCode     string   `xml:"err_code,omitempty"`
	ErrCodeDes  string   `xml:"err_code_des,omitempty"`
	DeviceInfo  string   `xml:"device_info,omitempty"`
	Openid      string   `xml:"openid,omitempty"`
	IsSubscribe string   `xml:"is_subscribe,omitempty"`
	TradeType   string   `xml:"trade_type,omitempty"`
	//Trade_state: SUCCESS—支付成功
	//REFUND—转入退款
	//NOTPAY—未支付
	//CLOSED—已关闭
	//REVOKED—已撤销（付款码支付）
	//USERPAYING--用户支付中（付款码支付）
	//PAYERROR--支付失败(其他原因，如银行返回失败)
	TradeState         string `xml:"trade_state,omitempty"`
	BankType           string `xml:"bank_type,omitempty"`
	FeeType            string `xml:"fee_type,omitempty"`
	TotalFee           int    `xml:"total_fee,omitempty"`
	SettlementTotalFee int    `xml:"settlement_total_fee,omitempty"`
	CouponFee          string `xml:"coupon_fee,omitempty"`
	CashFeeType        string `xml:"cash_fee_type,omitempty"`
	CashFee            string `xml:"cash_fee,omitempty"`
	TransactionID      string `xml:"transaction_id,omitempty"`
	OutTradeNo         string `xml:"out_trade_no,omitempty"`
	Attach             string `xml:"attach,omitempty"`
	TimeEnd            string `xml:"time_end,omitempty"`
	PromotionDetail    string `xml:"promotion_detail,omitempty"`
	TradeStateDesc     string `xml:"trade_state_desc,omitempty"` //对当前查询订单状态的描述和下一步操作的指引
	ReturnCode
}

// Authcodetoopenid 付款码获取openid
type Authcodetoopenid struct {
	XMLName  xml.Name `xml:"xml,omitempty"`
	Appid    string   `xml:"appid,omitempty"`
	Mchid    string   `xml:"mch_id,omitempty"`
	NonceStr string   `xml:"nonce_str,omitempty"`
	Sign     string   `xml:"sign,omitempty"`
	AuthCode string   `xml:"auth_code,omitempty"` //付款码
}

// AuthcodetoopenidRsp 付款码获取openid返回
type AuthcodetoopenidRsp struct {
	XMLName    xml.Name `xml:"xml,omitempty"`
	Appid      string   `xml:"appid,omitempty"`
	MchID      string   `xml:"mch_id,omitempty"`
	NonceStr   string   `xml:"nonce_str,omitempty"`
	Sign       string   `xml:"sign,omitempty"`
	ResultCode string   `xml:"result_code,omitempty"`
	ErrCode    string   `xml:"err_code,omitempty"`
	Openid     string   `xml:"openid,omitempty"`
	ReturnCode
}

//Refundquery 退款订单查询
type Refundquery struct {
	XMLName     xml.Name `xml:"xml,omitempty"`
	Appid       string   `xml:"appid,omitempty"`
	Mchid       string   `xml:"mch_id,omitempty"`
	NonceStr    string   `xml:"nonce_str,omitempty"`
	Sign        string   `xml:"sign,omitempty"`
	Signtype    string   `xml:"sign_type,omitempty"`
	RefundID    string   `xml:"refund_id,omitempty"`
	OutRefundNo string   `xml:"out_refund_no,omitempty"`
	Offset      int      `xml:"offset,omitempty"`
}

//Refund 退款信息
type Refund struct {
	RefundID            string `xml:"refund_id,omitempty"`     //微信退款单号
	OutRefundNo         string `xml:"out_refund_no,omitempty"` //商户退款单号
	RefundFee           string `xml:"refund_fee,omitempty"`
	SettlementRefundFee int64  `xml:"settlement_refund_fee,omitempty"`
	RefundStatus        string `xml:"refund_status,omitempty"` //SUCCESS-退款成功 CHANGE-退款异常REFUNDCLOSE—退款关闭
	SuccessTime         string `xml:"success_time,omitempty"`
	RefundRecvAccount   string `xml:"refund_recv_account,omitempty"` //退款入账账户
	RefundAccount       string `xml:"refund_account,omitempty"`      //退款资金来源
}

//RefundqueryRsp 退款信息查询返回
//！！！如果微信修改返回字段，验签会失败！！！
// https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_5&index=5
type RefundqueryRsp struct {
	XMLName            xml.Name `xml:"xml,omitempty"`
	ResultCode         string   `xml:"result_code,omitempty"`
	ErrCode            string   `xml:"err_code,omitempty"`
	ErrCodeDes         string   `xml:"err_code_des,omitempty"`
	Appid              string   `xml:"appid,omitempty"`
	MchID              string   `xml:"mch_id,omitempty"`
	NonceStr           string   `xml:"nonce_str,omitempty"`
	Sign               string   `xml:"sign,omitempty"`
	TotalRefundCount   int64    `xml:"total_refund_count,omitempty"`
	TransactionID      string   `xml:"transaction_id,omitempty"` //微信订单号
	OutTradeNo         string   `xml:"out_trade_no,omitempty"`   //商户系统内部的订单号
	TotalFee           int64    `xml:"total_fee,omitempty"`
	SettlementTotalFee int64    `xml:"settlement_total_fee,omitempty"`
	FeeType            string   `xml:"fee_type,omitempty"` //微信退款单号
	CashFee            int64    `xml:"cash_fee,omitempty"`
	RefundCount        int64    `xml:"refund_count,omitempty"` //退款笔数
	Refunds            []Refund //退款记录
	ReturnCode
}

//GetPublicKey 查询publickey
//https://pay.weixin.qq.com/wiki/doc/api/tools/mch_pay.php?chapter=24_7
type GetPublicKey struct {
	Mchid    string `xml:"mch_id,omitempty"`
	NonceStr string `xml:"nonce_str,omitempty"`
	Sign     string `xml:"sign,omitempty"`
	Signtype string `xml:"sign_type,omitempty"`
}

//GetPublicKeyRsp 查询publickey返回
type GetPublicKeyRsp struct {
	XMLName    xml.Name `xml:"xml,omitempty"`
	ResultCode string   `xml:"result_code,omitempty"`
	ErrCode    string   `xml:"err_code,omitempty"`
	ErrCodeDes string   `xml:"err_code_des,omitempty"`
	Mchid      string   `xml:"mch_id,omitempty"`
	PubKey     string   `xml:"pub_key,omitempty"`
	ReturnCode
}
