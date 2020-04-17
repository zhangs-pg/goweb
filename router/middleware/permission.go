package middleware

import (
	"github.com/gin-gonic/gin"
)

func CheckPerm(perm uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		permission := claims.(*CustomClaims).Permission
		if (permission & perm) == 0 {
			// 无权限
			c.JSON(403, gin.H{
				"status": -1,
				"msg":    "无权限访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
