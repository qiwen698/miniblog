package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/pkg/known"
	"github.com/qiwen698/miniblog/pkg/core"
	"github.com/qiwen698/miniblog/pkg/errno"
	"github.com/qiwen698/miniblog/pkg/token"
)

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 jwt token
		username, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Set(known.XUsernameKey, username)
		c.Next()
	}
}
