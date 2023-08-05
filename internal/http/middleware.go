package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionContext struct {
	UserID int64
}

func (s *Server) AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie(AuthCookieKey)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "session cookie not found"})
			return
		}

		user := SessionContext{}

		err = s.cache.GetStruct(c.Request.Context(), authCookie, &user)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if user.UserID < 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(gin.AuthUserKey, user.UserID)
		c.Next()
	}
}
