package middlewares

import (
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/gin-gonic/gin"
)

func DelayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.AppConfig.Delay != 0 {
			time.Sleep(time.Duration(utils.AppConfig.Delay) * time.Millisecond)
		}
		c.Next()
	}
}
