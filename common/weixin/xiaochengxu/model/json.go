package model

type Watermark struct {
	Timestamp int64 `json:"timestamp"`
	Appid     string    `json:"appid"`
}

type UserWX struct {
	OpenId    string      `json:"openId"`
	NickName  string      `json:"nickName"`
	Gender    int         `json:"gender"`
	Language  string      `json:"language"`
	City      string      `json:"city"`
	Province  string      `json:"province"`
	Country   string      `json:"country"`
	AvatarUrl string      `json:"avatarUrl"`
	Watermark []Watermark `json:"watermark"`
	Unionid   string      `json:"unionId"`
}

type UserPhone struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string       `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}
