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
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return "[API-GATEWAY] " +
		param.TimeStamp.Format("2006/01/02 - 15:04:05") +
		" |" + statusColor + " " + string(rune(param.StatusCode)) + " " + resetColor +
		"| " + param.Latency.String() +
		" | " + param.ClientIP +
		" |" + methodColor + " " + param.Method + " " + resetColor +
		" " + param.Path +
		" " + param.ErrorMessage + "\n"
}
