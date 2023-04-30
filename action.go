package triggerable

import (
	"context"
)

func (a *actionImpl) Run(ctx context.Context) error {
	return a.runFunc(ctx)
}

func (a *actionImpl) Name() string {
	return a.name
}

func (a *actionImpl) RetryOnError(ctx context.Context, err error) bool {
	return a.retryOnError(ctx, err)
}

func Action(runFunc func(ctx context.Context) error, opts ...actionOption) *actionImpl {
	a := &actionImpl{
		runFunc:      runFunc,
		name:         "unnamed",
		retryOnError: func(ctx context.Context, err error) bool { return false },
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

func WithRetryOnError(retryOnError func(ctx context.Context, err error) bool) actionOption {
	return func(a *actionImpl) {
		a.retryOnError = retryOnError
	}
}

type actionOption func(a *actionImpl)

type actionImpl struct {
	runFunc      func(ctx context.Context) error
	name         string
	retryOnError func(ctx context.Context, err error) bool
}

type action interface {
	Name() string
	Run(ctx context.Context) error
	RetryOnError(ctx context.Context, err error) bool
}
