package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/musorok/server/internal/core/auth"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized); return
		}
		tok := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.ParseClaims(secret, tok)
		if err != nil { c.AbortWithStatus(http.StatusUnauthorized); return }
		c.Set("uid", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
