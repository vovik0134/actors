package triggerable

import "context"

type logger interface {
	Debug(ctx context.Context, args ...any)
	Info(ctx context.Context, args ...any)
}
