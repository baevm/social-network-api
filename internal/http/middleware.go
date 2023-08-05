package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (s *Server) ZapLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if query != "" {
			path = path + "?" + query
		}
		param.Path = path

		// fields := []zapcore.Field{
		// 	zap.Int("status", c.Writer.Status()),
		// 	zap.String("method", c.Request.Method),
		// 	zap.String("path", path),
		// 	zap.String("query", query),
		// 	zap.String("ip", c.ClientIP()),
		// 	zap.String("user-agent", c.Request.UserAgent()),
		// 	zap.Duration("latency", param.Latency),
		// }

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Errorw(e,
					"status", c.Writer.Status(),
					"method", c.Request.Method,
				)
			}
		} else {
			logger.Infow(path,
				"status", c.Writer.Status(),
				"method", c.Request.Method)
		}
	}
}
