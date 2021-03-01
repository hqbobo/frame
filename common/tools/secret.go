package tools

import (
	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
	"encoding/hex"
)

//EncryPwdViaBcrypt 使用bcrypt对密码进行加密
func EncryPwdViaBcrypt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	return encodePW, nil
}

//DecryPwdViaBcrypt 将加密后的字符串解密与pwd比较，相同则返回true
func DecryPwdViaBcrypt(encodePw, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePw), []byte(pwd))
	if err != nil {
		return false
	} else {
		return true
	}
}

func Md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}