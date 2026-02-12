package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/utils"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	db *sql.DB
}

func NewRateLimiter(db *sql.DB) RateLimiter {
	return RateLimiter{db: db}
}

func (r RateLimiter) Middleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		ctx := context.Background()

		// Limpar registros expirados periodicamente (opcional, mas recomendado)
		go r.cleanupExpiredRecords(ctx)

		count, resetTime, err := r.incrementAndGet(ctx, clientIP, window)
		if err != nil {
			apiErr := utils.NewHTTPError(http.StatusInternalServerError, "It was not possible to check user resource limits")
			c.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			return
		}

		remaining := max(limit-count, 0)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		if count > limit {
			apiErr := utils.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			c.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
			return
		}

		c.Next()
	}
}

func (r RateLimiter) incrementAndGet(ctx context.Context, clientIP string, window time.Duration) (int, time.Time, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, time.Time{}, err
	}
	defer tx.Rollback()

	now := time.Now()
	expiresAt := now.Add(window)

	query := `
        INSERT INTO rate_limits (client_ip, request_count, window_start, expires_at)
        VALUES ($1, 1, $2, $3)
        ON CONFLICT (client_ip) 
        DO UPDATE SET 
            request_count = CASE 
                WHEN rate_limits.expires_at > $2 THEN rate_limits.request_count + 1
                ELSE 1
            END,
            window_start = CASE 
                WHEN rate_limits.expires_at > $2 THEN rate_limits.window_start
                ELSE $2
            END,
            expires_at = CASE 
                WHEN rate_limits.expires_at > $2 THEN rate_limits.expires_at
                ELSE $3
            END
        RETURNING request_count, expires_at
    `

	var count int
	var resetTime time.Time

	err = tx.QueryRowContext(ctx, query, clientIP, now, expiresAt).Scan(&count, &resetTime)
	if err != nil {
		return 0, time.Time{}, err
	}

	if err := tx.Commit(); err != nil {
		return 0, time.Time{}, err
	}

	return count, resetTime, nil
}

func (r RateLimiter) cleanupExpiredRecords(ctx context.Context) {
	query := `DELETE FROM rate_limits WHERE expires_at < $1`
	_, _ = r.db.ExecContext(ctx, query, time.Now())
}

func (r RateLimiter) StartCleanupJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			r.cleanupExpiredRecords(context.Background())
		}
	}()
}
