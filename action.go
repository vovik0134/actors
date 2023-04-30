package triggerable

import (
	"context"
	"time"
)

const NoRetryTimeout = -1 * time.Second

func (a *actionImpl) Run(ctx context.Context) error {
	return a.runFunc(ctx)
}

func (a *actionImpl) Name() string {
	return a.name
}

func (a *actionImpl) RetryAfterTimeout() time.Duration {
	return a.retryAfterTimeout
}

func Action(runFunc func(ctx context.Context) error, opts ...actionOption) *actionImpl {
	a := &actionImpl{
		runFunc:           runFunc,
		name:              "unnamed",
		retryAfterTimeout: NoRetryTimeout,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

func WithName(name string) actionOption {
	return func(a *actionImpl) {
		a.name = name
	}
}

func WithRetryAfterTimeout(retryAfterTimeout time.Duration) actionOption {
	return func(a *actionImpl) {
		a.retryAfterTimeout = retryAfterTimeout
	}
}

type actionOption func(a *actionImpl)

type actionImpl struct {
	runFunc           func(ctx context.Context) error
	name              string
	retryAfterTimeout time.Duration
}
