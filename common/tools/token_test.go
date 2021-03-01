package tools_test

import (
	"github.com/hqbobo/frame/common/tools"
	"testing"
)

const (
	Sign   = "sign123456123456"
	AesKey = "aes_key12345123456"
)

func TestValidateToken(t *testing.T) {
	token, err := tools.ValidateToken("pZVrbRk2pakA/qi8AezGZspnsqB5pfvM+XTlXv62jO8KhSsLNjYAxni164cbqED6AvZc9tkKe6Fe4b4E6MLiBF9h9kv7wROhVd8zyqHHLA==", Sign, AesKey)
	if err != nil {
		panic(err)
	}
	t.Log("token = ", token)

	token, err = tools.ValidateToken("pZVrbRk2pakA/qi8AezGZspnsqB5pfvM+XTlXv62jO8KhSsLNjYAxni164cbqED6AvZc9tkKe6Fe4b4E6MLiBF9h9kv7wROhV9A2x6DDLA==", Sign, AesKey)
	if err != nil {
		panic(err)
	}
	t.Log("token = ", token)
}
