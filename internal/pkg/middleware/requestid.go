package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qiwen698/miniblog/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(known.XRequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// 将RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(known.XRequestIDKey, requestID)
		// 将RequestID 保存在 HTTP 返回头中，Header 的键位 `X-Request-ID`
		c.Writer.Header().Set(known.XRequestIDKey, requestID)
		c.Next()
	}
}
