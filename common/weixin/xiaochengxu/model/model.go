package model

type ErrRsp struct {
	Errmsg  string `json:"errmsg"`
	Errcode int64  `json:"errcode"`
}

type AccessTokenRsp struct {
	Access_token string `json:"access_token"`
	Expires_in   int64  `json:"expires_in"`
	ErrRsp
}

type SubscribeMessage struct {
	Touser            string      `json:"touser"`
	Template_id       string      `json:"template_id"`
	Page              string      `json:"page,omitempty"`
	Miniprogram_state string      `json:"miniprogram_state,omitempty"` //developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Data              interface{} `json:"data"`
}

type MiniToken struct {
	ErrRsp
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
}

type SubcribeValue struct {
	Value string `json:"value"`
}

const SubcribeDateFormat = "2006年01月02日 15:04"

//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
type QrcodeGetUnlimitedReq struct {
	Scene      string `json:"scene"`
	Page       string `json:"page"`
	Checkpath  bool   `json:"check_path"` //检查page 是否存在，为 true 时 page 必须是已经发布的小程序存在的页面（否则报错）；为 false 时允许小程序未发布或者 page 不存在， 但page 有数量上限（60000个）请勿滥用
	Envversion string `json:"envVersion"`
}
