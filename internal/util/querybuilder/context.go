package querybuilder

import "context"

type contextKey string

const builderKey contextKey = "query_builder"

// WithBuilder returns a new context with the given builder attached.
func WithBuilder(ctx context.Context, b *Builder) context.Context {
	return context.WithValue(ctx, builderKey, b)
}

// Get extracts the builder from the context, or returns nil if not found.
func Get(ctx context.Context) *Builder {
	if b, ok := ctx.Value(builderKey).(*Builder); ok {
		return b
	}
	return nil
}
