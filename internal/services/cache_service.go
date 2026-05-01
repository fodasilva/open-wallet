package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

type CacheService interface {
	Set(ctx context.Context, domain string, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, domain string, key string) (any, error)
	Delete(ctx context.Context, domain string, key string) error
	Incr(ctx context.Context, domain string, key string, ttl time.Duration) (int, error)
}

type postgresCacheService struct {
	db *sql.DB
}

func NewCacheService(db *sql.DB) CacheService {
	return &postgresCacheService{db: db}
}

func (s *postgresCacheService) cleanup() {
	// #nosec G404 -- this is just for background cleanup probability, not cryptographic
	if rand.Float32() < 0.01 {
		go func() {
			_, err := s.db.Exec("DELETE FROM cache_entries WHERE expires_at < NOW()")
			if err != nil {
				log.Printf("failed to clean up cache: %v\n", err)
			}
		}()
	}
}

func (s *postgresCacheService) Set(ctx context.Context, domain string, key string, value any, ttl time.Duration) error {
	defer s.cleanup()
	expiresAt := time.Now().Add(ttl)
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO cache_entries (domain, key, value, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (domain, key) DO UPDATE
		SET value = $3, expires_at = $4;
	`
	_, err = s.db.ExecContext(ctx, query, domain, key, b, expiresAt)
	return err
}

func (s *postgresCacheService) Get(ctx context.Context, domain string, key string) (any, error) {
	defer s.cleanup()
	query := `
		SELECT value FROM cache_entries 
		WHERE domain = $1 AND key = $2 AND expires_at >= NOW();
	`
	var b []byte
	err := s.db.QueryRowContext(ctx, query, domain, key).Scan(&b)
	if err != nil {
		return nil, err
	}

	var val any
	if err := json.Unmarshal(b, &val); err != nil {
		return nil, err
	}
	return val, nil
}

func (s *postgresCacheService) Delete(ctx context.Context, domain string, key string) error {
	query := `DELETE FROM cache_entries WHERE domain = $1 AND key = $2;`
	_, err := s.db.ExecContext(ctx, query, domain, key)
	return err
}

func (s *postgresCacheService) Incr(ctx context.Context, domain string, key string, ttl time.Duration) (int, error) {
	defer s.cleanup()
	expiresAt := time.Now().Add(ttl)
	query := `
		INSERT INTO cache_entries (domain, key, count, expires_at)
		VALUES ($1, $2, 1, $3)
		ON CONFLICT (domain, key) DO UPDATE
		SET count = CASE WHEN cache_entries.expires_at < NOW() THEN 1 ELSE COALESCE(cache_entries.count, 0) + 1 END,
		    expires_at = CASE WHEN cache_entries.expires_at < NOW() THEN $3 ELSE cache_entries.expires_at END
		RETURNING count;
	`
	var currentCount int
	err := s.db.QueryRowContext(ctx, query, domain, key, expiresAt).Scan(&currentCount)
	return currentCount, err
}
