package tools

import (
	"errors"
	"fmt"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/hqbobo/frame/common/aes"
	"github.com/hqbobo/frame/common/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

type jwtCustomClaims struct {
	jwt.StandardClaims
	// 追加自己需要的信息
	UID int64
	//Admin bool
}

//SecretKey 加密密钥
var SecretKey = []byte("SecretKey")

//CreateToken 创建token
func CreateToken(uid int64) (tokenString string, err error) {
	claims := &jwtCustomClaims{
		UID: uid,
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()
	claims.Issuer = utils.NewLenChars(10)
	//claims := jwt.StandardClaims{
	//	ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
	//	Issuer: issuer,
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

//ParseToken 解析token
func ParseToken(tokenSrt string) (uid int64, err error) {
	//var token jwt.Token
	token, err := jwt.Parse(tokenSrt, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if token == nil {
		return 0, errors.New("token nil ")
	}
	claims := token.Claims
	fmt.Println("ParseToken uid : ", claims.(jwt.MapClaims))
	if claims.(jwt.MapClaims)["UID"] == nil {
		err = errors.New("token err : " + tokenSrt)
		return
	}
	uidF := claims.(jwt.MapClaims)["UID"].(float64)
	uid = int64(uidF)
	fmt.Println("uid : ", uid)
	return
}

type Token struct {
	Uid    int64
	Sign   string //签名
	Indate int64  //有效期
	Time   int64  //时间
}

/*
*初始化一个Token
 */
func NewToken(uid int64, sign, aesKey string, indate int64) (string, error) {
	token := &Token{
		Uid:    uid,
		Sign:   sign,
		Indate: indate,
		Time:   time.Now().Unix(),
	}
	ciphertext, err := token.Encrypt(aesKey)
	if err != nil {
		return "", err
	}
	return ciphertext, nil
}

/*
*验证token
 */
func ValidateToken(token, sign, aesKey string) (*Token, error) {
	t := &Token{}
	err := t.Decrypt(token, aesKey)
	if err != nil {
		return t, errors.New("非法token")
	}
	if time.Now().Unix()-t.Time > t.Indate {
		return t, errors.New("token已过期")
	}
	if t.Sign != sign {
		return t, errors.New("签名错误")
	}
	return t, nil
}

/*
*加密一个验证码
 */
func (t *Token) Encrypt(aesKey string) (string, error) {
	key := aesKey
	enc := &aes.AesEncrypt{}
	buffer, err := json.Marshal(*t)
	if err != nil {
		return "", err
	}
	src, _ := enc.Encrypt(key, string(buffer))
	encodeString := base64.StdEncoding.EncodeToString(src)
	return encodeString, nil
}

/*
*解密一个验证码
 */
func (t *Token) Decrypt(src, aesKey string) error {
	decodeBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return err
	}
	key := aesKey
	enc := &aes.AesEncrypt{}
	dst, _ := enc.Decrypt(key, decodeBytes)
	e := json.Unmarshal(dst, t)
	if e != nil {
		return e
	}
	return nil
}
