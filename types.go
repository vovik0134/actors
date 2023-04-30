package triggerable

import (
	"context"
	"time"
)

type action interface {
	Name() string
	Run(ctx context.Context) error
	RetryAfterTimeout() time.Duration
}

type logger interface {
	Debug(ctx context.Context, args ...any)
	Info(ctx context.Context, args ...any)
}
