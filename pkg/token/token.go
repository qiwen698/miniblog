package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	key         string
	identityKey string
}

var ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")
var (
	config = Config{
		key:         "gsdfsdKKxon",
		identityKey: "identityKey",
	}
	once sync.Once
)

func Init(key string, identityKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
	})
}

func Parse(tokenString string, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//确保token 加密算法是预期的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})
	//解析失败
	if err != nil {
		return "", err
	}
	var identityKey string
	// 如果解析成功，从token 中取出 token 的主题
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
	}
	return identityKey, nil
}

func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", ErrMissingHeader
	}
	var t string
	// 从请求头中取出token
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, config.key)
}

func Sign(identityKey string) (tokenString string, err error) {
	// token 的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(100000 * time.Hour).Unix(),
	})
	//签发 token
	tokenString, err = token.SignedString([]byte(config.key))
	return
}
