package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin middleware for logging requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return formatLog(param)
	})
}

func formatLog(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return "[AUTH-SERVICE] " +
		param.TimeStamp.Format("2006/01/02 - 15:04:05") +
		" | " + intToString(param.StatusCode) +
		" | " + param.Latency.String() +
		" | " + param.ClientIP +
		" | " + param.Method +
		" " + param.Path +
		" " + param.ErrorMessage + "\n"
}

func intToString(n int) string {
	return string(rune('0'+(n/100)%10)) + string(rune('0'+(n/10)%10)) + string(rune('0'+n%10))
}
