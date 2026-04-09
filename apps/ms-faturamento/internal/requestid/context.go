package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type contextKey string

const key contextKey = "request_id"

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(key).(string)
	return v
}

func New() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(b)
}
