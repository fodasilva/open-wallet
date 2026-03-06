package middlewares

import (
	"time"

	"github.com/felipe1496/open-wallet/infra"
	"github.com/gin-gonic/gin"
)

func DelayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if infra.AppConfig.Delay != 0 {
			time.Sleep(time.Duration(infra.AppConfig.Delay) * time.Millisecond)
		}
		c.Next()
	}
}
