package middleware

import (
	"net/http"
	"strings"
	"whiteboard/server/internal/auth"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

// AuthRequired 认证中间件
func AuthRequired(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		// 检查Authorization头是否存在
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		// 检查Authorization头是否以Bearer "开头
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(header, prefix))

		// 解析token字符串，这里的token是一个*jwt.Token结构体
		claims, err := jwtManager.Parse(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		// 放进上下文
		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// CurrentUserID 获取当前用户ID
// 符合谁设置谁提供的原则
func CurrentUserID(c *gin.Context) string {
	v, ok := c.Get(UserIDKey)
	if !ok {
		return ""
	}
	userID, _ := v.(string)
	return userID
}
