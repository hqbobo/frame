package encode

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/config/encoder"
)

var js = jsoniter.ConfigCompatibleWithStandardLibrary

type jsonEncoder struct{}

func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return js.Marshal(v)
}

func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return js.Unmarshal(d, v)
}

func (j jsonEncoder) String() string {
	return "json"
}

func NewJEncoder() encoder.Encoder {
	return jsonEncoder{}
}
