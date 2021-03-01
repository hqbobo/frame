package encode

import (
	"github.com/micro/go-micro/config/encoder"
	"github.com/micro/go-micro/config/encoder/json"
)

var encode encoder.Encoder

//使用json做encoder
func init() {
	encode = json.NewEncoder()
}

// Encoder 获取Encoder
func Encoder() encoder.Encoder{
	return encode
}

