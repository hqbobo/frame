package model

type miniPath struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}
type WeixinTemp struct {
	ToUser     string      `json:"touser"`
	TemplateId string      `json:"template_id"`
	Url        string      `json:"url"`
	Topcolor   string      `json:"topcolor"`
	Data       interface{} `json:"data"`
	Mini       miniPath    `json:"miniprogram"`
}

type WeixinTempItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type WeixinKfText struct {
	ToUser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}
