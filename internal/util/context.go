package util

import (
	"context"
	"fmt"
)

type ContextKey string

const (
	ContextKeyUserID       ContextKey = "user_id"
	ContextKeyPage         ContextKey = "page"
	ContextKeyPerPage      ContextKey = "per_page"
	ContextKeyQueryBuilder ContextKey = "query_builder"
)

// Get is a generic helper to retrieve values from context.
func Get[T any](ctx context.Context, key ContextKey) (T, bool) {
	val, ok := ctx.Value(key).(T)
	return val, ok
}

// MustGet is a generic helper that panics if the key is not found.
func MustGet[T any](ctx context.Context, key ContextKey) T {
	val, ok := Get[T](ctx, key)
	if !ok {
		panic(fmt.Sprintf("key %s not found in context or type mismatch", key))
	}
	return val
}

// GetString mimics Gin's GetString behavior.
func GetString(ctx context.Context, key ContextKey) string {
	val, _ := ctx.Value(key).(string)
	return val
}

// GetInt mimics Gin's GetInt behavior.
func GetInt(ctx context.Context, key ContextKey) int {
	val, _ := ctx.Value(key).(int)
	return val
}
