package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/felipe1496/open-wallet/infra"
)

func DelayMiddleware(cfg *infra.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Delay > 0 {
			time.Sleep(time.Duration(cfg.Delay) * time.Millisecond)
		}
		c.Next()
	}
}
