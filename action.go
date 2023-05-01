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

func (a *actionImpl) RetryOnError(err error) (bool, func(ctx context.Context)) {
	return a.retryOnError(err)
}

func Action(runFunc func(ctx context.Context) error, opts ...actionOption) *actionImpl {
	a := &actionImpl{
		runFunc: runFunc,
		name:    "unnamed",
		retryOnError: func(err error) (bool, func(ctx context.Context)) {
			return false, nil
		},
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

func WithRetryOnError(retryOnError func(err error) (bool, func(ctx context.Context))) actionOption {
	return func(a *actionImpl) {
		a.retryOnError = retryOnError
	}
}

type actionOption func(a *actionImpl)

type actionImpl struct {
	runFunc      func(ctx context.Context) error
	name         string
	retryOnError func(err error) (bool, func(ctx context.Context))
}

type action interface {
	Name() string
	Run(ctx context.Context) error
	RetryOnError(err error) (retry bool, retryFunc func(ctx context.Context))
}
