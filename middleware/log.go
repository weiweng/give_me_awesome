package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"give_me_awesome/logs"
)

func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		logs.Infof("| %3d | %13v | %s | %s |",
			statusCode,
			latencyTime,
			reqMethod,
			reqUri,
		)
	}
}
