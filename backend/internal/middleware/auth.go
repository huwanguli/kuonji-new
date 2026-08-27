package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuonji/blog/config"
	"github.com/kuonji/blog/pkg/jwt"
	"github.com/kuonji/blog/pkg/response"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 移除 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		cfg := config.Get()
		j := jwt.New(cfg.JWT.Secret, cfg.JWT.Expire)
		claims, err := j.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
